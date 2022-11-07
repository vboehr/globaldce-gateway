package main
import(
	"github.com/globaldce/globaldce-gateway/cli"
	//"time"
	"os"
	"github.com/globaldce/globaldce-gateway/gui"
)
func main(){
	if len(os.Args)>=2{
		//fmt.Println("Launching the cli")
		cli.Start("globaldce-gateway"+cli.AppVersion)
	} else {
		//fmt.Println("Launching the gui")
		gui.Start("globaldce-gateway"+cli.AppVersion)

	}

}
