package net_plugin

import (
	"bytes"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"time"
)

const (
	// default value initializers
	defSendBufferSizeMb        = 4
	defSendBufferSize          = 1024 * 1024 * defSendBufferSizeMb
	defMaxClients              = 25 // 0 for unlimited clients
	defMaxNodesPerHost         = 1
	defConnRetryWait           = 30
	defTxnExpireWait           = time.Duration(3 * time.Second)
	defRespExpectedWait        = time.Duration(5 * time.Second)
	defSyncFetchSpan           = 100
	defMaxJustSend      uint32 = 1500 // roughly 1 "mtu"
	largeMsgNotify      bool   = false

	messageHeaderSize = 4

	/*	   For a while, network version was a 16 bit value equal to the second set of 16 bits
		   of the current build's git commit id. We are now replacing that with an integer protocol
		   identifier. Based on historical analysis of all git commit identifiers, the larges gap
		   between ajacent commit id values is shown below.
		   these numbers were found with the following commands on the master branch:

		   git log | grep "^commit" | awk '{print substr($2,5,4)}' | sort -u > sorted.txt
		   rm -f gap.txt; prev=0; for a in $(cat sorted.txt); do echo $prev $((0x$a - 0x$prev)) $a >> gap.txt; prev=$a; done; sort -k2 -n gap.txt | tail

		   DO NOT EDIT net_version_base OR net_version_range!
	*/
	netVersionBase  uint16 = 0x04b5
	netVersionRange uint16 = 106

	//If there is a change to network protocol or behavior, increment net version to identify
	//the need for compatibility hooks
	protoBase         uint16 = 0
	protoExplicitSync uint16 = 1

	netVersion uint16 = protoExplicitSync
)

/*
   Index by id
   Index by is_known, block_num, validated_time, this is the order we will broadcast
   to peer.
   Index by is_noticed, validated_time
*/
type transactionState struct {
	id              common.TransactionIdType
	isKnownByPeer   bool //true if we sent or received this trx to this peer or received notice from peer
	isNoticedToPeer bool //have we sent peer notice we know it (true if we receive from this peer)
	blockNum        uint32
	expires         common.TimePointSec
	requestedTime   common.TimePoint
}

// typedef multi_index_container<
//    transaction_state,
//    indexed_by<
//       ordered_unique< tag<by_id>, member<transaction_state, transaction_id_type, &transaction_state::id > >,
//       ordered_non_unique< tag< by_expiry >, member< transaction_state,fc::time_point_sec,&transaction_state::expires >>,
//       ordered_non_unique<
//          tag<by_block_num>,
//          member< transaction_state,
//                  uint32_t,
//                  &transaction_state::block_num > >
//       >

//    > transaction_state_index;

type peerBlockState struct {
	id          common.BlockIdType
	blockNum    uint32
	isKnown     bool
	isNoticed   bool
	requestTime common.TimePoint
}

// typedef multi_index_container<
//    eosio::peer_block_state,
//    indexed_by<
//       ordered_unique< tag<by_id>, member<eosio::peer_block_state, block_id_type, &eosio::peer_block_state::id > >,
//       ordered_unique< tag<by_block_num>, member<eosio::peer_block_state, uint32_t, &eosio::peer_block_state::block_num > >
//       >
//    > peer_block_state_index;

type nodeTransactionState struct {
	id            common.TransactionIdType
	expires       common.TimePointSec //time after which this may be purged.Expires increased while the txn is "in flight" to another peer
	packedTxn     types.PackedTransaction
	serializedTxn common.HexBytes // the received raw bundle
	blockNum      uint32          // block transaction was included in
	trueBlock     uint32          // used to reset block_uum when request is 0
	requests      uint16          // the number of "in flight" requests for this txn
}

// typedef multi_index_container<
//    node_transaction_state,
//    indexed_by<
//       ordered_unique<
//          tag< by_id >,
//          member < node_transaction_state,
//                   transaction_id_type,
//                   &node_transaction_state::id > >,
//       ordered_non_unique<
//          tag< by_expiry >,
//          member< node_transaction_state,
//                  fc::time_point_sec,
//                  &node_transaction_state::expires >
//          >,
//       ordered_non_unique<
//          tag<by_block_num>,
//          member< node_transaction_state,
//                  uint32_t,
//                  &node_transaction_state::block_num > >
//       >
//    >
// node_transaction_index;

func (n *nodeTransactionState) ElementObject() {}

func (n *transactionState) ElementObject() {}

func (n *peerBlockState) ElementObject() {}

func CompareById(first common.ElementObject, second common.ElementObject) int {
	result := -2
	switch first.(type) {
	case *nodeTransactionState:
		fir := first.(*nodeTransactionState)
		sec := second.(*nodeTransactionState)
		result = bytes.Compare(fir.id.Bytes(), sec.id.Bytes())
	case *transactionState:
		fir := first.(*transactionState)
		sec := second.(*transactionState)
		result = bytes.Compare(fir.id.Bytes(), sec.id.Bytes())
	case *peerBlockState:
		fir := first.(*transactionState)
		sec := second.(*transactionState)
		result = bytes.Compare(fir.id.Bytes(), sec.id.Bytes())
	}
	return result
}

func CompareByBlockNum(first common.ElementObject, second common.ElementObject) int {
	result := -2
	switch first.(type) {
	case *nodeTransactionState:
		fir := first.(*nodeTransactionState)
		sec := second.(*nodeTransactionState)
		if fir.blockNum == sec.blockNum {
			result = 0
		} else if fir.blockNum < sec.blockNum {
			result = -1
		} else {
			result = 1
		}
	case *transactionState:
		fir := first.(*transactionState)
		sec := second.(*transactionState)
		if fir.blockNum == sec.blockNum {
			result = 0
		} else if fir.blockNum < sec.blockNum {
			result = -1
		} else {
			result = 1
		}
	case *peerBlockState:
		fir := first.(*transactionState)
		sec := second.(*transactionState)
		if fir.blockNum == sec.blockNum {
			result = 0
		} else if fir.blockNum < sec.blockNum {
			result = -1
		} else {
			result = 1
		}
	}
	return result
}

func CompareByExpiry(first common.ElementObject, second common.ElementObject) int {
	fir := first.(*transactionState)
	sec := second.(*transactionState)
	if fir.expires == sec.expires {
		return 0
	} else if fir.expires < sec.expires {
		return -1
	} else {
		return 1
	}
}

//Index by start_block_num
type syncState struct {
	startBlock uint32
	endBlock   uint32
	last       uint32           //last sent or received
	startTime  common.TimePoint //time request made or received
}

func newSyncState(start, end, lastActed uint32) *syncState {
	return &syncState{
		startBlock: start,
		endBlock:   end,
		last:       lastActed,
		startTime:  common.Now(),
	}
}

type stages byte

const (
	libCatchup = stages(iota)
	headCatchup
	inSync
)

type syncManager struct {
	syncKnownLibNum      uint32
	syncLastRequestedNum uint32
	syncNextExpectedNum  uint32
	syncReqSpan          uint32
	source               *Peer
	state                stages
	_blocks              common.BlockIdType //<deque<block_id_type>>
	//chainPlugin *chainPlugin
	myImpl *netPluginIMpl
}

func NewSyncManager(impl *netPluginIMpl, span uint32) *syncManager {
	//chainPlugin :=
	return &syncManager{
		syncKnownLibNum:      0,
		syncLastRequestedNum: 0,
		syncNextExpectedNum:  1,
		syncReqSpan:          span,
		//source:
		state:  inSync,
		myImpl: impl,
	}
}

func stageStr(s stages) string {
	switch s {
	case libCatchup:
		return "lib catchup"
	case headCatchup:
		return "head catchup"
	case inSync:
		return "in sync"
	default:
		return "unkown"
	}
}

func (s *syncManager) setStage(newstate stages) {
	if s.state == newstate {
		return
	}
	s.myImpl.log.Info("old state %s becoming %s", stageStr(s.state), stageStr(newstate))
	s.state = newstate
}

func (s *syncManager) syncRequired() bool {
	s.myImpl.log.Info("last req = %d,last recv = %d known = %d our head %d\n", +s.syncLastRequestedNum, s.syncNextExpectedNum, s.syncKnownLibNum, 100) //chain_plug->chain( ).head_block_num( )
	return s.syncLastRequestedNum < s.syncKnownLibNum || 0 < s.syncLastRequestedNum                                                                    //100  ---->  chain_plug->chain( ).head_block_num( )
}

func (s *syncManager) isActive(p *Peer) bool {
	if s.state == headCatchup && p != nil { //TODO
		fhset := p.forkHead != common.BlockIdType(*crypto.NewSha256Nil())
		s.myImpl.log.Info("fork_head_num = %d fork_head set = %s\n", p.forkHeadNum, fhset)

		return p.forkHead != common.BlockIdNil()
		//&& p.forkHeadNum < chain_plug->chain().head_block_num()

		//         return c->fork_head != block_id_type() && c->fork_head_num < chain_plug->chain().head_block_num();
	}
	return s.state != inSync
}

func (s *syncManager) resetLibNum(myImpl *netPluginIMpl, p *Peer) {
	if s.state == inSync {
		s.source.reset()
	}
	if p.current() {
		if p.lastHandshakeRecv.LastIrreversibleBlockNum > s.syncKnownLibNum {
			s.syncKnownLibNum = p.lastHandshakeRecv.LastIrreversibleBlockNum
		}
	} else if p == s.source {
		s.syncLastRequestedNum = 0
		s.requestNextChunk(myImpl, p)
	}
}

func (s *syncManager) requestNextChunk(myImpl *netPluginIMpl, p *Peer) {
	//syncRequest := SyncRequestMessage{
	//	StartBlock: s.syncLastRequestedNum,
	//	EndBlock:   s.syncNextExpectedNum,
	//}
	//p.write(&syncRequest)
	//
	////uint32_t head_block = chain_plug->chain().fork_db_head_block_num();
	//headBlock := uint32(100)
	//if headBlock < s.syncLastRequestedNum && s.source != nil && s.source.current() {
	//	//fc_ilog (logger, "ignoring request, head is ${h} last req = ${r} source is ${p}",
	//	// ("h",head_block)("r",sync_last_requested_num)("p",source->peer_name()));
	//	netlog.Info("ignoring request,head is %d last req = %d source is %s\n", +headBlock, s.syncLastRequestedNum, p.peerAddr)
	//	return
	//}
	///* ----------
	// * next chunk provider selection criteria
	// * a provider is supplied and able to be used, use it.
	// * otherwise select the next available from the list, round-robin style.
	// */
	//
	//if p != nil && p.current() {
	//	s.source = p
	//} else {
	//	if len(myImpl.peers) == 1 {
	//		if s.source == nil {
	//			for _, p := range myImpl.peers {
	//				s.source = p
	//			}
	//		}
	//	} else {
	//		//// init to a linear array search
	//		//auto cptr = my_impl->connections.begin();
	//		//auto cend = my_impl->connections.end();
	//		//// do we remember the previous source?
	//		//if (source) {
	//		// //try to find it in the list
	//		// cptr = my_impl->connections.find(source);
	//		// cend = cptr;
	//		// if (cptr == my_impl->connections.end()) {
	//		//	 //not there - must have been closed! cend is now connections.end, so just flatten the ring.
	//		//	 source.reset();
	//		//	 cptr = my_impl->connections.begin();
	//		// } else {
	//		//	 //was found - advance the start to the next. cend is the old source.
	//		//	 if (++cptr == my_impl->connections.end() && cend != my_impl->connections.end() ) {
	//		//	 cptr = my_impl->connections.begin();
	//		//	 }
	//		// }
	//		//}
	//		//scan the list of peers looking for another able to provide sync blocks.
	//		// auto cstart_it = cptr;
	//		// do {
	//		//	 //select the first one which is current and break out.
	//		//	 if((*cptr)->current()) {
	//		//	 source = *cptr;
	//		//	 break;
	//		// }
	//		//	 if(++cptr == my_impl->connections.end())
	//		//	 cptr = my_impl->connections.begin();
	//		// } while(cptr != cstart_it);
	//		// // no need to check the result, either source advanced or the whole list was checked and the old source is reused.
	//		//}
	//
	//	}
	//
	//	// verify there is an available source
	//	if s.source != nil || !s.source.current() {
	//		//elog("Unable to continue syncing at this time")
	//		netlog.Error("Unable to continue syncing at this time")
	//		//sync_known_lib_num = chain_plug->chain().last_irreversible_block_num()
	//		s.syncLastRequestedNum = 0
	//		s.setStage(inSync) // probably not, but we can't do anything else
	//		return
	//	}

	if s.syncLastRequestedNum != s.syncKnownLibNum {
		start := s.syncNextExpectedNum
		end := start + s.syncReqSpan - 1
		if end > s.syncKnownLibNum {
			end = s.syncKnownLibNum
		}
		if end > 0 && end >= start {
			//fc_ilog(logger, "requesting range ${s} to ${e}, from ${n}",
			//	("n",source->peer_name())("s",start)("e",end));
			s.myImpl.log.Info("requesting range %s to %d, from %d\n", s.source.peerAddr, start, end)
			s.source.requestSyncBlocks(start, end)
			p.requestSyncBlocks(start, end)
			s.syncLastRequestedNum = end

		}
	}

	//}
}

func (s *syncManager) sendHandshakes(impl *netPluginIMpl) {
	for _, p := range impl.peers {
		if p.current() {
			p.sendHandshake(impl)
		}
	}
}

func (s *syncManager) recvHandshake(myImpl *netPluginIMpl, p *Peer, msg *HandshakeMessage) {
	//controller& cc = chain_plug->chain();
	//libNum := cc.last_irreversible_block_num()
	//libNum := uint32(100) //TODO
	libNum := uint32(0) //TODO
	peerLib := msg.LastIrreversibleBlockNum
	s.resetLibNum(myImpl, p)
	p.syncing = false

	//--------------------------------
	// sync need checks; (lib == last irreversible block)
	//
	// 0. my head block id == peer head id means we are all caugnt up block wise
	// 1. my head block num < peer lib - start sync locally
	// 2. my lib > peer head num - send an last_irr_catch_up notice if not the first generation
	//
	// 3  my head block num <= peer head block num - update sync state and send a catchup request
	// 4  my head block num > peer block num ssend a notice catchup if this is not the first generation
	//
	//-----------------------------

	//head := cc.headBlockNum()
	//headID :=cc.headBlockID()

	head := uint32(100) //TODO
	headID := common.BlockIdType(*crypto.NewSha256Nil())

	if headID == msg.HeadID {
		s.myImpl.log.Info("sync check statue 0")
		// notify peer of our pending transactions

		note := NoticeMessage{}
		note.KnownBlocks.Mode = none
		note.KnownTrx.Mode = catchUp
		//note.KnownTrx.Pending = my_impl->local_txns.size()//TODO
		note.KnownBlocks.Pending = uint32(len(myImpl.localTxns.indexs)) //TODO
		p.write(&note)
		return

	}

	if head < peerLib {
		s.myImpl.log.Info("sync check state 1")
		//wait for receipt of a notice message before initiating sync
		//if p.protocolVersion < protoExplicitSync {
		s.startSync(myImpl, p, peerLib)
		//}
		return
	}

	if libNum > msg.HeadNum {
		s.myImpl.log.Info("sync check state 2")
		if msg.Generation > 1 || p.protocolVersion > protoBase {
			note := NoticeMessage{}
			note.KnownBlocks.Mode = lastIrrCatchUp
			note.KnownBlocks.Pending = head
			note.KnownTrx.Mode = lastIrrCatchUp
			note.KnownTrx.Pending = libNum
			p.write(&note)
		}
		p.syncing = true
		return
	}

	if head <= msg.HeadNum {
		s.myImpl.log.Info("sync check state 3")
		s.verifyCatchup(myImpl, p, msg.HeadNum, msg.HeadID)
		return
	} else {
		s.myImpl.log.Info("sync check state 4")
		if msg.Generation > 1 || p.protocolVersion > protoBase {
			note := NoticeMessage{}
			note.KnownBlocks.Mode = catchUp
			note.KnownBlocks.Pending = head
			note.KnownBlocks.IDs = append(note.KnownBlocks.IDs, &headID)
			note.KnownTrx.Mode = none
			p.write(&note)
		}
		p.syncing = true
		return
	}

	s.myImpl.log.Error("sync check failed to resolve status")
}

func (s *syncManager) startSync(myImpl *netPluginIMpl, p *Peer, target uint32) {
	if target > s.syncKnownLibNum {
		s.syncKnownLibNum = target
	}
	//if !s.syncRequired() {
	//	bnum := 100 //chain_plug->chain().last_irreversible_block_num()
	//	hnum := 100 //chain_plug->chain().head_block_num()
	//	netlog.Info("we are already caught up, my irr = %d,head =%d,target = %d\n", bnum, hnum, target)
	//	return
	//}
	if s.state == inSync {
		s.setStage(libCatchup)
		//s.syncNextExpectedNum = 99 + 1 //TODO  chain_plug->chain().last_irreversible_block_num() + 1
		s.syncNextExpectedNum = p.lastHandshakeSent.HeadNum + 1
	}
	s.myImpl.log.Warn("Catching up with chain, our last req is %d, theirs is %d peer %s", +s.syncLastRequestedNum, target, p.peerAddr)

	s.requestNextChunk(myImpl, p)
}
func (s *syncManager) reassignFetch(myImpl *netPluginIMpl, p *Peer, reason GoAwayReason) {
	s.myImpl.log.Info("reassign_fetch, our last req is %d, next expected is %d peer %s\n", +s.syncLastRequestedNum, s.syncNextExpectedNum, p.peerAddr)
	if p == s.source {
		p.cancelSync(reason)
		s.syncLastRequestedNum = 0
		s.requestNextChunk(myImpl, p)
	}
}

func (s *syncManager) verifyCatchup(myImpl *netPluginIMpl, p *Peer, num uint32, id common.BlockIdType) {
	req := RequestMessage{}
	req.ReqBlocks.Mode = catchUp

	for _, peer := range myImpl.peers {
		if peer.forkHead == id || peer.forkHeadNum > num {
			req.ReqBlocks.Mode = none
		}
		break
	}

	if req.ReqBlocks.Mode == catchUp {
		p.forkHead = id
		p.forkHeadNum = num
		s.myImpl.log.Info("got a catch_up notice while in %s, fork head num = %d target LIB = %d next_expected = %d",
			stageStr(s.state), num, s.syncKnownLibNum, s.syncNextExpectedNum)
		if s.state == libCatchup {
			return
		}
		s.setStage(headCatchup)

	} else {
		p.forkHead = common.BlockIdNil()
		p.forkHeadNum = 0
	}

	req.ReqTrx.Mode = none
	p.write(&req)
}

func (s *syncManager) recvNotice(myImpl *netPluginIMpl, p *Peer, msg *NoticeMessage) {
	s.myImpl.log.Info("sync_manager got %s block notice", modeTostring[msg.KnownBlocks.Mode])
	if msg.KnownBlocks.Mode == catchUp {
		IDsCount := len(msg.KnownBlocks.IDs)
		if IDsCount == 0 {
			s.myImpl.log.Error("got a catch up with ids size = 0")
		} else {
			s.verifyCatchup(myImpl, p, msg.KnownBlocks.Pending, *msg.KnownBlocks.IDs[IDsCount-1])
		}
	} else {
		p.lastHandshakeRecv.LastIrreversibleBlockNum = msg.KnownTrx.Pending
		s.resetLibNum(myImpl, p)
		s.startSync(myImpl, p, msg.KnownBlocks.Pending)
	}
}

func (s *syncManager) rejectedBlock(myImpl *netPluginIMpl, p *Peer, blkNum uint32) {
	if s.state != inSync {
		s.myImpl.log.Info("block %d not accepted from %s", blkNum, p.peerAddr)
		s.syncLastRequestedNum = 0
		s.source.reset()
		myImpl.close(p)
		s.setStage(inSync)
		s.sendHandshakes(myImpl)
	}
}

func (s *syncManager) recvBlock(myImpl *netPluginIMpl, p *Peer, blkID common.BlockIdType, blkNum uint32) { //TODO impl

	s.myImpl.log.Info("got block %d from %s", blkNum, p.peerAddr)
	if s.state == libCatchup {
		if blkNum != s.syncNextExpectedNum { //TODO ??
			s.myImpl.log.Info("expected block %d but got %d", s.syncNextExpectedNum, blkNum)
			//myImpl.close(p)
			return
		}
		s.syncNextExpectedNum = blkNum + 1
	}

	if s.state == headCatchup {
		s.myImpl.log.Info("sync_manager in head_catchup state")
		s.setStage(inSync)
		s.source.reset()

		nullID := common.BlockIdType(*crypto.NewSha256Nil())
		for _, p := range myImpl.peers {
			if p.forkHead == nullID {
				continue
			}
			if p.forkHead == blkID || p.forkHeadNum < blkNum {
				p.forkHead = nullID
				p.forkHeadNum = 0
			} else {
				s.setStage(headCatchup)
			}
		}
	} else if s.state == libCatchup {
		if blkNum == s.syncKnownLibNum {
			s.myImpl.log.Info("All caught up with last known last irreversible block resending handshake")
			s.setStage(inSync)
			s.sendHandshakes(myImpl)
		} else if blkNum == s.syncLastRequestedNum {
			s.requestNextChunk(myImpl, p) //TODO        request_next_chunk();
		} else {
			s.myImpl.log.Info("calling sync_wait on connecting %s", p.peerAddr)
			p.syncWait()
		}
	}
}
