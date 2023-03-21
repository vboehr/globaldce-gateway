package rpc

import (
    //"encoding/json"
    //"log"
	"fmt"
	"strconv"
)
func runCacheTorrent(tmpparms []string) Result {
    tmpmagnet:=tmpparms[0]
    tmpdappname:=tmpparms[1]
    tmpdirectory:=tmpparms[2]
	if Mncc==nil{
		//fmt.Println("Uninitiated Client")
		return Result{Type:"CacheTorrentAborted",Data:fmt.Sprintf("Warning: Uninitiated Client")}
	}
    Mncc.AddCacheTorrentRequest(tmpdappname,tmpdirectory,tmpmagnet)
	return Result{Type:"CacheTorrentProcessed"}
}
func runDropTorrent(tmpparms []string) Result {
    tmpmagnet:=tmpparms[0]
	tmpdappname:=tmpparms[1]
    tmpdirectory:=tmpparms[2]
	tmperasefilesflag:=(tmpparms[3]=="true")
	if Mncc==nil{
		//fmt.Println("Uninitiated Client")
		return Result{Type:"DropTorrentAborted",Data:fmt.Sprintf("Warning: Uninitiated Client")}
	}
    Mncc.DropTorrent(tmpmagnet,tmpdappname,tmpdirectory,tmperasefilesflag)
	return Result{Type:"DropTorrentProcessed"}
}

func runProtorizeTorrentPiecesInterval(tmpparms []string) Result {
	//log.Printf("cacheTorrent params %s",request.Params[0])
	tmpmagnet:=tmpparms[0]
	tmpfilepath:=tmpparms[1]
	tmpstartpiece, cerr1 := strconv.Atoi(tmpparms[2])
	if cerr1 != nil {
		return Result{Type:"ProtorizeTorrentPiecesIntervalAborted",Data:fmt.Sprintf("Warning:", cerr1)}
	}
	tmpendpiece, cerr2 := strconv.Atoi(tmpparms[3])
	if cerr2 != nil {
		return Result{Type:"ProtorizeTorrentPiecesIntervalAborted",Data:fmt.Sprintf("Warning:", cerr2)}
	}
	tmpcancelflag:=(tmpparms[3]=="true")
	Mncc.ProtorizeTorrentPiecesInterval(tmpmagnet,tmpfilepath,tmpstartpiece,tmpendpiece,tmpcancelflag)
	//maincontentclient.ProtorizeTorrentPiecesInterval(tmpmagnet,".mp4",0,20)
	return Result{Type:"ProtorizeTorrentPiecesIntervalProcessed"}
}
//
func runProtorizeTorrentDurationPercentageInterval(tmpparms []string) Result {
	//log.Printf("cacheTorrent params %s",request.Params[0])
	tmpmagnet:=tmpparms[0]
	tmpfilepath:=tmpparms[1]
	tmpstartpercentage, cerr1 := strconv.Atoi(tmpparms[2])
	if cerr1 != nil {
		return Result{Type:"ProtorizeTorrentPiecesIntervalAborted",Data:fmt.Sprintf("Warning:", cerr1)}
	}
	tmpendpercentage, cerr2 := strconv.Atoi(tmpparms[3])
	if cerr2 != nil {
		return Result{Type:"ProtorizeTorrentPiecesIntervalAborted",Data:fmt.Sprintf("Warning:", cerr2)}
	}

	tmpcancelflag:=(tmpparms[3]=="true")
	Mncc.ProtorizeTorrentDurationPercentageInterval(tmpmagnet,tmpfilepath,tmpstartpercentage,tmpendpercentage,tmpcancelflag)
	//maincontentclient.ProtorizeTorrentPiecesInterval(tmpmagnet,".mp4",0,20)
	return Result{Type:"ProtorizeTorrentDurationPercentageIntervalProcessed"}
}
  /*      
func runProtorizeTorrentAllPieces(tmpparms []string) Result {
	//log.Printf("cacheTorrent params %s",request.Params[0])
	tmpmagnet:=tmpparms[0]
	tmpfilepath:=tmpparms[1]
	Mncc.ProtorizeTorrentAllPieces(tmpmagnet,tmpfilepath)
	//maincontentclient.ProtorizeTorrentAllPieces(tmpmagnet,".mp4")
	return Result{Type:"ProtorizeTorrentAllPiecesProcessed"}
}*/
//
func runGetTorrentDetails(tmpparms []string) Result {
    tmpmagnet:=tmpparms[0]
    _=tmpmagnet
	tmpTorrentDetailsString:=Mncc.GetTorrentDetails(tmpmagnet)
	return Result{Type:"TorrentDetails",Data:tmpTorrentDetailsString}
}