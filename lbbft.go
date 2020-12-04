package lbbft

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/joe-zxh/lbbft/data"
	"github.com/joe-zxh/lbbft/util"
	"log"
	"net"
	"sync"
	"time"

	"github.com/joe-zxh/lbbft/config"
	"github.com/joe-zxh/lbbft/consensus"
	"github.com/joe-zxh/lbbft/internal/logging"
	"github.com/joe-zxh/lbbft/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	go func() {
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

			ent.Mut.Lock()
			lbbft.PutEntry(ent)
			ent.Mut.Unlock()
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
	}()
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
				if ent.Prepared == false {
					ent.PreparedCert.Sigs[id] = *dPS
					if len(ent.PreparedCert.Sigs) >= int(lbbft.Q) && lbbft.SigCache.VerifyQuorumCert(ent.PreparedCert) {

						// collector先处理自己的entry的commit的签名
						ent.Prepared = true
						lbbft.UpdateLastPreparedID(ent)

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
				if ent.Committed == false {
					ent.CommittedCert.Sigs[id] = *dPS
					if len(ent.CommittedCert.Sigs) >= int(lbbft.Q) && lbbft.SigCache.VerifyQuorumCert(ent.CommittedCert) {

						ent.Committed = true

						// 收拾收拾，准备broadcast commit
						pC := &proto.CommitArgs{
							View: ent.PP.View,
							Seq:  ent.PP.Seq,
							QC:   proto.QuorumCertToProto(ent.CommittedCert),
						}

						ent.Mut.Unlock()
						go lbbft.BroadcastCommitRequest(pC)
						lbbft.ApplyCommands(pC.Seq)
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
		lbbft.Mut.Unlock()
		ent := lbbft.GetEntryBySeq(pP.Seq)

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

		lbbft.UpdateLastPreparedID(ent)

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
		lbbft.Mut.Unlock()
		ent := lbbft.GetEntryBySeq(pC.Seq)

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

// view change...
func (lbbft *LBBFT) StartViewChange() {
	logger.Printf("StartViewChange: \n")

	lbbft.Mut.Lock()
	lbbft.View += 1

	ent := lbbft.GetEntryBySeq(lbbft.LastPreparedID.N)

	ent.Mut.Lock()
	if ent.PP.View != lbbft.LastPreparedID.V {
		panic(`ent.PP.View!=pbft.LastPreparedID.V`)
	}

	pRV := &proto.RequestVoteArgs{
		CandidateID:      lbbft.ID,
		NewView:          lbbft.View,
		LastPreparedView: ent.PP.View,
		LastPreparedSeq:  ent.PP.Seq,
		Digest:           ent.GetDigest().ToSlice(),
		PreparedCert:     proto.QuorumCertToProto(ent.PreparedCert),
	}
	ent.Mut.Unlock()
	lbbft.Mut.Unlock()

	// 添加candidate自己的签名
	sig, sigcontent := lbbft.CreateRequestVoteReplySig(pRV)

	rvqc := &proto.QuorumCert{
		Sigs:       make([]*proto.PartialSig, 0),
		SigContent: *sigcontent,
	}

	rvqc.Sigs = append(rvqc.Sigs, sig)

	lbbft.BroadcastRequestVote(pRV, rvqc)
}

func (lbbft *LBBFT) BroadcastRequestVote(pRV *proto.RequestVoteArgs, rvqc *proto.QuorumCert) {

	logger.Printf("Broadcast RequestVote: New view: %d, Candidate ID: %d\n", pRV.NewView, pRV.CandidateID)

	broadcastNewView := false

	for rid, client := range lbbft.nodes {
		if rid != lbbft.Config.ID {
			go func(cli *proto.LBBFTClient) {
				pRVReply, err := (*cli).RequestVote(context.TODO(), pRV)
				util.PanicErr(err)
				lbbft.Mut.Lock()
				if !broadcastNewView && pRVReply.Sig != nil {
					rvqc.Sigs = append(rvqc.Sigs, pRVReply.Sig)
				}
				if !broadcastNewView && len(rvqc.Sigs) >= int(lbbft.Q) && lbbft.SigCache.VerifyQuorumCert(rvqc.Proto2QuorumCert()) {
					pNV := &proto.NewViewArgs{
						NewView:     pRV.NewView,
						CandidateID: lbbft.ID,
						NewViewCert: rvqc,
					}
					broadcastNewView = true

					lbbft.View = pRV.NewView
					lbbft.Leader = lbbft.ID
					lbbft.IsLeader = true
					lbbft.VoteFor[pRV.NewView] = lbbft.ID

					lbbft.Mut.Unlock()
					lbbft.ViewChangeChan <- struct{}{}
					lbbft.BroadcastNewView(pNV)

				} else {
					lbbft.Mut.Unlock()
				}
			}(client)
		}
	}
}

func (lbbft *LBBFT) BroadcastNewView(pNV *proto.NewViewArgs) {
	logger.Printf("Broadcast NewView: NewView: %d, CandidateID: %d \n", pNV.NewView, pNV.CandidateID)
	for rid, client := range lbbft.nodes {
		if rid != lbbft.Config.ID {
			go func(cli *proto.LBBFTClient) {
				_, err := (*cli).NewView(context.TODO(), pNV)
				util.PanicErr(err)
			}(client)
		}
	}
}

func (lbbft *LBBFT) RequestVote(_ context.Context, pRV *proto.RequestVoteArgs) (*proto.RequestVoteReply, error) {
	logger.Printf("Receive RequestVote: new view: %d\n", pRV.NewView)

	_, ok := lbbft.VoteFor[pRV.NewView]

	if !ok && pRV.NewView > lbbft.View &&
		lbbft.LastPreparedID.IsOlderOrEqual(&data.EntryID{V: pRV.LastPreparedView, N: pRV.LastPreparedSeq}) &&
		lbbft.IsRequestVotePreparedCertValid(pRV) {

		ps, _ := lbbft.CreateRequestVoteReplySig(pRV)

		return &proto.RequestVoteReply{Sig: ps}, nil
	}

	return &proto.RequestVoteReply{}, nil
}

func (lbbft *LBBFT) NewView(_ context.Context, pNV *proto.NewViewArgs) (*proto.NewViewReply, error) {
	logger.Printf("Receive NewView: new view: %d\n", pNV.NewView)

	_, ok := lbbft.VoteFor[pNV.NewView]

	if !ok && pNV.NewView > lbbft.View &&
		lbbft.IsNewViewCertValid(pNV) {
		lbbft.Mut.Lock()

		lbbft.View = pNV.NewView
		lbbft.Leader = pNV.CandidateID
		lbbft.IsLeader = false

		lbbft.VoteFor[pNV.NewView] = pNV.CandidateID

		lbbft.Mut.Unlock()
		logger.Printf("enter new view: %d\n", pNV.NewView)
		lbbft.ViewChangeChan <- struct{}{}

		return &proto.NewViewReply{
			View: lbbft.LastPreparedID.V,
			Seq:  lbbft.LastPreparedID.N,
		}, nil
	}

	return &proto.NewViewReply{View: ^uint32(0), Seq: ^uint32(0)}, nil // ^uint32(0)表示不接受NewView
}

func newLBBFTServer(lbbft *LBBFT) *lbbftServer {
	pbftSrv := &lbbftServer{
		LBBFT:   lbbft,
		clients: make(map[context.Context]config.ReplicaID),
	}
	return pbftSrv
}
