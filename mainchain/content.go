
package mainchain

import(
	"encoding/json"
	//"encoding/hex"
	"path/filepath"
	"os"
	"bufio"
	"io"
	"fmt"
	"github.com/globaldce/globaldce-gateway/utility"
)




const (
	ContentTypeUniformPieceSize=1
	ContentDefaultUniformPieceSize=1024
	ContentCurrentVersion=1
)
/*
type Content struct {
	RegistredName []byte
	ContentId []byte
	Priority uint32
}*/
/*
func NewContent(name []byte, contentid []byte) *Content {
	tmpcontent:=new(Content)
	tmpcontent.RegistredName=name
	tmpcontent.ContentId=contentid

	return tmpcontent
}
*/

func GetContentIdWithUniformPieceSize(contentfilepath string,piecesize uint64) ([]byte,error){



	f, err := os.Open(filepath.FromSlash(contentfilepath))
	if err != nil {
		//
		fmt.Println("file open error:", err)
		return nil,err
	}
	defer f.Close()

	stats,serr:=f.Stat()
	if serr !=nil{
		return nil,serr
	}

	filesize:=stats.Size()

	filepiecebytes := make([]byte,piecesize)

	bufreader:=bufio.NewReader(f)


	//func ComputeRoot(hashes *[]Hash) Hash{
	var Hashes []utility.Hash
	//fmt.Printf(" %v len %d\n",contentbytes,len(contentbytes))
	var count int=0
	var rerr error
	for i:=0;i<int(filesize);i+=count{
		//var piece []byte
		//if i+int(piecesize)<filesize{
			count,rerr =bufreader.Read(filepiecebytes)
			if rerr!= nil{
				return nil,rerr
			}
			//piece=contentbytes[i:i+int(piecesize)]
			Hashes=append(Hashes,utility.ComputeHash(filepiecebytes[:count]))
		//} else {
		//	piece=contentbytes[i:]
		//}
		//fmt.Printf(" %v Hash %v\n",piece,utility.ComputeHash(piece))
		
	}
	tmproot:=utility.ComputeRoot(&Hashes)
	//fmt.Printf(" Root Hash %v\n",tmproot)
	//
	
	tmpbw:=utility.NewBufferWriter()
	tmpbw.PutUint32(ContentTypeUniformPieceSize)// 
	tmpbw.PutVarUint(uint64(piecesize))
	tmpbw.PutHash(tmproot)
	contentid:=tmpbw.GetContent()
	return contentid,nil
}


type ContentFileInfo struct {
	FilePath string
	FileName string
	ContentId []byte
	Priority uint32
}

type ContentDirInfo struct {
	Version uint32
	ContentSubDirPathArray []string
	ContentFileInfoArray []ContentFileInfo
}




func CacheExistingDirectoryWithUniformPieceSize(dirpath string,piecesize uint64,registredname []byte) ([]byte,error){

	//var files []string
	var oldcontentdirinfo ContentDirInfo
	oldcontentdirinfo.Version=ContentCurrentVersion
	werr := filepath.Walk(dirpath, func(path string, fileinfo os.FileInfo, err error) error {
		if err!=nil{
			return err
		}
		if !(fileinfo.IsDir()) {
			if fileinfo.Name()=="contentroot.json"{
				return nil
			}
			fmt.Printf("path %s file %s\n",path,fileinfo.Name())
			finalfilepath:=filepath.ToSlash(path)//for accessing the file use: filepath.FromSlash(path string)
			tmpcontentfileinfo:=new(ContentFileInfo)
			tmpcontentfileinfo.FilePath=finalfilepath
			tmpcontentfileinfo.FileName=fileinfo.Name()
			
			contentid,cerr:=GetContentIdWithUniformPieceSize(finalfilepath,piecesize)
			if cerr != nil {
				//
				fmt.Println("GetContentIdWithUniformPieceSize error:", cerr)
				return cerr
			}
			
			tmpcontentfileinfo.ContentId=contentid
			
			oldcontentdirinfo.ContentFileInfoArray = append(oldcontentdirinfo.ContentFileInfoArray, *tmpcontentfileinfo)

		} else if fileinfo.Name()!=dirpath {
			fmt.Printf("Root %s\n",dirpath)

			oldcontentdirinfo.ContentSubDirPathArray=append(oldcontentdirinfo.ContentSubDirPathArray,path)
		}
		return nil
	})
	if werr != nil {
		//
		fmt.Println("Walking error:", werr)
		return nil,werr
	}
	
	//
	var newcontentdirinfo ContentDirInfo
	newcontentdirinfo.Version=ContentCurrentVersion
	newcontentdirinfo.ContentSubDirPathArray=append(newcontentdirinfo.ContentSubDirPathArray,"")
	basepath:=oldcontentdirinfo.ContentSubDirPathArray[0]
	for i:=1;i<len(oldcontentdirinfo.ContentSubDirPathArray);i++{
		
		//relpath, rperr := filepath.Rel(basepath, oldcontentdirinfo.ContentSubDirPathArray[i])
		relpath, rperr := FindRelativePath(basepath, oldcontentdirinfo.ContentSubDirPathArray[i])
		_=rperr
		//newcontentdirinfo.ContentSubDirPathArray[i]=relpath
		newcontentdirinfo.ContentSubDirPathArray=append(newcontentdirinfo.ContentSubDirPathArray,relpath)
		//newpath:=filepath.Join("./","Cache","Content",relpath)
		//if _, err := os.Stat(newpath); os.IsNotExist(err) {
		//	os.Mkdir(newpath, os.ModePerm)
		//	fmt.Printf("Path created %s\n",newpath)
		//}
	}
	//
	
	//basepath:=oldcontentdirinfo.ContentSubDirPathArray[0]


	for i:=0;i<len(oldcontentdirinfo.ContentFileInfoArray);i++{
		//relpath, rperr := filepath.Rel(basepath, oldcontentdirinfo.ContentFileInfoArray[i].FilePath)
		relpath, rperr := FindRelativePath(basepath, oldcontentdirinfo.ContentFileInfoArray[i].FilePath)
		_=rperr
		tmpcontentfileinfo:=new(ContentFileInfo)
		tmpcontentfileinfo.FilePath=relpath
		tmpcontentfileinfo.FileName=oldcontentdirinfo.ContentFileInfoArray[i].FileName
		tmpcontentfileinfo.ContentId=oldcontentdirinfo.ContentFileInfoArray[i].ContentId
		tmpcontentfileinfo.Priority=uint32(len(tmpcontentfileinfo.FilePath)-len(tmpcontentfileinfo.FileName)+1) //
		newcontentdirinfo.ContentFileInfoArray=append(newcontentdirinfo.ContentFileInfoArray,*tmpcontentfileinfo)
		//newfilepath:=filepath.Join("./","Cache","Content",string(registredname),/***/newcontentdirinfo.ContentFileInfoArray[i].FileName)

	}
	//
	newcontentdirinfobytes, merr := json.Marshal(newcontentdirinfo)
	if merr != nil {
		fmt.Println("newcontentdirinfo serialize error:", merr)
		return nil,merr
	}
	contentjsonpath:=basepath+"/"+"contentroot.json"
	fmt.Println("Saving ",contentjsonpath)
	serr:=utility.SaveBytesFile(newcontentdirinfobytes,filepath.FromSlash(contentjsonpath))
	if serr != nil {
		fmt.Println("contentroot.json error:", serr)
		return nil,serr
	}
	//
	dircontentid,cerr:=GetContentIdWithUniformPieceSize(contentjsonpath,piecesize)
	if cerr != nil {
		//
		fmt.Println("GetContentIdWithUniformPieceSize error:", cerr)
		return nil,cerr
	}

	//
	for i:=0;i<len(newcontentdirinfo.ContentSubDirPathArray);i++{
		
		newpath:="./Cache/Content/"+string(registredname)+"/"+newcontentdirinfo.ContentSubDirPathArray[i]
		fmt.Printf("New Path ?? %s\n",newpath)
		if _, err := os.Stat(newpath); os.IsNotExist(err) {
			os.Mkdir(newpath, os.ModePerm)
			fmt.Printf("Path created %s\n",newpath)
		}
	}
	//
	CacheExistingFile(basepath+"/contentroot.json","./Cache/Content/"+string(registredname)+"/contentroot.json")
	for i:=0;i<len(oldcontentdirinfo.ContentFileInfoArray);i++{
		cachingpath:="./Cache/Content/"+string(registredname)+"/"+newcontentdirinfo.ContentFileInfoArray[i].FilePath
		CacheExistingFile(oldcontentdirinfo.ContentFileInfoArray[i].FilePath,cachingpath)
	}
	//CacheExistingFile("./Tmp/dapptest","main.go","Cache/Content/dapptest/123")
	fmt.Printf("\n Finally ... new %s\n",newcontentdirinfobytes)
	return dircontentid,nil

}



func CacheExistingFile(initialpath string,cachepath string) (error){
		fmt.Println("CacheExistingFile",initialpath,cachepath)
	    // open input file 
		fi, err := os.Open( filepath.FromSlash(initialpath))
		if err != nil {
			return err
		}
		// close fi on exit and check for its returned error
		defer func() {
			if err := fi.Close(); err != nil {
				return 
			}
		}()
	
		// open output file
		fo, err := os.Create( filepath.FromSlash(cachepath))
		if err != nil {
			return err
		}
		// close fo on exit and check for its returned error
		defer func() {
			if err := fo.Close(); err != nil {
				return 
			}
		}()
	
		// make a buffer to keep chunks that are read
		buf := make([]byte, 1024)
		for {
			// read a chunk
			n, err := fi.Read(buf)
			if err != nil && err != io.EOF {
				return err
			}
			if n == 0 {
				break
			}
	
			// write a chunk
			if _, err := fo.Write(buf[:n]); err != nil {
				return err
			}
		}

	/*
	f, err := os.Open(path)
	if err != nil {
		//
		fmt.Println("file open error:", err)
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
	*/

	//if _, err := os.Stat( filepath.Join(mn.path,"Data","Data000")); os.IsNotExist(err) {
		// path does not exist
	//	mn.dataf = utility.OpenChunkStorage( filepath.Join(mn.path,"Data","Data"))
	//	mn.dataf.AddChunk([]byte("emptydata"))
	//} else {
	//	mn.dataf = utility.OpenChunkStorage(filepath.Join(mn.path,"Data","Data"))
	//}
	/*
	newdatafilename:=fmt.Sprintf("%x",contentid)
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
	*/
	//////////////////////////////////////////


	return nil	
}
func FindRelativePath(firstpath string, secondpath string) (string,error){
	relpath, rperr := filepath.Rel(filepath.FromSlash(firstpath),filepath.FromSlash(secondpath))
	return filepath.ToSlash(relpath), rperr
}