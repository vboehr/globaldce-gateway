package wire
import (
	//"github.com/globaldce/go-globaldce/applog"
	"net"
)

type Peer struct {
	// 
	Address string
	Connection net.Conn
	SyncingMainchainlength uint32
	GoodIPArray []string
	BadIPArray []string
	BannedNameArray []string
}
func NewPeer(peeraddress string,conn net.Conn) * Peer {
	np:=new(Peer)
	np.Address=peeraddress
	np.Connection=conn
	np.SyncingMainchainlength=uint32(0)	

	return np
}
func(p *Peer) WriteMessage(msg * Message){
	p.WriteTCPMessage(msg)
}
func(p *Peer) WriteTCPMessage(msg * Message){
	msg.WriteBytes(p.Connection)
}
func(p * Peer) ReadMessage() (*Message,error){
		msg,err:= p.ReadTCPMessage()
		
		return msg,err
	
}