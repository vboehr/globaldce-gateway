package wire

import (
	"github.com/globaldce/globaldce-toolbox/applog"
	"github.com/globaldce/globaldce-toolbox/mainchain"
	"github.com/globaldce/globaldce-toolbox/utility"
)
func (sw *Swarm) HandlePeerMessage(mn * mainchain.Maincore,rmsg *  Message) bool{
	applog.Trace("\n new message to be handled",rmsg)
	switch{
	/*case (rmsg.CheckIdentifier( MsgIdentifierReplyMainchainLength)):
		//msg:= NewMessage( MsgIdentifierReplyMainchainLength)
		//msg.PutContent(mn.GetSerializedMainchainLength())
		//*rmsg.OriginPeer.WriteMessage(msg)
		
		op:=*rmsg.OriginPeer
		//op.WriteMessage(msg)
		//msg.WriteBytes(*rmsg.Connection)
		tmpbr:=utility.NewBufferReader(rmsg.GetContent())

		op.SyncingMainchainlength=tmpbr.GetUint32()
		sw.Peers[op.Address]=op
		applog.Trace("\n GOT MsgIdentifierReplyMainchainLength %d",op.SyncingMainchainlength)
		return true
	*/
	case (rmsg.CheckIdentifier( MsgIdentifierRequestMainchainLength)):
		applog.Trace("\n Sending MsgIdentifierReplyMainchainLength")
		msg:= NewMessage( MsgIdentifierReplyMainchainLength)
		msg.PutContent(mn.GetSerializedMainchainLength())
		//*rmsg.OriginPeer.WriteMessage(msg)
		op:=*rmsg.OriginPeer
		op.WriteMessage(msg)
		//msg.WriteBytes(*rmsg.Connection)
		return true
	//////////////////////////////////
	case (rmsg.CheckIdentifier(MsgIdentifierRequestMainheaders)):
		applog.Trace("\nVERY Good headers request")
		correctness,first,last:=DecodeRequestMainheaders(rmsg)
		if correctness && first>0 && last <mn.GetMainchainLength() {
			applog.Trace("first %d last %d",first,last)
			msg:=NewMessage(MsgIdentifierReplyMainheaders)
			msg.PutContent(mn.GetSerializedMainheaders(first,last))
			//msg.WriteBytes(peer.Connection)
			op:=*rmsg.OriginPeer
			op.WriteMessage(msg)
			return true
		} else{
			applog.Warning("incorrect request for mainheader - also first %d last %d",first,last)
			return false
		}

	case (rmsg.CheckIdentifier(MsgIdentifierRequestMainblockTransactions)):
		applog.Trace("\nVERY Good mainblock transactions request")
		correctness,requestedblockheight:=DecodeRequestMainblockTransactions(rmsg)
		if correctness{
			applog.Trace("sending mainblock transactions request %d  ",requestedblockheight)
			msg:=NewMessage(MsgIdentifierReplyMainblockTransactions)
			msg.PutContent(mn.GetSerializedMainblockTransactions(requestedblockheight))
			//msg.WriteBytes(peer.Connection)
			op:=*rmsg.OriginPeer
			op.WriteMessage(msg)
			return true
		} else {
			applog.Warning("incorrect request for mainblocktransactions")
			return false
		}

	//////////////////////////////////
	case (rmsg.CheckIdentifier( MsgIdentifierBroadcastMainblock)):
		applog.Trace("VERY Good mainblock broadcast")
		correctness,nbhops,height,mb:= DecodeBroadcastMainblock(rmsg)
		if correctness{
			applog.Trace("mainblock received nbhops %d height %d mainblock %x",nbhops,height,mb)
			if mn.ValidatePropagatingMainblock(mb){
				applog.Notice("Received valid propagating mainblock")
				mn.AddInMemoryBlock(mb)
				mn.ConfirmBlocks()
				applog.Trace("Mainchainlength %d Confirmedmainchainlength %d",mn.GetMainchainLength(),mn.GetConfirmedMainchainLength())
				//relaying block
				nbhops++
				relayedmsg:=EncodeBroadcastMainblock(nbhops,mb) //*rmsg.OriginPeer
				sw.RelayMessage(relayedmsg,rmsg.OriginPeer)
				//if propagating mainblock is valide increase credibility
				return true
			} else {
				applog.Warning("incorrect mainblock broadcast")
				return false
			}
			//URGENT TODO ban peer that relyed invalide propagating mainblock
			// if mainblock already existes and decrease credibility of peer
		}
	//////////////////////////////////
	case (rmsg.CheckIdentifier( MsgIdentifierBroadcastTransaction)):
		applog.Trace("VERY Good transaction broadcast")
		correctness,nbhops,seize,tx :=DecodeBroadcastTransaction(rmsg)
		//applog.Trace("************ seize of tx %d",seize)
		if correctness{
			validity,fee:= mn.ValidateTransaction(tx)
			//URGENT TODO take into account the weight of the transaction
			if validity {

				priority:=fee/uint64(seize)/1000000
				applog.Trace("************ fee %d priority of tx %d nbhops %d",fee,priority,nbhops)
				//applog.Trace("valid received nbhops %d height %d mainblock %x",nbhops,height,mb)
				mn.AddTransactionToTxsPool(tx,fee,priority)
				//relaying transaction
				nbhops++
				relayedmsg:=EncodeBroadcastTransaction(nbhops,tx) //
				sw.RelayMessage(relayedmsg,rmsg.OriginPeer)
				return true
			} else {
				applog.Warning("invalid transaction in transaction broadcast")
				return false
			}
			
		} else {
			applog.Warning("incorrect transaction broadcast")
			return false
		}
	//////////////////////////////////
	case (rmsg.CheckIdentifier( MsgIdentifierRequestData)):
		applog.Trace("VERY Good request data")
		correctness,phash :=DecodeRequestData(rmsg)
		
		//
		if correctness && phash!=nil {
			applog.Trace("Requested data hash %x",(*phash))
			data,err:=mn.GetData(*phash)
			//
			applog.Trace("Data sent %s",data)
			if err==nil {
				replayedmsg:=EncodeReplyData(data) //
				sw.ReplyMessage(replayedmsg,rmsg.OriginPeer)
				return true
			} else {
				applog.Warning("incorrect request data hash - unkown data hash ")
				return false
			}	
		} else {
			applog.Warning("incorrect request data")
			return false
		}
		//////////////////////////////////
	case (rmsg.CheckIdentifier( MsgIdentifierRequestDataFile)):
		applog.Trace("VERY Good request data file")
		correctness,phash :=DecodeRequestDataFile(rmsg)
		
		//
		if correctness && phash!=nil {
			applog.Trace("Requested data file hash %x",(*phash))
			data,err:=mn.GetDataFile(*phash)
			//
			//applog.Trace("Data file sent %s",data)
			if err==nil {
				replayedmsg:=EncodeReplyDataFile(data) //
				sw.ReplyMessage(replayedmsg,rmsg.OriginPeer)
				return true
			} else {
				applog.Warning("incorrect request data file hash - unkown data file hash ")
				return false
			}	
		} else {
			applog.Warning("incorrect request data file")
			return false
		}
	//////////////////////////////////
	case (rmsg.CheckIdentifier( MsgIdentifierReplyData )):
		applog.Trace("VERY Good reply data")
		correctness,data :=DecodeReplyData(rmsg)
		//
		if correctness {
			hash:=utility.ComputeHash(data)
			if mn.IsMissingData(hash){
				applog.Trace("Adding data")
				mn.AddData(hash,data)
				return true
			} else {

				applog.Trace("incorrect data reply - data is not missing")
				return false
				//TODO decrease reputation of peer
			}
			
		} else {
			applog.Trace("incorrect data reply")
			return false
			//TODO decrease reputation of peer
		}
	//////////////////////////////////
	case (rmsg.CheckIdentifier( MsgIdentifierReplyDataFile )):
		applog.Trace("VERY Good reply data file")
		correctness,data :=DecodeReplyDataFile(rmsg)
		//
		if correctness {
			hash:=utility.ComputeHash(data)
			if mn.IsMissingDataFile(hash){
				applog.Trace("Adding data")
				mn.AddDataFile(hash,data)
				return true
			} else {

				applog.Trace("incorrect data file reply - data file is not missing")
				return false
				//TODO decrease reputation of peer
			}
			
		} else {
			applog.Trace("incorrect data file reply")
			return false
			//TODO decrease reputation of peer
		}
	//////////////////////////////////
	default:
		applog.Warning("unkown indentifier")
		return false
	}
	applog.Warning("unexpected flow of instructions in HandlePeerMessage")
	return false
}