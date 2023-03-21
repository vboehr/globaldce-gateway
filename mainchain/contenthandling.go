package mainchain
import(
	"path/filepath"
	//"os"
	//"bufio"
	"fmt"
	//"github.com/globaldce/globaldce-gateway/utility"
	"net/http"
    "html/template"
)

func  (mn *Maincore) ServeContent(name string)  error{
	//p:="/"+name+"/sytle.css"
	//http.HandleFunc(p, HandleContent)
	//http.Handle("/", http.FileServer("./style.css"))
	//http.Handle("/"+name+"/", fileServer("./style.css"))
	//http.Handle("/"+name+"/",http.FileServer(http.Dir("./Cache/Content/"+name+"/")))

	_=name
	return nil
}

func HandleContent(w http.ResponseWriter, r *http.Request) {
    contentfilepath:= filepath.Join("Cache","Content",filepath.FromSlash(r.URL.Path[1:]))////filepath.FromSlash() for win compatibility
    fmt.Printf("Content file path %s loading ... url %s \n",contentfilepath,r.URL.Path)
	//os.Exit(0)
	doc, err := template.ParseFiles(contentfilepath)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	doc.Execute(w, nil)
}
/*
func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/index.html", HandleContent)
	//r.HandleFunc("/requestfreecoins", createRequestFreeCoinsHandler).Methods("POST")
	//r.HandleFunc("/requestfreecoinscaptcha", createRequestFreeCoinsCaptchaHandler).Methods("POST")
	return r
}
*/
func InitLocalhost(){
	//http.HandleFunc("/", HandleContent)/////////////////Was removed
	//http.Handle("/static/", http.FileServer(http.Dir("/")))
	
	fs := http.FileServer(http.Dir("./Cache/Content"))
	http.Handle("/", http.StripPrefix("/", fs))

	fmt.Println(http.ListenAndServe(":8088", nil))

}
/*
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
	//
	//if _, err := os.Stat( filepath.Join(mn.path,"Data","Data000")); os.IsNotExist(err) {
		// path does not exist
	//	mn.dataf = utility.OpenChunkStorage( filepath.Join(mn.path,"Data","Data"))
	//	mn.dataf.AddChunk([]byte("emptydata"))
	//} else {
	//	mn.dataf = utility.OpenChunkStorage(filepath.Join(mn.path,"Data","Data"))
	//}
	//
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
*/

func  (mn *Maincore) PushRegistredNameCommit(name []byte,contentid []byte)  {

	mn.PutRegistredNameCommitState(name,contentid) 

}
/*
func  (mn *Maincore) GetDataFile(hash utility.Hash) ([]byte,error) {
	//TODO check if file has been add to mainstate
	datafilesdirpath:=filepath.Join(mn.path,"Data","DataFiles")
	datafilename:=fmt.Sprintf("%x",hash)
	datafilepath:=filepath.Join(datafilesdirpath,datafilename)
	f, err := os.Open(datafilepath)
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
	//ed:=utility.NewExtradataFromBytes(filebytes)
	//////////////////////////////////////////
	return filebytes,nil

}
*/