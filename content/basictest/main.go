package main

import (
	"fmt"
	"github.com/globaldce/globaldce-gateway/content"
	"context"
	"time"
)

func main() {
	fmt.Println("Hello")

	// Create a context with a cancellation function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	maincontentclient:=content.Newcontentclient(ctx,"./")
	go maincontentclient.Initcontentclient()
	tmpmagnet:="magnet:?xt=urn:btih:08ada5a7a6183aae1e09d831df6748d566095a10&dn=Sintel&tr=udp%3A%2F%2Fexplodie.org%3A6969&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Ftracker.empire-js.us%3A1337&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=wss%3A%2F%2Ftracker.btorrent.xyz&tr=wss%3A%2F%2Ftracker.fastcast.nz&tr=wss%3A%2F%2Ftracker.openwebtorrent.com&ws=https%3A%2F%2Fwebtorrent.io%2Ftorrents%2F&xs=https%3A%2F%2Fwebtorrent.io%2Ftorrents%2Fsintel.torrent"
	//maincontentclient.AddMagnetWithDownloadDir(tmpmagnet,"./Cache/cooldapp")
	time.Sleep(5 * time.Second)
	maincontentclient.AddCacheTorrentRequest("cooldapp","",tmpmagnet)
	time.Sleep(5 * time.Second)
	//maincontentclient.ProtorizeTorrentPiecesInterval(tmpmagnet,".mp4",0,20)
	maincontentclient.ProtorizeTorrentDurationPercentageInterval(tmpmagnet,".mp4","0","100")
	time.Sleep(200 * time.Second)
	fmt.Println("closing")
	cancel()

}