package main
import(
	"github.com/globaldce/globaldce-gateway/cli"

	//"time"
	"os"
	"strings"
	"github.com/globaldce/globaldce-gateway/gui"
)
const (
	AppName="globaldce-gateway"
    AppVersion="0.2.5"
	AppPath="."
)
func main(){
	if (len(os.Args)==1)||(strings.ToLower(os.Args[1])=="gui"){
		//fmt.Println("Launching the cli")er
		gui.Start(AppName,AppVersion,AppPath)
	} else {
		//fmt.Println("Launching the gui")
		cli.Start(AppName,AppVersion,AppPath)

	}

}
