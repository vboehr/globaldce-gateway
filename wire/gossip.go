package wire
import (
	"github.com/globaldce/globaldce/applog"
	"github.com/globaldce/globaldce/mainchain"
	"github.com/globaldce/globaldce/utility"
)


func (sw *Swarm) BroadcastMainblock(mb *mainchain.Mainblock) {
	//applog.Trace("our Mainchainlength is %d", currentmainchainlength)
	applog.Notice("Broadcasting Mainblock %d",mb.Height)

	//blockmsg:=NewMessage(MsgIdentifierBroadcastMainblock)
	blockmsg :=EncodeBroadcastMainblock(uint32(0),mb)
	//blockmsg.WriteBytes(p.Connection)
	//
	for address, p := range sw.Peers {
		p.WriteMessage(blockmsg)
		applog.Trace("Writing to peer %v",address)
	}
}
func (sw *Swarm) BroadcastTransaction(tx *utility.Transaction) {
	//applog.Trace("our Mainchainlength is %d", currentmainchainlength)
	applog.Notice("Broadcasting Transaction !")
	//blockmsg:=NewMessage(MsgIdentifierBroadcastMainblock)
	txmsg :=EncodeBroadcastTransaction(uint32(0),tx)
	//blockmsg.WriteBytes(p.Connection)
	//
	for address, p := range sw.Peers {

		p.WriteMessage(txmsg)
		applog.Trace("Writing to peer %v",address)
	}
}

func (sw *Swarm) RelayMessage(msg *Message,originpeer *Peer) {

	applog.Trace("Relaying message ")

	for address, p := range sw.Peers {

		if address!=originpeer.Address{
			p.WriteMessage(msg)
		}
	}
}
