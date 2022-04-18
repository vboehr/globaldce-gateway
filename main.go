package main
import(
	"github.com/globaldce/go-globaldce/cli"
	//"time"
	"os"
	"github.com/globaldce/go-globaldce/gui"
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
