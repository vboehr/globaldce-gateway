package mainchain
import(
	"path/filepath"
	"os"
	"bufio"
	"fmt"
	"github.com/globaldce/globaldce-toolbox/utility"
)


func  (mn *Maincore) CacheExistingFile(path string) (*utility.Extradata,error){
	
	f, err := os.Open(path)
	if err != nil {
		//
		fmt.Println("error:", err)
		return nil,err
	}
	defer f.Close()

	stats,serr:=f.Stat()
	if serr !=nil{
		return nil,serr
	}
	filesize:=stats.Size()
	filebytes := make([]byte,filesize)

	bufreader:=bufio.NewReader(f)
	_,rerr :=bufreader.Read(filebytes)
	if rerr!= nil{
		return nil,rerr
	}
	ed:=utility.NewExtradataFromBytes(filebytes)
	//////////////////////////////////////////

	datafilesdirpath:=filepath.Join(mn.path,"Data","DataFiles")
	
	if _, err := os.Stat(datafilesdirpath); os.IsNotExist(err) {
		os.Mkdir(datafilesdirpath, os.ModePerm)
	}
	/*
	if _, err := os.Stat( filepath.Join(mn.path,"Data","Data000")); os.IsNotExist(err) {
		// path does not exist
		mn.dataf = utility.OpenChunkStorage( filepath.Join(mn.path,"Data","Data"))
		mn.dataf.AddChunk([]byte("emptydata"))
	} else {
		mn.dataf = utility.OpenChunkStorage(filepath.Join(mn.path,"Data","Data"))
	}
	*/
	newdatafilename:=fmt.Sprintf("%x",ed.Hash)
	newdatafilepath:=filepath.Join(datafilesdirpath,newdatafilename)
	fmt.Println("creating file",newdatafilepath)
	cf, err := os.OpenFile(newdatafilepath, os.O_WRONLY|os.O_CREATE, 0755)
	
	if err != nil {
		//
		fmt.Println("error:", err)
	}
	defer cf.Close()
	_, wterr :=  cf.Write(filebytes)
	if wterr != nil {
		//
		fmt.Println("error:", wterr)
	}

	//////////////////////////////////////////
	mn.PutDataFileState(ed.Hash,ed.Size)

	return &ed,nil	
}