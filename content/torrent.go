package content

import (
	"fmt"
	"context"
	"path/filepath"
	"strings"
	//"time"
	"io/fs"
	"os"
	"log"
	"encoding/json"
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
	fmt.Println(cachetorrentrequest)
	contentclient.CacheTorrentRequestChannel <- *cachetorrentrequest

}
type ContentClient struct {
	torrentclient * torrent.Client
	ContentLocation	string
	ctx context.Context
	//ClosingChannel chan bool
	CacheTorrentRequestChannel chan CacheTorrentRequest
	//UncacheTorrentMagnetChannel chan string
}

func Newcontentclient(ctx context.Context,applocation string) *ContentClient{
    contentclient:=new(ContentClient)
	contentclient.ctx=ctx
	contentclient.ContentLocation=filepath.Join(applocation,"Cache","Content")//applocation
	tmpcachelocation:=filepath.Join(applocation,"Cache")
	if _, err := os.Stat(tmpcachelocation); os.IsNotExist(err) {
		os.Mkdir(tmpcachelocation, os.ModePerm)
		fmt.Printf("Creating :%s\n",tmpcachelocation)
	}

	if _, err := os.Stat(contentclient.ContentLocation); os.IsNotExist(err) {
		os.Mkdir(contentclient.ContentLocation, os.ModePerm)
		fmt.Printf("Creating :%s\n",contentclient.ContentLocation)
	}


	//contentclient.AppIsClosing=false
	contentclient.CacheTorrentRequestChannel=make(chan CacheTorrentRequest)
	//contentclient.UncacheTorrentMagnetChannel=make(chan string)
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
		//FilePathMaker: func(opts storage.FilePathMakerOpts) string {
		//	return filepath.Join(opts.File.Path...)
		//},
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
	cfg.DataDir = contentclient.ContentLocation //filepath.Join(contentclient.ContentLocation,"Cache","Content")//
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
		//log.Println("Content client is up...")
		select {
		case tmpcachetorrentrequest:=<-contentclient.CacheTorrentRequestChannel:
			//log.Println("Received a new magnet",nmagnet)
			//=tmpdappname
			//=tmppath
			//cachetorrentrequest.Magnet=tmpmagnet
			go contentclient.CacheTorrent(tmpcachetorrentrequest.DAppName,tmpcachetorrentrequest.Path,tmpcachetorrentrequest.Magnet)
			
			//break
		//case <-contentclient.ClosingChannel:
		//	log.Println("Received signal, breaking...")
		//	break
		//
		case <-contentclient.ctx.Done():
			log.Println("Cancellation signal received. Exiting...")
			break
		//default:
			// Do some work here
		//	time.Sleep(1 * time.Second)
		}
	}
}

////////////////////
////////////////////
func (contentclient *ContentClient)  DeleteFileSystemObject(tmpdappname string, tmppath string) string{
	tmpDir:=filepath.Join(contentclient.ContentLocation,tmpdappname,filepath.FromSlash(tmppath))

	fileInfo, err := os.Stat(tmpDir)
    if err != nil {
        fmt.Println("Error:", err)
        return ""
    }

    // Check if the path refers to a directory
    if fileInfo.Mode().IsDir() {
        fmt.Println("removing directory",tmpDir)
		_=os.RemoveAll(tmpcacheFileDir)
		return ""
    } else {
        fmt.Println("removing file",tmpDir)
		_=os.Remove(tmpcacheFileDir)
		return ""
    }

}
//
func (contentclient *ContentClient)  ScanDirectory(tmpdappname string, tmppath string) string{
	tmpDir:=filepath.Join(contentclient.ContentLocation,tmpdappname,filepath.FromSlash(tmppath))
	
	var err error
	var tmpDirectoryDetails []string
	err = filepath.Walk(tmpDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			if path==tmpDir { // Skipping root path 
				return nil
			}
			tmprelpath,terr:=filepath.Rel(tmpDir, path)
			fmt.Printf("directory : %+v err %v \n",filepath.ToSlash(tmprelpath),terr)//info.Name())
			tmpDirectoryDetails=append(tmpDirectoryDetails,filepath.ToSlash(tmprelpath)+"/")
		} else {
			tmprelpath,terr:=filepath.Rel(tmpDir, path)
			fmt.Printf("file : %+v err %v \n",filepath.ToSlash(tmprelpath),terr)//info.Name())
			tmpDirectoryDetails=append(tmpDirectoryDetails,filepath.ToSlash(tmprelpath))
		}
		//fmt.Printf("visited file or dir: %q\n", path)
		return nil
	})
	if err != nil {
		return fmt.Sprintf("error walking the path %q: %v\n", tmpDir, err)
	}
	tmpDirectoryDetailsBytes, err := json.Marshal(tmpDirectoryDetails)
    if err != nil {
        return fmt.Sprintf("Error:", err)
    }
	return string(tmpDirectoryDetailsBytes)
}
//
func (contentclient *ContentClient)  CacheRawString(tmpdappname string, tmppath string,tmpfilename string, tmprawstring string) {
	tmpcacheDir:=filepath.Join(contentclient.ContentLocation,tmpdappname,filepath.FromSlash(tmppath))
	if _, err := os.Stat(tmpcacheDir); os.IsNotExist(err) {
		os.Mkdir(tmpcacheDir, os.ModePerm)
		//TODO better error handling
	}
	tmpcacheFileDir:=filepath.Join(contentclient.ContentLocation,tmpdappname,filepath.FromSlash(tmppath),tmpfilename)
	_=os.Remove(tmpcacheFileDir)
	f, err := os.OpenFile(tmpcacheFileDir, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		//log.Fatal(err)
		fmt.Println("error:", err)
	}
	defer f.Close()
	rawstringbytes := []byte(tmprawstring)
	_, wserr := f.Write(rawstringbytes)
	if wserr != nil {
		//log.Fatal(err)
		fmt.Println("error:", wserr)
	}
}
////////////////////
func (contentclient *ContentClient)  CacheTorrent(tmpdappname string, tmppath string, tmpmagneturi string) {
	//////////
	tmpdownloadDir:=filepath.Join(contentclient.ContentLocation,tmpdappname,filepath.FromSlash(tmppath))
	t, err := contentclient.AddMagnetWithDownloadDir(tmpmagneturi,tmpdownloadDir)
	//t, err := contentclient.torrentclient.AddMagnet(tmpmagneturi)
	if err != nil {
		log.Print("new torrent error: %w", err)
	}
	//_=tmppath
	
	//_=tmpdappname
	//t.SetDownloadDir(filepath.Join(contentclient.ContentLocation,"Cache",tmpdappname,filepath.FromSlash(tmppath)))// "/path/to/download/directory")
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
//
//
func (contentclient *ContentClient)  DropTorrent(tmpmagneturi string,tmpdappname string,tmpdirectory string,tmperasefilesflag bool) string{
	log.Printf("Droping torrent %s erase files%v",tmpmagneturi,tmperasefilesflag)
	tmpmagnet, perr := metainfo.ParseMagnetUri(tmpmagneturi)
	if perr != nil {
		return ""
	}

	t, ok := contentclient.torrentclient.Torrent(tmpmagnet.InfoHash)
	if !ok {
		return ""
	}
	if t == nil {
		return ""
	}
	if t.Info() == nil {
		return ""
	}
	files := t.Files()
	if files == nil {
		return ""
	}
	t.Drop()
	if tmperasefilesflag {
		for _, filei := range files {
			tmpfpath:=filepath.Join(contentclient.ContentLocation,tmpdappname,filepath.FromSlash(tmpdirectory),filepath.FromSlash(filei.Path()))
			tmpfdir := filepath.Dir(tmpfpath)
			//if _, tmpfdirerr := os.Stat(tmpfdir); !os.IsNotExist(tmpfdirerr) {
			terr:=os.RemoveAll(tmpfdir)//os.Remove(tmpfpath)// 
			log.Printf("Removed folder %s error %v",tmpfdir,terr)
			//}
		}
	}

	
	return "TorrentDropped magnet "+tmpmagneturi
}
//
/*
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
*/
func (contentclient *ContentClient)  ProtorizeTorrentPiecesInterval(tmpmagnet string, tmppath string, beginprioritizedpiece int,endprioritizedpiece int,cancelflag bool) {
	

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
			if cancelflag {
				t.CancelPieces(lastpiece, filei.EndPieceIndex())
			}
		} else {
			filei.SetPriority(torrent.PiecePriorityNone)
		}
	}
}
//
func (contentclient *ContentClient) ProtorizeTorrentDurationPercentageInterval(tmpmagnet string, tmppath string, beginprioritizeddurationpercentage int,endprioritizeddurationpercentage int,cancelflag bool) {
	

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
			d:=int(filei.EndPieceIndex())-int(filei.BeginPieceIndex())
			firstpiece := int(filei.BeginPieceIndex())+(beginprioritizeddurationpercentage*d/100)//int(filei.BeginPieceIndex()+beginprioritizedpiece)
			lastpiece := int(filei.BeginPieceIndex())+(endprioritizeddurationpercentage*d/100)////CustomMin(firstpiece+(endprioritizedpiece-beginprioritizedpiece), int(filei.EndPieceIndex()))
			log.Println("Priority for ",firstpiece,lastpiece)
			t.DownloadPieces(firstpiece, lastpiece)
			if cancelflag {
				t.CancelPieces(lastpiece, filei.EndPieceIndex())
			}
		} else {
			filei.SetPriority(torrent.PiecePriorityNone)
		}
	}
}
//

//
type TorrentFileInfo struct {
	Path string
	Size int
	Progress int
}

type TorrentDetails struct {
	Magnet string  //    `json:"magnet"`
	Name string    //  `json:"name"`
	Nbpeers int    //     `json:"nbpeers"`
	FileInfoArray  []TorrentFileInfo //interface{} `json:"name,omitempty"`
}

func (contentclient *ContentClient)  GetTorrentDetails(tmpmagneturi string) string{
	tmpmagnet, perr := metainfo.ParseMagnetUri(tmpmagneturi)
	if perr != nil {
		return ""
	}

	t, ok := contentclient.torrentclient.Torrent(tmpmagnet.InfoHash)
	if !ok {
		return ""
	}
	if t == nil {
		return ""
	}
	if t.Info() == nil {
		return ""
	}
	files := t.Files()
	if files == nil {
		return ""
	}
	var tmpTorrentDetails TorrentDetails
	tmpTorrentDetails.Magnet=tmpmagneturi
	tmpTorrentDetails.Name=""
	tmpTorrentDetails.Nbpeers=len(t.PeerConns())
	//tmpreturnstring += "*" + tmpmagneturi
	//tmpreturnstring += "*" + "TORRENTNAME"
	//tmpreturnstring += "*" + fmt.Sprintf("%d", len(t.PeerConns())) //"333"//nbpeers

	
	for _, filei := range files {
		//tmpreturnstring += "*" + fmt.Sprintf("%s*%d", filei.Path(), filei.BytesCompleted()*100/filei.Length())
		var tmpTorrentFileInfo TorrentFileInfo
		tmpTorrentFileInfo.Path=filei.Path()
		tmpTorrentFileInfo.Size=int(filei.Length())
		tmpTorrentFileInfo.Progress=int(filei.BytesCompleted()*100/filei.Length())//fmt.Sprintf("%d",filei.BytesCompleted()*100/filei.Length())
		tmpTorrentDetails.FileInfoArray=append(tmpTorrentDetails.FileInfoArray,tmpTorrentFileInfo)
	}
    tmpTorrentDetailsBytes, err := json.Marshal(tmpTorrentDetails)
    if err != nil {
        
        return fmt.Sprintf("Error:", err)
    }
	return string(tmpTorrentDetailsBytes)
}
//
func CustomMin(i int, j int) int {
	if i > j {
		return j
	} else {
		return i
	}
}

