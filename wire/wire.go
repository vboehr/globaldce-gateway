package wire
import (
	"github.com/globaldce/globaldce/applog"
	//"fmt"
	"net"
	"io"
	"time"
	"log"
)


func (sw *Swarm) SetupListener() ( err error) {
	//newpeersChan := make(chan * wire.Peer)
	//ipaddrChan := make(chan  string)


	// port dynamic
	//var localaddress string
	//if localport!=""{
	//	localaddress="127.0.0.1:"+localport
	//} else {
	//	localaddress=":0"
	//}
	
	ls, err := net.Listen("tcp", ":0")//TODO allow user to choose a local port
    if err != nil {
		//log.Fatal(err)
		panic(err)
	}
	sw.Listener=ls
	go sw.ListenConnections()	
    go sw.StartMDNSServer()//Wireswarm.GetLocalIP())
    
    
    go sw.StartMDNSClient()
	
	//applog.Notice("\nAccepting swarm connection on port: %d\n", sw.Listener.Addr().(*net.TCPAddr).Port)
	applog.Notice("Accepting swarm connections on port: %d - listening local address: %s" ,sw.Listener.Addr().(*net.TCPAddr).Port, sw.Listener.Addr().String())

	return nil
}

func (sw *Swarm) ListenConnections()  {

	for {
		conn, err := sw.Listener.Accept()
		if err != nil {
			//
			applog.Fatal("Listener accept failed:%v", err)
			//
		} else {
			applog.Trace("Calling handleConnection")
			sw.handleConnection(conn)
		}

	}

}

func (sw *Swarm) handleConnection(conn net.Conn) {
	//defer conn.Close()
	// set SetReadDeadline
	err := conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Println("SetReadDeadline failed:", err)
		// 
		conn.Close()
		return
	}
	rmsg,rerr:=ReadConnectionMessage(conn)
	if rerr != nil {
		log.Println("read handshake request failed:", rerr)
		// 
		conn.Close()
		return
	}
	if (!rmsg.CheckIdentifier(MsgIdentifierRequestHandshake)){
		applog.Trace("Peer rejected")
		conn.Close()
		return
	} else {
		applog.Trace("Peer accepted")
		msghandshake:=NewMessage(MsgIdentifierReplyHandshake)
		msghandshake.WriteBytes(conn)
		sw.NewpeersChan <- &Peer{
			Address: conn.RemoteAddr().String(),
			Connection: conn,
			SyncingMainchainlength:0,
		}
		applog.Trace("new peer address %s",conn.RemoteAddr().String())
	}
	//applog.Trace("Connection was closed")
}

func (sw *Swarm) HintNewPeer(addrstring string){
	applog.Trace("HintNewPeer %s",addrstring)
	if sw.CheckPeerAlreadyExist(addrstring){
		
		//return nil,fmt.Errorf("Peer address %s already exist",addrstring)
	}
	//applog.Trace("++++ %s",addrstring)

	addrstate,addrstateok:=sw.AddrState[addrstring]

	if addrstateok{
		return
		applog.Trace("Peer address state %s",addrstate)
		//return nil,fmt.Errorf("Peer address state %s",addrstate)
	} else {
		//applog.Trace("******** %s",addrstring)
		
		sw.AddrState[addrstring]="CONNECTING"
		
	}
    //addr, _ := net.ResolveTCPAddr("tcp", addrstring)
    conn, err := net.Dial("tcp",  addrstring)
    if err != nil {
		//panic(err.Error())
		//return nil,err
		return
	}
	//defer conn.Close()
	// set SetReadDeadline
	derr := conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if derr != nil {
		log.Println("SetReadDeadline failed:", derr)
		// do something else, for example create new conn
		return
		//return nil,derr
	}
	//////////////////////////////////////
	msghandshake:=NewMessage(MsgIdentifierRequestHandshake)
	
	msghandshake.WriteBytes(conn)

	for {

		rmsg,rerr:=ReadConnectionMessage(conn)
		if (rerr==nil)&&(rmsg!=nil){
			//applog.Trace("ReadConnectionMessage %v",rmsg)
			
			if (!rmsg.CheckIdentifier(MsgIdentifierReplyHandshake)){
				applog.Trace("Peer rejected")
				conn.Close()
				return
				//return nil,rerr
			} else {
				applog.Trace("Peer accepted")
				conn.SetReadDeadline(time.Time{})
				newpeer:=NewPeer(conn.RemoteAddr().String(),conn)
				sw.AddPeer(newpeer)
                if sw.Syncingdone{
                    go sw.ListenPeerMessages(newpeer) 
                } 
				return
				//return newpeer,nil
			}
		}
				

		if rerr != nil {
			if netErr, ok := rerr.(net.Error); ok && netErr.Timeout() {
				log.Println("read timeout:", rerr)
				// read time out
				conn.Close()
				return
				//return nil,rerr
				//
			} else if (rerr!=io.EOF){
				log.Println("read error:", rerr)
				conn.Close()
				return
				//return nil,rerr
				//
				//
			}
		}
	
		
	}

}

func (sw *Swarm) ListenPeerMessages(peer *Peer) {
	for {
		peer.Connection.SetReadDeadline(time.Now().Add(5 * time.Second))
		rmsg,rerr:=peer.ReadMessage()
		if (rerr==nil) && (rmsg!=nil){
			applog.Trace("NewMessage %x",rmsg.GetContent())
			sw.PeersmsgChan<-rmsg
		}
		if (rerr==io.EOF){

			applog.Trace("peer disconnected - EOF ")
			//peer.Connection.Close()
			sw.RemovePeer(peer)
			break

		}
		
		if netErr, ok := rerr.(net.Error); ok && netErr.Timeout() {
			//log.Println("read timeout:", rerr)
			// time out
			//peer.Connection.Close()
			continue
			//break
		}	
		
	}
}
