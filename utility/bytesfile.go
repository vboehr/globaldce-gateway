package utility 

import
(
	//"encoding/binary"
	//"bytes"
	//"fmt"
	//"os"
	"io/ioutil"
)
func SaveBytesFile(bytesfilebytes []byte,bytesfilepath string) error{
	err :=ioutil.WriteFile(bytesfilepath,bytesfilebytes,0644)
	return err
}
func LoadBytesFile(path string) (*[]byte,error){
	b,err :=ioutil.ReadFile(path)
	return &b,err
}
/*
func SaveBytesFile(bytesfilebytes []byte,bytesfilepath string) error{

	bufferBytesfiletype := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufferBytesfiletype, uint32(1))// type 1

	bufferBytesfileseize := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufferBytesfileseize, uint32(len(bytesfilebytes)))
	
	f, err := os.OpenFile(bytesfilepath, os.O_WRONLY|os.O_CREATE, 0755)
	
	if err != nil {
		//
		fmt.Println("error:", err)
	}
	//defer f.Close()
	_, wterr := f.Write(bufferBytesfiletype)
	if wterr != nil {
		//
		fmt.Println("error:", wterr)
	}
	_, werr := f.Write(bufferBytesfileseize)
	if werr != nil {
		//
		fmt.Println("error:", werr)
	}
	_, wserr := f.Write(bytesfilebytes)
	if wserr != nil {
		//
		fmt.Println("error:", wserr)
	}
	//fmt.Println("Bytes saved.")
	return nil
}
func LoadBytesFile(path string) (*[]byte,error){
	
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	
	if err != nil {
		//
		fmt.Println("error:", err)
	}
	//defer f.Close()
	bufferBytesfiletype := make([]byte, 4)
	_, rterr := f.Read(bufferBytesfiletype)
	if rterr != nil {
		//
		fmt.Println("error:", rterr)
	}
	var bytesfiletype uint32
	readerBytesfiletype := bytes.NewReader(bufferBytesfiletype)

	binary.Read(readerBytesfiletype, binary.LittleEndian, &bytesfiletype)
	fmt.Println("type:", bytesfiletype)

	bufferBytesfileseize := make([]byte, 4)
	_, rserr := f.Read(bufferBytesfileseize)
	if rserr != nil {
		//
		fmt.Println("error:", rserr)
	}
	var bytesfileseize uint32
	readerBytesfileseize := bytes.NewReader(bufferBytesfileseize)

	binary.Read(readerBytesfileseize, binary.LittleEndian, &bytesfileseize)
	fmt.Println("seize:", bytesfileseize)
	bytesfilerawbytes := make([]byte, bytesfileseize)
	_, rerr := f.Read(bytesfilerawbytes)
	if rerr != nil {
		//
		fmt.Println("error:", rerr)
	}

	return &bytesfilerawbytes,nil
}
*/