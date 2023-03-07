package rpc

import (
    //"encoding/json"
    //"log"
	"fmt"
	"strconv"
)
func runCmdCacheTorrent(tmpparms []string) string {
    tmpmagnet:=tmpparms[0]
    tmpdappname:=tmpparms[1]
    tmpdirectory:=tmpparms[2]
	if Mncc==nil{
		fmt.Println("Uninitiated Client")
	}
    Mncc.AddCacheTorrentRequest(tmpdappname,tmpdirectory,tmpmagnet)
	return "CacheTorrentProcessed"
}

func runCmdProtorizeTorrentPiecesInterval(tmpparms []string) string {
	//log.Printf("cacheTorrent params %s",request.Params[0])
	tmpmagnet:=tmpparms[0]
	tmpfilepath:=tmpparms[1]
	tmpstartpiece, cerr1 := strconv.Atoi(tmpparms[2])
	if cerr1 != nil {
		return fmt.Sprintf("Error:", cerr1)
	}
	tmpendpiece, cerr2 := strconv.Atoi(tmpparms[3])
	if cerr2 != nil {
		return fmt.Sprintf("Error:", cerr2)
	}
	Mncc.ProtorizeTorrentPiecesInterval(tmpmagnet,tmpfilepath,tmpstartpiece,tmpendpiece)
	//maincontentclient.ProtorizeTorrentPiecesInterval(tmpmagnet,".mp4",0,20)
	return "ProtorizeTorrentPiecesIntervalProcessed"
}
        
func runCmdProtorizeTorrentAllPieces(tmpparms []string) string {
	//log.Printf("cacheTorrent params %s",request.Params[0])
	tmpmagnet:=tmpparms[0]
	tmpfilepath:=tmpparms[1]
	Mncc.ProtorizeTorrentAllPieces(tmpmagnet,tmpfilepath)
	//maincontentclient.ProtorizeTorrentAllPieces(tmpmagnet,".mp4")
	return "ProtorizeTorrentAllPiecesProcessed"
}	
//
func runCmdGetTorrentDetails(tmpparms []string) string {
    tmpmagnet:=tmpparms[0]
    _=tmpmagnet
	return "TorrentDetails"
}