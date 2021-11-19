package main
import(
	"github.com/globaldce/globaldce/cli"
	//"time"
	"os"
	"github.com/globaldce/globaldce/gui"
)
func main(){
	if len(os.Args)>=2{
		// cli is needed
		cli.Start("globaldce")
		/*for {
			time.Sleep(time.Minute)
	    	}*/
	} else {
		//fmt.Println("Launching the gui")
		gui.Start("globaldce")

	}

}
