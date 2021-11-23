package main
import(
	"github.com/globaldce/globaldce-toolbox/cli"
	//"time"
	"os"
	"github.com/globaldce/globaldce-toolbox/gui"
)
func main(){
	if len(os.Args)>=2{
		// cli is needed
		cli.Start("globaldce-toolbox")
		/*for {
			time.Sleep(time.Minute)
	    	}*/
	} else {
		//fmt.Println("Launching the gui")
		gui.Start("globaldce-toolbox")

	}

}
