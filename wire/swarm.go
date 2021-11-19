package wire
import (
	"github.com/globaldce/globaldce-toolbox/applog"
	"net"
	"sync"
	//"github.com/globaldce/globaldce-toolbox/utility"
	
)

type Swarm struct {
	Syncingdone bool
	NewpeersChan chan * Peer
	IpaddrChan chan  string
	PeersmsgChan chan * Message
	Listener net.Listener
	Peers map[string] Peer
	AddrState map[string] string
	mu sync.RWMutex
}

func NewSwarm() *Swarm {
	ns:=new(Swarm)
	ns.Syncingdone=false
	ns.NewpeersChan = make(chan * Peer)
	ns.IpaddrChan = make(chan  string)
	ns.PeersmsgChan = make(chan * Message)//TODO rename PeersmsgChan to MessageChan
	ns.Peers= make(map[string]Peer) 
	ns.AddrState= make(map[string]string) 
	return ns
}

func (sw *Swarm) AddPeer(newpeer * Peer) {
	if sw.CheckPeerAlreadyExist(newpeer.Address){
		return
	}
	
	applog.Trace("Adding peer with address %s",newpeer.Address)
	sw.Peers[newpeer.Address]=*newpeer
}
func (sw *Swarm) RemovePeer(peer * Peer) {
	peer.Connection.Close()
	delete(sw.Peers,peer.Address)
}
func (sw *Swarm) RemovePeerByAddress(addr string) {
	peer:=sw.Peers[addr]
	peer.Connection.Close()
	delete(sw.Peers,peer.Address)
}
func (sw *Swarm) CheckPeerAlreadyExist(peeraddr string) bool {
	if _,ok:=sw.Peers[peeraddr];ok {
		applog.Trace("Peer address %s already exist",peeraddr)
		return true
	}
	return false
}
/*
func (sw *Swarm) Range() {
    for k, v := range sw.Peers {
        //fmt.Println("k:", k, "v  Add:", v.Address)
    }

}
*/
func (sw *Swarm) NbPeers() int {
	return len(sw.Peers)
}
func (sw *Swarm) GetListeningPort() int {
	return sw.Listener.Addr().(*net.TCPAddr).Port
}
func (sw *Swarm) GetLocalIP() string {
	
	return sw.Listener.Addr().String()
}

