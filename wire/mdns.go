package wire
import (
	"github.com/globaldce/globaldce-toolbox/applog"
	"os"
	"time"
	"github.com/globaldce/mdns-1"
)
func (sw *Swarm) StartMDNSServer(){

		// Setup our service export
		host, _ := os.Hostname()
		info := []string{sw.GetLocalIP()}//"My awesome service569"}
		service, _ := mdns.NewMDNSService(host, "_globaldce_MDNS._udp", "", "", 8008, nil, info)
	
		// Create the mDNS server, defer shutdown
		server, _ := mdns.NewServer(&mdns.Config{Zone: service})
		defer server.Shutdown()
		for {
			
			time.Sleep(time.Minute)

		}
		
}
func (sw *Swarm) StartMDNSClient(){


		// Make a channel for results and start listening
		entriesCh := make(chan *mdns.ServiceEntry, 4)
		defer close(entriesCh)
		//for {
			go func() {
				for entry := range entriesCh {
				applog.Trace("Got new entry: %v", entry.InfoFields[0])
				//for _,addr := range entry.InfoFields {
				sw.IpaddrChan<-entry.InfoFields[0]
				//}
				}
			}()
		
			// Start the lookup
			mdns.Lookup("_globaldce_MDNS._udp", entriesCh)
			time.Sleep(time.Minute)
			//
			//applog.Notice("\n**********%s",entriesCh)
		//}
		
		
	
		
}