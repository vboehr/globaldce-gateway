package main
import(
	"github.com/globaldce/globaldce-gateway/cli"

	//"time"
	"log"
	"os"
	"path/filepath"
	"strings"
	"github.com/globaldce/globaldce-gateway/gui"
	"github.com/globaldce/globaldce-gateway/daemon"
)
const (
	AppName="globaldce-gateway"
    AppVersion="0.2.6"
	AppID="globaldcegateway.app.testapp"
	AppReleaseType=""//"mobile"
)
func main(){
	switch AppReleaseType {
	case "":
		daemon.AppPath="."
	case "mobile":
		daemon.AppPath=os.TempDir()
	default:
			tmpapppath, err := filepath.Abs(filepath.Dir(os.Args[0]))
				if err != nil {
					log.Printf("%v",err)
					os.Exit(0)
				}
				daemon.AppPath=tmpapppath	
	}
	log.Printf("AppPath: %s \n",daemon.AppPath)
	if (len(os.Args)==1)||(strings.ToLower(os.Args[1])=="gui"){
		//fmt.Println("Launching the cli")er
		gui.Start(AppName,AppVersion,AppID)
	} else {
		//fmt.Println("Launching the gui")
		cli.Start(AppName,AppVersion)

	}

}
