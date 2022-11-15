package main
import(
	"github.com/globaldce/globaldce-gateway/cli"

	//"time"
	"os"
	"strings"
	"github.com/globaldce/globaldce-gateway/gui"
)
func main(){
	if (len(os.Args)==1)||(strings.ToLower(os.Args[1])=="gui"){
		//fmt.Println("Launching the cli")
		gui.Start("globaldce-gateway"+cli.AppVersion)
	} else {
		//fmt.Println("Launching the gui")
		cli.Start("globaldce-gateway"+cli.AppVersion)

	}

}
