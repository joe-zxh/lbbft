package consensus

import (
	"context"
	"fmt"
	"go.uber.org/atomic"
	"log"
	"sync"
	"time"

	"github.com/joe-zxh/lbbft/config"
	"github.com/joe-zxh/lbbft/data"
	"github.com/joe-zxh/lbbft/internal/logging"
	"github.com/joe-zxh/lbbft/util"
)

const (
	changeViewTimeout = 60 * time.Second
	checkpointDiv     = 2000000
)

var logger *log.Logger

func init() {
	logger = logging.GetLogger()
}

// LBBFTCore is the safety core of the LBBFTCore protocol
type LBBFTCore struct {
	// from hotstuff

	cmdCache *data.CommandSet // Contains the commands that are waiting to be proposed
	Config   *config.ReplicaConfig
	SigCache *data.SignatureCache
	cancel   context.CancelFunc // stops goroutines

	Exec chan []data.Command

	// from lbbft
	Mut    sync.Mutex // Lock for all internal data
	ID     uint32
	TSeq   atomic.Uint32           // Total sequence number of next request
	seqmap map[data.EntryID]uint32 // Use to map {Cid,CSeq} to global sequence number for all prepared message
	View   uint32
	Apply  uint32 // Sequence number of last executed request

	Log    []*data.Entry
	LogMut sync.Mutex

	// Log        map[data.EntryID]*data.Entry // bycon的log是一个数组，因为需要保证连续，leader可以处理log inconsistency，而pbft不需要。client只有执行完上一条指令后，才会发送下一条请求，所以顺序 并没有问题。
	cps        map[int]*CheckPoint
	WaterLow   uint32
	WaterHigh  uint32
	F          uint32
	Q          uint32
	N          uint32
	monitor    bool
	Change     *time.Timer
	Changing   bool                // Indicate if this node is changing view
	state      interface{}         // Deterministic state machine's state
	ApplyQueue *util.PriorityQueue // 因为PBFT的特殊性(log是一个map，而不是list)，所以这里需要一个applyQueue。
	vcs        map[uint32][]*ViewChangeArgs
	lastcp     uint32

	Leader   uint32 // view改变的时候，再改变
	IsLeader bool   // view改变的时候，再改变

	waitEntry *sync.Cond
}

func (lbbft *LBBFTCore) AddCommand(command data.Command) {
	lbbft.cmdCache.Add(command)
}

func (lbbft *LBBFTCore) CommandSetLen(command data.Command) int {
	return lbbft.cmdCache.Len()
}

// CreateProposal creates a new proposal
func (lbbft *LBBFTCore) GetProposeCommands(timeout bool) *([]data.Command) {

	var batch []data.Command

	if timeout { // timeout的时候，不管够不够batch都要发起共识。
		batch = lbbft.cmdCache.RetriveFirst(lbbft.Config.BatchSize)
	} else {
		batch = lbbft.cmdCache.RetriveExactlyFirst(lbbft.Config.BatchSize)
	}

	return &batch
}

func (lbbft *LBBFTCore) InitLog() {
	lbbft.LogMut.Lock()
	defer lbbft.LogMut.Unlock()
	dPP := data.PrePrepareArgs{
		View:     1,
		Seq:      0,
		Commands: nil,
	}

	nonce := &data.Entry{
		PP:            &dPP,
		PreparedCert:  nil,
		Prepared:      true,
		CommittedCert: nil,
		Committed:     true,
		PreEntryHash:  &data.EntryHash{},
		Digest:        &data.EntryHash{},
		PrepareHash:   &data.EntryHash{},
		CommitHash:    &data.EntryHash{},
	}

	lbbft.Log = append([]*data.Entry{}, nonce)
}

// New creates a new LBBFTCore instance
func New(conf *config.ReplicaConfig) *LBBFTCore {
	logger.SetPrefix(fmt.Sprintf("hs(id %d): ", conf.ID))

	ctx, cancel := context.WithCancel(context.Background())

	lbbft := &LBBFTCore{
		// from hotstuff
		Config:   conf,
		cancel:   cancel,
		SigCache: data.NewSignatureCache(conf),
		cmdCache: data.NewCommandSet(),
		Exec:     make(chan []data.Command, 1),

		// lbbft
		ID:         uint32(conf.ID),
		seqmap:     make(map[data.EntryID]uint32),
		View:       1,
		Apply:      0,
		cps:        make(map[int]*CheckPoint),
		WaterLow:   0,
		WaterHigh:  2 * checkpointDiv,
		F:          uint32(len(conf.Replicas)-1) / 3,
		N:          uint32(len(conf.Replicas)),
		monitor:    false,
		Change:     nil,
		Changing:   false,
		state:      make([]interface{}, 1),
		ApplyQueue: util.NewPriorityQueue(),
		vcs:        make(map[uint32][]*ViewChangeArgs),
		lastcp:     0,
	}

	lbbft.InitLog()

	lbbft.waitEntry = sync.NewCond(&lbbft.LogMut)

	lbbft.Q = lbbft.F*2 + 1
	lbbft.Leader = (lbbft.View-1)%lbbft.N + 1
	lbbft.IsLeader = (lbbft.Leader == lbbft.ID)

	// Put an initial stable checkpoint
	cp := lbbft.getCheckPoint(-1)
	cp.Stable = true
	cp.State = lbbft.state

	go lbbft.proposeConstantly(ctx)

	return lbbft
}

func (lbbft *LBBFTCore) proposeConstantly(ctx context.Context) {
	for {
		select {
		// todo: 一个计时器，如果是leader，就开始preprepare
		case <-ctx.Done():
			return
		}
	}
}

func (lbbft *LBBFTCore) Close() {
	lbbft.cancel()
}

func (lbbft *LBBFTCore) GetExec() chan []data.Command {
	return lbbft.Exec
}

func (lbbft *LBBFTCore) PutEntry(ent *data.Entry) {
	lbbft.LogMut.Lock()
	defer lbbft.LogMut.Unlock()

	listLen := uint32(len(lbbft.Log))

	if listLen > ent.PP.Seq {
		panic(`try to overwrite an entry`)
	}

	for listLen < ent.PP.Seq {
		lbbft.waitEntry.Wait()
		listLen = uint32(len(lbbft.Log))
	}

	// listLen == ent.PP.Seq
	ent.PreEntryHash = lbbft.Log[listLen-1].Digest
	ent.GetCommitHash() // 把Digest都算出来，以免后一个entry添加的时候 PreEntryHash是空的。记得在Put的外部，给ent加锁。
	lbbft.Log = append(lbbft.Log, ent)

	lbbft.waitEntry.Broadcast()
}

func (lbbft *LBBFTCore) GetEntryBySeq(seq uint32) *data.Entry {
	lbbft.LogMut.Lock()
	defer lbbft.LogMut.Unlock()

	listLen := uint32(len(lbbft.Log))

	for listLen <= seq {
		lbbft.waitEntry.Wait()
		listLen = uint32(len(lbbft.Log))
	}

	return lbbft.Log[seq]
}

func (lbbft *LBBFTCore) GetApplyCmds(commitSeq uint32) (*[]data.Command, uint32) {
	lbbft.LogMut.Lock()
	defer lbbft.LogMut.Unlock()

	commands := make([]data.Command, 0)

	for commitSeq < uint32(len(lbbft.Log)) {
		ent := lbbft.Log[commitSeq]
		if ent.Committed {
			commands = append(commands, *ent.PP.Commands...)
			commitSeq++
		} else {
			break
		}
	}
	return &commands, commitSeq - 1
}
