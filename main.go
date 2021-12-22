package main
import(
	"github.com/globaldce/globaldce-toolbox/cli"
	//"time"
	"os"
	"github.com/globaldce/globaldce-toolbox/gui"
)
func main(){
	if len(os.Args)>=2{
		//fmt.Println("Launching the cli")
		cli.Start("globaldce-toolbox")
	} else {
		//fmt.Println("Launching the gui")
		gui.Start("globaldce-toolbox")

	}

}
