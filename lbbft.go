package lbbft

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/joe-zxh/lbbft/data"
	"github.com/joe-zxh/lbbft/util"
	"log"
	"math/big"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/joe-zxh/lbbft/config"
	"github.com/joe-zxh/lbbft/consensus"
	"github.com/joe-zxh/lbbft/internal/logging"
	"github.com/joe-zxh/lbbft/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var logger *log.Logger

func init() {
	logger = logging.GetLogger()
}

// LBBFT is a thing
type LBBFT struct {
	*consensus.LBBFTCore
	proto.UnimplementedLBBFTServer

	tls bool

	nodes map[config.ReplicaID]*proto.LBBFTClient
	conns map[config.ReplicaID]*grpc.ClientConn

	server *lbbftServer

	closeOnce      sync.Once
	connectTimeout time.Duration
}

//New creates a new backend object.
func New(conf *config.ReplicaConfig, tls bool, connectTimeout, qcTimeout time.Duration) *LBBFT {
	lbbft := &LBBFT{
		LBBFTCore:      consensus.New(conf),
		nodes:          make(map[config.ReplicaID]*proto.LBBFTClient),
		conns:          make(map[config.ReplicaID]*grpc.ClientConn),
		connectTimeout: connectTimeout,
	}
	return lbbft
}

//Start starts the server and client
func (lbbft *LBBFT) Start() error {
	addr := lbbft.Config.Replicas[lbbft.Config.ID].Address
	err := lbbft.startServer(addr)
	if err != nil {
		return fmt.Errorf("Failed to start GRPC Server: %w", err)
	}
	err = lbbft.startClient(lbbft.connectTimeout)
	if err != nil {
		return fmt.Errorf("Failed to start GRPC Clients: %w", err)
	}
	return nil
}

// 作为rpc的client端，调用其他hsserver的rpc。
func (lbbft *LBBFT) startClient(connectTimeout time.Duration) error {

	grpcOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithReturnConnectionError(),
	}

	if lbbft.tls {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(lbbft.Config.CertPool, "")))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}

	for rid, replica := range lbbft.Config.Replicas {
		if replica.ID != lbbft.Config.ID {
			conn, err := grpc.Dial(replica.Address, grpcOpts...)
			if err != nil {
				log.Fatalf("connect error: %v", err)
				conn.Close()
			} else {
				lbbft.conns[rid] = conn
				c := proto.NewLBBFTClient(conn)
				lbbft.nodes[rid] = &c
			}
		}
	}

	return nil
}

// startServer runs a new instance of lbbftServer
func (lbbft *LBBFT) startServer(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("Failed to listen to port %s: %w", port, err)
	}

	grpcServerOpts := []grpc.ServerOption{}

	if lbbft.tls {
		grpcServerOpts = append(grpcServerOpts, grpc.Creds(credentials.NewServerTLSFromCert(lbbft.Config.Cert)))
	}

	lbbft.server = newLBBFTServer(lbbft)

	s := grpc.NewServer(grpcServerOpts...)
	proto.RegisterLBBFTServer(s, lbbft.server)

	go s.Serve(lis)
	return nil
}

// Close closes all connections made by the LBBFT instance
func (lbbft *LBBFT) Close() {
	lbbft.closeOnce.Do(func() {
		lbbft.LBBFTCore.Close()
		for _, conn := range lbbft.conns { // close clients connections
			conn.Close()
		}
	})
}

// 这个server是面向 集群内部的。
type lbbftServer struct {
	*LBBFT

	mut     sync.RWMutex
	clients map[context.Context]config.ReplicaID
}

func (lbbft *LBBFT) Ordering(_ context.Context, pO *proto.OrderingArgs) (*proto.OrderingReply, error) {

	logger.Printf("Ordering: leader id:%d\n", lbbft.ID)

	if !lbbft.IsLeader {
		return &proto.OrderingReply{}, errors.New(`I am not the leader`)
	}

	tseq := lbbft.TSeq.Inc()

	ent := &data.Entry{
		PP: &data.PrePrepareArgs{
			View:     lbbft.View,
			Seq:      tseq,
			Commands: pO.GetDataCommands(),
		},
	}

	ent.Mut.Lock()
	lbbft.PutEntry(ent)
	ps, err := lbbft.SigCache.CreatePartialSig(lbbft.Config.ID, lbbft.Config.PrivateKey, ent.GetPrepareHash().ToSlice())
	if err != nil {
		panic(err)
	}
	ent.Mut.Unlock()

	return &proto.OrderingReply{
		Seq: tseq,
		Sig: proto.PartialSig2Proto(ps),
	}, nil
}

func (lbbft *LBBFT) Propose(timeout bool) {
	cmds := lbbft.GetProposeCommands(timeout)
	if (*cmds) == nil {
		return
	}

	pOA := proto.Commands2OrderingArgs(cmds) // todo: 如果是leader，其实可以不用转2次。

	var pOR *proto.OrderingReply
	var err error

	if lbbft.IsLeader {
		pOR, err = lbbft.Ordering(context.TODO(), pOA)
		util.PanicErr(err)
	} else {
		pOR, err = (*lbbft.nodes[config.ReplicaID(lbbft.Leader)]).Ordering(context.TODO(), pOA)
		util.PanicErr(err)

		dPP := &data.PrePrepareArgs{
			View:     lbbft.View,
			Seq:      pOR.Seq,
			Commands: cmds,
		}

		ent := &data.Entry{
			PP: dPP,
		}

		lbbft.PutEntry(ent)
	}

	ent := lbbft.GetEntryBySeq(pOR.Seq)
	ent.Mut.Lock()
	if ent.PreparedCert != nil {
		panic(`collector: ent.PreparedCert != nil`)
	}
	qc := &data.QuorumCert{
		Sigs:       make(map[config.ReplicaID]data.PartialSig),
		SigContent: ent.GetPrepareHash(),
	}
	ent.PreparedCert = qc

	ent.Mut.Unlock()

	go func() {
		leaderPs := *pOR.Sig.Proto2PartialSig()

		if !lbbft.IsLeader { // collector的签名
			ps, err := lbbft.SigCache.CreatePartialSig(lbbft.Config.ID, lbbft.Config.PrivateKey, qc.SigContent.ToSlice())
			if err != nil {
				panic(err)
			}

			// 认证leader的签名
			if !lbbft.SigCache.VerifySignature(leaderPs, qc.SigContent) {
				panic(`ordering signature from leader is not correct`)
			}

			ent.Mut.Lock()
			ent.PreparedCert.Sigs[lbbft.Config.ID] = *ps                     // collector的签名
			ent.PreparedCert.Sigs[config.ReplicaID(lbbft.Leader)] = leaderPs // leader的签名
			ent.Mut.Unlock()
		} else {
			ent.Mut.Lock()
			ent.PreparedCert.Sigs[config.ReplicaID(lbbft.Leader)] = leaderPs
			ent.Mut.Unlock()
		}
	}()

	pPP := &proto.PrePrepareArgs{
		View:     ent.PP.View,
		Seq:      ent.PP.Seq,
		Commands: pOA.Commands,
	}
	lbbft.BroadcastPrePrepareRequest(pPP, ent)
}

func (lbbft *LBBFT) BroadcastPrePrepareRequest(pPP *proto.PrePrepareArgs, ent *data.Entry) {
	logger.Printf("[B/PrePrepare]: view: %d, seq: %d, (%d commands)\n", pPP.View, pPP.Seq, len(pPP.Commands))

	for rid, client := range lbbft.nodes {
		if rid != lbbft.Config.ID && rid != config.ReplicaID(lbbft.Leader) { // 向leader发送的ordering，相当于PrePrepare了，所以不需要重复发送
			go func(id config.ReplicaID, cli *proto.LBBFTClient) {
				pPPR, err := (*cli).PrePrepare(context.TODO(), pPP)
				if err != nil {
					panic(err)
				}
				dPS := pPPR.Sig.Proto2PartialSig()
				ent.Mut.Lock()
				if ent.Prepared == false &&
					lbbft.SigCache.VerifySignature(*dPS, ent.GetPrepareHash()) {
					ent.PreparedCert.Sigs[id] = *dPS
					if len(ent.PreparedCert.Sigs) > int(2*lbbft.F) {

						// collector先处理自己的entry的commit的签名
						ent.Prepared = true
						if ent.CommittedCert != nil {
							panic(`leader: ent.CommittedCert != nil`)
						}

						qc := &data.QuorumCert{
							Sigs:       make(map[config.ReplicaID]data.PartialSig),
							SigContent: ent.GetCommitHash(),
						}
						ent.CommittedCert = qc

						// 收拾收拾，准备broadcast prepare
						pP := &proto.PrepareArgs{
							View: ent.PP.View,
							Seq:  ent.PP.Seq,
							QC:   proto.QuorumCertToProto(ent.PreparedCert),
						}
						ent.Mut.Unlock()

						go func() { // 签名比较耗时，所以用goroutine来进行
							ps, err := lbbft.SigCache.CreatePartialSig(lbbft.Config.ID, lbbft.Config.PrivateKey, qc.SigContent.ToSlice())
							if err != nil {
								panic(err)
							}
							ent.Mut.Lock()
							ent.CommittedCert.Sigs[lbbft.Config.ID] = *ps
							ent.Mut.Unlock()
						}()

						lbbft.BroadcastPrepareRequest(pP, ent)
					} else {
						ent.Mut.Unlock()
					}
				} else {
					ent.Mut.Unlock()
				}
			}(rid, client)
		}
	}
}

func (lbbft *LBBFT) BroadcastPrepareRequest(pP *proto.PrepareArgs, ent *data.Entry) {
	logger.Printf("[B/Prepare]: view: %d, seq: %d\n", pP.View, pP.Seq)

	for rid, client := range lbbft.nodes {

		if rid != lbbft.Config.ID {
			go func(id config.ReplicaID, cli *proto.LBBFTClient) {
				pPR, err := (*cli).Prepare(context.TODO(), pP)
				if err != nil {
					panic(err)
				}
				dPS := pPR.Sig.Proto2PartialSig()
				ent.Mut.Lock()
				if ent.Committed == false &&
					lbbft.SigCache.VerifySignature(*dPS, ent.GetCommitHash()) {
					ent.CommittedCert.Sigs[id] = *dPS
					if len(ent.CommittedCert.Sigs) > int(2*lbbft.F) {

						ent.Committed = true

						// 收拾收拾，准备broadcast commit
						pC := &proto.CommitArgs{
							View: ent.PP.View,
							Seq:  ent.PP.Seq,
							QC:   proto.QuorumCertToProto(ent.CommittedCert),
						}
						lbbft.BroadcastCommitRequest(pC)

						ent.Mut.Unlock()
						go lbbft.ApplyCommands(pC.Seq)
					} else {
						ent.Mut.Unlock()
					}
				} else {
					ent.Mut.Unlock()
				}

			}(rid, client)
		}
	}
}

func (lbbft *LBBFT) BroadcastCommitRequest(pC *proto.CommitArgs) {
	logger.Printf("[B/Commit]: view: %d, seq: %d\n", pC.View, pC.Seq)

	for rid, client := range lbbft.nodes {
		if rid != lbbft.Config.ID {
			go func(id config.ReplicaID, cli *proto.LBBFTClient) {
				_, err := (*cli).Commit(context.TODO(), pC)
				if err != nil {
					panic(err)
				}
			}(rid, client)
		}
	}
}

func (lbbft *LBBFT) PrePrepare(_ context.Context, pPP *proto.PrePrepareArgs) (*proto.PrePrepareReply, error) {

	logger.Printf("PrePrepare: view:%d, seq:%d\n", pPP.View, pPP.Seq)

	dPP := pPP.Proto2PP()

	if !lbbft.Changing && lbbft.View == dPP.View {

		ent := &data.Entry{
			PP: dPP,
		}
		ent.Mut.Lock()
		lbbft.PutEntry(ent)

		ent.PP = dPP
		ps, err := lbbft.SigCache.CreatePartialSig(lbbft.Config.ID, lbbft.Config.PrivateKey, ent.GetPrepareHash().ToSlice())
		util.PanicErr(err)

		ent.Mut.Unlock()

		ppReply := &proto.PrePrepareReply{
			Sig: proto.PartialSig2Proto(ps),
		}
		return ppReply, nil

	} else {
		return nil, errors.New(`正在view change 或者 view不匹配`)
	}
}

func (lbbft *LBBFT) Prepare(_ context.Context, pP *proto.PrepareArgs) (*proto.PrepareReply, error) {
	logger.Printf("Receive Prepare: seq: %d, view: %d\n", pP.Seq, pP.View)

	lbbft.Mut.Lock()
	if !lbbft.Changing && lbbft.View == pP.View {
		ent := lbbft.GetEntryBySeq(pP.Seq)
		lbbft.Mut.Unlock()

		ent.Mut.Lock()

		if ent.Prepared == true {
			panic(`already prepared...`)
		}

		// 检查qc
		dQc := pP.QC.Proto2QuorumCert()
		if !lbbft.SigCache.VerifyQuorumCert(dQc) {
			logger.Println("Prepared QC not verified!: ", dQc)
			return nil, errors.New(`Prepared QC not verified!`)
		}

		if ent.PreparedCert != nil {
			panic(`receiver: ent.PreparedCert != nil`)
		}
		ent.PreparedCert = &data.QuorumCert{
			Sigs:       dQc.Sigs,
			SigContent: dQc.SigContent,
		}
		ent.Prepared = true
		ent.PrepareHash = &dQc.SigContent // 这里应该做检查的，如果先收到PP，PHash需要相等。PP那里，如果有PHash和CHash需要检查是否相等。这里简化了。

		ps, err := lbbft.SigCache.CreatePartialSig(lbbft.Config.ID, lbbft.Config.PrivateKey, ent.GetCommitHash().ToSlice())
		if err != nil {
			panic(err)
		}
		ent.Mut.Unlock()

		pPR := &proto.PrepareReply{Sig: proto.PartialSig2Proto(ps)}
		return pPR, nil

	} else {
		lbbft.Mut.Unlock()
	}
	return nil, nil
}

func (lbbft *LBBFT) Commit(_ context.Context, pC *proto.CommitArgs) (*empty.Empty, error) {
	logger.Printf("Receive Commit: seq: %d, view: %d\n", pC.Seq, pC.View)
	lbbft.Mut.Lock()
	if !lbbft.Changing && lbbft.View == pC.View {
		ent := lbbft.GetEntryBySeq(pC.Seq)
		lbbft.Mut.Unlock()

		ent.Mut.Lock()

		if ent.Committed == true {
			panic(`already committed...`)
		}

		// 检查qc
		dQc := pC.QC.Proto2QuorumCert()
		if !lbbft.SigCache.VerifyQuorumCert(dQc) {
			logger.Println("Commit QC not verified!: ", dQc)
			return &empty.Empty{}, errors.New(`Commit QC not verified!`)
		}

		if ent.CommittedCert != nil {
			panic(`follower: ent.CommittedCert != nil`)
		}
		ent.CommittedCert = &data.QuorumCert{
			Sigs:       dQc.Sigs,
			SigContent: dQc.SigContent,
		}
		ent.Committed = true
		ent.CommitHash = &dQc.SigContent // 这里应该做检查的，如果先收到P，CHash需要相等。PP那里，如果有PHash和CHash需要检查是否相等。这里简化了。

		go lbbft.ApplyCommands(pC.Seq)

		ent.Mut.Unlock()
		return &empty.Empty{}, nil
	} else {
		lbbft.Mut.Unlock()
	}
	return &empty.Empty{}, nil
}

func (lbbft *LBBFT) ApplyCommands(commitSeq uint32) {
	lbbft.Mut.Lock()
	cmds, applyUpdate := lbbft.GetApplyCmds(lbbft.Apply + 1)

	if applyUpdate != lbbft.Apply {
		lbbft.Apply = applyUpdate
		lbbft.Exec <- *cmds
	}

	lbbft.Mut.Unlock()
}

func (lbbft *LBBFT) ApplyCommands2(elem *util.PQElem) {
	lbbft.Mut.Lock()
	inserted := lbbft.ApplyQueue.Insert(*elem)
	if !inserted {
		panic("Already insert some request with same sequence")
	}

	for i, sz := 0, lbbft.ApplyQueue.Length(); i < sz; i++ { // commit需要按global seq的顺序
		m, err := lbbft.ApplyQueue.GetMin()
		if err != nil {
			break
		}
		if int(lbbft.Apply+1) == m.Pri {
			lbbft.Apply++
			cmds, ok := m.C.([]data.Command)
			if ok {
				lbbft.Exec <- cmds
			}
			lbbft.ApplyQueue.ExtractMin()

		} else if int(lbbft.Apply+1) > m.Pri {
			panic("This should already done")
		} else {
			break
		}
	}
	lbbft.Mut.Unlock()
}

func newLBBFTServer(lbbft *LBBFT) *lbbftServer {
	pbftSrv := &lbbftServer{
		LBBFT:   lbbft,
		clients: make(map[context.Context]config.ReplicaID),
	}
	return pbftSrv
}

func (lbbft *lbbftServer) getClientID(ctx context.Context) (config.ReplicaID, error) {
	lbbft.mut.RLock()
	// fast path for known stream
	if id, ok := lbbft.clients[ctx]; ok {
		lbbft.mut.RUnlock()
		return id, nil
	}

	lbbft.mut.RUnlock()
	lbbft.mut.Lock()
	defer lbbft.mut.Unlock()

	// cleanup finished streams
	for ctx := range lbbft.clients {
		if ctx.Err() != nil {
			delete(lbbft.clients, ctx)
		}
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, fmt.Errorf("getClientID: metadata not available")
	}

	v := md.Get("id")
	if len(v) < 1 {
		return 0, fmt.Errorf("getClientID: id field not present")
	}

	id, err := strconv.Atoi(v[0])
	if err != nil {
		return 0, fmt.Errorf("getClientID: cannot parse ID field: %w", err)
	}

	info, ok := lbbft.Config.Replicas[config.ReplicaID(id)]
	if !ok {
		return 0, fmt.Errorf("getClientID: could not find info about id '%d'", id)
	}

	v = md.Get("proof")
	if len(v) < 2 {
		return 0, fmt.Errorf("getClientID: No proof found")
	}

	var R, S big.Int
	v0, err := base64.StdEncoding.DecodeString(v[0])
	if err != nil {
		return 0, fmt.Errorf("getClientID: could not decode proof: %v", err)
	}
	v1, err := base64.StdEncoding.DecodeString(v[1])
	if err != nil {
		return 0, fmt.Errorf("getClientID: could not decode proof: %v", err)
	}
	R.SetBytes(v0)
	S.SetBytes(v1)

	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], uint32(lbbft.Config.ID))
	hash := sha512.Sum512(b[:])

	if !ecdsa.Verify(info.PubKey, hash[:], &R, &S) {
		return 0, fmt.Errorf("Invalid proof")
	}

	lbbft.clients[ctx] = config.ReplicaID(id)
	return config.ReplicaID(id), nil
}
