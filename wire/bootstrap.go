package wire
import (
	"github.com/globaldce/globaldce-toolbox/applog"
	//"fmt"
	//"net"
	//"io"
	//"time"
	//"log"
)

func (sw *Swarm)  Bootstrap(){
	bootstrapaddrs :=[]string{"***","***"}
	applog.Notice("Bootstraping ... ")
	for _,addr := range bootstrapaddrs {
		sw.HintNewPeer(addr)
		/*
		conn,err:=net.Dial("tcp","****.ddns.net:15555")
		if err!=nil {
			fmt.Printf("err %v\n",err)
		} else {
			fmt.Printf("ip %s\n",conn.RemoteAddr())
		}
		*/
	} 

}