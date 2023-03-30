package main

import(
	//"path/filepath"
	//"os"
	//"bufio"
	"fmt"
	"github.com/globaldce/globaldce-gateway/mainchain"
	//"path/filepath"
)


func main(){
	
	//contentid,err:=mainchain.GetContentIdWithUniformPieceSize("./Tmp/dapptest/main.go",32)
	//_=mainchain.CacheExistingFile("./Tmp/dapptest/main.go","Cache/Content/dapptest/main.go")
	contentid,err:=mainchain.CacheExistingDirectoryWithUniformPieceSize("./Tmp/cooldapp/",1024,[]byte("cooldapp"))
	fmt.Printf("contentid  %v err %v\n",contentid,err)

}


