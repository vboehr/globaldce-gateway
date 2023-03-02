package content

import (
	//"fmt"
	"context"
	"path/filepath"
	"strings"
	"time"
	"log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
	"github.com/anacrolix/torrent/metainfo"
)

//var mainclient *torrent.Client
type CacheTorrentRequest struct {
	DAppName string
	Path string
	Magnet string
}
func (contentclient *ContentClient) AddCacheTorrentRequest(tmpdappname string,tmppath string,tmpmagnet string) {
	cachetorrentrequest:=new(CacheTorrentRequest)
	cachetorrentrequest.DAppName=tmpdappname
	cachetorrentrequest.Path=tmppath
	cachetorrentrequest.Magnet=tmpmagnet
	contentclient.CacheTorrentRequestChannel<-*cachetorrentrequest

}
type ContentClient struct {
	torrentclient * torrent.Client
	AppLocation	string
	ctx context.Context
	//ClosingChannel chan bool
	CacheTorrentRequestChannel chan CacheTorrentRequest
	UncacheTorrentMagnetChannel chan string
}

func Newcontentclient(ctx context.Context,applocation string) *ContentClient{
    contentclient:=new(ContentClient)
	contentclient.ctx=ctx
	contentclient.AppLocation=applocation
	//contentclient.AppIsClosing=false
	contentclient.CacheTorrentRequestChannel=make(chan CacheTorrentRequest)
	contentclient.UncacheTorrentMagnetChannel=make(chan string)
	//contentclient.ClosingChannel=make(chan bool)
	//wlt.Walletloaded=false     
	//go Gensequentialwallet(wlt,seedString)
	return contentclient
}

func (contentclient *ContentClient) AddMagnetWithDownloadDir(uri string,downloadDir string) (*torrent.Torrent, error) {

	spec, err := torrent.TorrentSpecFromMagnetUri(uri)
	if err != nil {
		return nil,err
	}
	

	pieceCompletion, err := storage.NewDefaultPieceCompletionForDir(downloadDir)
	if err != nil {
		log.Println("error creating piece completion:", err)
		return nil,err
	}
	spec.Storage=storage.NewFileOpts(storage.NewFileClientOpts{
		ClientBaseDir: downloadDir,
		FilePathMaker: func(opts storage.FilePathMakerOpts) string {
			return filepath.Join(opts.File.Path...)
		},
		TorrentDirMaker: nil,
		PieceCompletion: pieceCompletion,
	})

	//////////
	t, _, aerr := contentclient.torrentclient.AddTorrentSpec(spec)
	_=t

	 
	log.Println("error",aerr)
	//<-t.GotInfo()
	//t.DownloadAll()
	//contentclient.torrentclient.WaitAll()
	//(cl *torrent.Client) AddMagnetWithDownloadDir(uri string,downloadDir) (T *Torrent, err error) 
	return t,aerr
}
func (contentclient *ContentClient) Initcontentclient() {
	cfg := torrent.NewDefaultClientConfig()
	// cfg.Seed = true
	cfg.DataDir = contentclient.AppLocation //
	// cfg.NoDHT = true
	// cfg.DisableTCP = true
	 cfg.DisableUTP = true
	// cfg.DisableAggressiveUpload = false
	// cfg.DisableWebtorrent = false
	// cfg.DisableWebseeds = false
	var err error
	contentclient.torrentclient, err = torrent.NewClient(cfg)
	if err != nil {
		log.Print("new torrent client: %w", err)
		return //fmt.Errorf("new torrent client: %w", err)
	}
	log.Print("new torrent client INITIATED")
	defer contentclient.torrentclient.Close()
	//go func() {
		/*for {
			if contentclient.AppIsClosing {
				log.Print("closing contentclient")
				break
			}
			log.Print("torrentclient running ...")
			time.Sleep(1 * time.Second)
		}*/
	//}()
	for {
		log.Println("Content client is up...")
		select {
		case tmpcachetorrentrequest:=<-contentclient.CacheTorrentRequestChannel:
			//log.Println("Received a new magnet",nmagnet)
			//=tmpdappname
			//=tmppath
			//cachetorrentrequest.Magnet=tmpmagnet
			contentclient.CacheTorrent(tmpcachetorrentrequest.DAppName,tmpcachetorrentrequest.Path,tmpcachetorrentrequest.Magnet)
			
			//break
		//case <-contentclient.ClosingChannel:
		//	log.Println("Received signal, breaking...")
		//	break
		//
		case <-contentclient.ctx.Done():
			log.Println("Cancellation signal received. Exiting...")
			break
		default:
			// Do some work here
			time.Sleep(1 * time.Second)
		}
	}
}

////////////////////
func (contentclient *ContentClient)  CacheTorrent(tmpdappname string, tmppath string, tmpmagneturi string) {



	//////////
	tmpdownloadDir:=filepath.Join(contentclient.AppLocation,tmpdappname,filepath.FromSlash(tmppath))
	t, err := contentclient.AddMagnetWithDownloadDir(tmpmagneturi,tmpdownloadDir)
	//t, err := contentclient.torrentclient.AddMagnet(tmpmagneturi)
	if err != nil {
		log.Print("new torrent error: %w", err)
	}
	//_=tmppath
	
	//_=tmpdappname
	//t.SetDownloadDir(filepath.Join(contentclient.AppLocation,"Cache",tmpdappname,filepath.FromSlash(tmppath)))// "/path/to/download/directory")
	/////////



	///////
	<-t.GotInfo()

	log.Printf("added magnet %s\n", tmpmagneturi)
	files := t.Files()
	for _, filei := range files {
			filei.SetPriority(torrent.PiecePriorityNone)
	}
	/*
	totalsize := int64(0)
	tmppreviewfile := ""
	tmppreviewfilesize := int64(0)


	
	for _, filei := range files {
		if (filei.Length() > tmppreviewfilesize) && (strings.Contains(filei.Path(), ".mp4")) {
			tmppreviewfile = filei.Path()
			totalsize += filei.Length()
		}
	}

	for _, filei := range files {
		if tmppreviewfile == filei.Path() {
			firstprioritizedpiece := int(filei.BeginPieceIndex())
			lastprioritizedpiece := CustomMin(firstprioritizedpiece+20, int(filei.EndPieceIndex()))

			t.DownloadPieces(firstprioritizedpiece, lastprioritizedpiece)
			t.CancelPieces(lastprioritizedpiece, filei.EndPieceIndex())
		} else {
			filei.SetPriority(torrent.PiecePriorityNone)
		}
	}
	*/


	//for {
		/*
		if (!IsSavedItemWithMagnet(tmpmagneturi)) && (!s.IsMainTorrent(tmpmagneturi)) && (!IsPreviewingTorrent(tmpmagneturi)) {
			log.Println("Torrent removed", tmpmagneturi)
			t.Drop()
			return
		}
		*/
	//	time.Sleep(8 * time.Second)
	//}
}
func (contentclient *ContentClient)  ProtorizeTorrentAllPieces(tmpmagnet string, tmppath string) {
	tmpmagnetobj, perr := metainfo.ParseMagnetUri(tmpmagnet)
	if perr != nil {
		log.Println("Error ",perr)
		return 
	}
	t, ok := contentclient.torrentclient.Torrent(tmpmagnetobj.InfoHash)
	if !ok {
		log.Println("Torrent not found ")
		return 
	}
	files := t.Files()
	for _, filei := range files {
		if strings.Contains(filei.Path(), tmppath) {//tmppreviewfile == filei.Path() {
			firstpiece := int(filei.BeginPieceIndex())
			lastpiece := int(filei.EndPieceIndex())
			log.Println("Priority for ",firstpiece,lastpiece)
			t.DownloadPieces(firstpiece, lastpiece)
			t.CancelPieces(lastpiece, filei.EndPieceIndex())
		} else {
			filei.SetPriority(torrent.PiecePriorityNone)
		}
	}
}
func (contentclient *ContentClient)  ProtorizeTorrentPiecesInterval(tmpmagnet string, tmppath string, beginprioritizedpiece int,endprioritizedpiece int) {
	

	tmpmagnetobj, perr := metainfo.ParseMagnetUri(tmpmagnet)

	if perr != nil {
		log.Println("Error ",perr)
		return 
	}

	t, ok := contentclient.torrentclient.Torrent(tmpmagnetobj.InfoHash)
	if !ok {
		log.Println("Torrent not found ")
		return 
	}
	files := t.Files()
	for _, filei := range files {
		if strings.Contains(filei.Path(), tmppath) {//tmppreviewfile == filei.Path() {
			firstpiece := int(filei.BeginPieceIndex()+beginprioritizedpiece)
			lastpiece := CustomMin(firstpiece+(endprioritizedpiece-beginprioritizedpiece), int(filei.EndPieceIndex()))
			log.Println("Priority for ",firstpiece,lastpiece)
			t.DownloadPieces(firstpiece, lastpiece)
			t.CancelPieces(lastpiece, filei.EndPieceIndex())
		} else {
			filei.SetPriority(torrent.PiecePriorityNone)
		}
	}


}
func CustomMin(i int, j int) int {
	if i > j {
		return j
	} else {
		return i
	}
}

