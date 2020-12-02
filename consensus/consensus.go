package consensus

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"github.com/joe-zxh/lbbft/internal/proto"
	"github.com/joe-zxh/lbbft/util"
	"go.uber.org/atomic"
	"log"
	"sync"
	"time"

	"github.com/joe-zxh/lbbft/config"
	"github.com/joe-zxh/lbbft/data"
	"github.com/joe-zxh/lbbft/internal/logging"
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

	// from pbft
	Mut   sync.Mutex // Lock for all internal data
	ID    uint32
	TSeq  atomic.Uint32 // Total sequence number of next request
	View  uint32        // View和Seq都是从1开始的
	Apply uint32        // Sequence number of last executed request

	Log                   []*data.Entry
	LogMut                sync.Mutex
	LastPreparedID        data.EntryID
	UpdateLastPreparedMut sync.Mutex

	// Log        map[data.EntryID]*data.Entry // bycon的log是一个数组，因为需要保证连续，leader可以处理log inconsistency，而pbft不需要。client只有执行完上一条指令后，才会发送下一条请求，所以顺序 并没有问题。
	cps       map[int]*CheckPoint
	WaterLow  uint32
	WaterHigh uint32
	F         uint32
	Q         uint32
	N         uint32
	monitor   bool
	Change    *time.Timer
	Changing  bool        // Indicate if this node is changing view
	state     interface{} // Deterministic state machine's state
	vcs       map[uint32][]*ViewChangeArgs
	lastcp    uint32

	Leader   uint32 // view改变的时候，再改变
	IsLeader bool   // view改变的时候，再改变

	waitEntry      *sync.Cond
	ViewChangeChan chan struct{}
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
		View:     0,
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

		// from pbft
		ID:        uint32(conf.ID),
		View:      1,
		Apply:     0,
		cps:       make(map[int]*CheckPoint),
		WaterLow:  0,
		WaterHigh: 2 * checkpointDiv,
		F:         uint32(len(conf.Replicas)-1) / 3,
		N:         uint32(len(conf.Replicas)),
		monitor:   false,
		Change:    nil,
		Changing:  false,
		state:     make([]interface{}, 1),
		vcs:       make(map[uint32][]*ViewChangeArgs),
		lastcp:    0,

		LastPreparedID: data.EntryID{V: 0, N: 0},
		ViewChangeChan: make(chan struct{}, 1),
	}

	lbbft.InitLog()

	lbbft.waitEntry = sync.NewCond(&lbbft.LogMut)

	lbbft.Q = (lbbft.F+lbbft.N)/2 + 1
	lbbft.Leader = 1
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

func (lbbft *LBBFTCore) ApplyCommands(commitSeq uint32) {
	lbbft.Mut.Lock()
	defer lbbft.Mut.Unlock()

	lbbft.LogMut.Lock()
	defer lbbft.LogMut.Unlock()

	for lbbft.Apply < commitSeq {
		if lbbft.Apply+1 < uint32(len(lbbft.Log)) {
			ent := lbbft.Log[lbbft.Apply+1]
			lbbft.Exec <- *ent.PP.Commands
			lbbft.Apply++
		} else {
			break
		}
	}
}

func (lbbft *LBBFTCore) UpdateLastPreparedID(ent *data.Entry) {
	lbbft.UpdateLastPreparedMut.Lock()
	defer lbbft.UpdateLastPreparedMut.Unlock()

	if lbbft.LastPreparedID.IsOlder(&data.EntryID{V: ent.PP.View, N: ent.PP.Seq}) {
		lbbft.LastPreparedID.V = ent.PP.View
		lbbft.LastPreparedID.N = ent.PP.Seq
	}
}

func (lbbft *LBBFTCore) IsRequestVotePreparedCertValid(pRV *proto.RequestVoteArgs) bool {

	if !lbbft.SigCache.VerifyQuorumCert(pRV.PreparedCert.Proto2QuorumCert()) {
		return false
	}

	// 检查sigContent是否=hash(view+seq+"prepare"+Digest)
	var dDigest data.EntryHash
	copy(dDigest[:], pRV.Digest[:len(dDigest)])

	s512 := sha512.New()

	byte4 := make([]byte, 4)
	binary.LittleEndian.PutUint32(byte4, uint32(pRV.LastPreparedView))
	s512.Write(byte4[:])

	binary.LittleEndian.PutUint32(byte4, uint32(pRV.LastPreparedSeq))
	s512.Write(byte4[:])

	s512.Write([]byte("prepare"))
	s512.Write(pRV.Digest)

	sum := s512.Sum(nil)

	if bytes.Equal(sum, pRV.PreparedCert.SigContent) {
		return true
	} else {
		panic("IsRequestVotePreparedCertValid: check prepare hash do not pass...")
	}

	return false
}

func (lbbft *LBBFTCore) IsNewViewCertValid(pNV *proto.NewViewArgs) bool {

	if !lbbft.SigCache.VerifyQuorumCert(pNV.NewViewCert.Proto2QuorumCert()) {
		log.Println(`VerifyQuorumCert failed...`)
		return false
	}

	// 检查sigContent是否=hash(view+candidate id)
	s512 := sha512.New()

	byte4 := make([]byte, 4)
	binary.LittleEndian.PutUint32(byte4, uint32(pNV.NewView))
	s512.Write(byte4[:])

	binary.LittleEndian.PutUint32(byte4, uint32(pNV.CandidateID))
	s512.Write(byte4[:])

	sum := s512.Sum(nil)

	if bytes.Equal(sum, pNV.NewViewCert.SigContent) {
		return true
	} else {
		panic("IsNewViewCertValid: sigContent不等于hash(view+candidate id)...")
	}

	return false
}

// return sig, sigcontent
func (lbbft *LBBFTCore) CreateRequestVoteReplySig(pRV *proto.RequestVoteArgs) (*proto.PartialSig, *([]byte)) {

	s512 := sha512.New()

	byte4 := make([]byte, 4)
	binary.LittleEndian.PutUint32(byte4, uint32(pRV.NewView))
	s512.Write(byte4[:])

	binary.LittleEndian.PutUint32(byte4, uint32(pRV.CandidateID))
	s512.Write(byte4[:])

	sum := s512.Sum(nil)

	ps, err := lbbft.SigCache.CreatePartialSig(lbbft.Config.ID, lbbft.Config.PrivateKey, sum)
	util.PanicErr(err)
	return proto.PartialSig2Proto(ps), &sum
}
