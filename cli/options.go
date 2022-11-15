package cli
import (
	"github.com/globaldce/globaldce-gateway/applog"
	"github.com/globaldce/globaldce-gateway/daemon"
    "os"
	"strings"
)

func InterpretOptions() {

    for i := 2; i<len(os.Args); i++ {
        tmparg:= os.Args[i]
        if strings.Index(tmparg, "-path=")==0{
            daemon.AppPath=strings.TrimPrefix(tmparg, "-path=")
            applog.Notice("appPath set to: %s",daemon.AppPath)
        }
        if strings.Index(tmparg, "-port=")==0{
            daemon.AppLocalPort=strings.TrimPrefix(tmparg, "-port=")
            applog.Notice("appLocalPort set to: %s",daemon.AppLocalPort)
        }
        if strings.Index(tmparg, "-seed")==0{
            daemon.Seed=true
        }
        if strings.Index(tmparg, "-miningaddressesfile=")==0{
            tmpMiningaddrressesfilepath:=strings.TrimPrefix(tmparg, "-miningaddressesfile=")
            daemon.MAddresses=new(daemon.MiningAddresses)
            daemon.MAddresses.LoadJSONMiningAddressesFile(tmpMiningaddrressesfilepath)
            daemon.Miningaddrressesfileloaded=true
            applog.EnableDisplayTrace()
        }
        if strings.Index(tmparg, "-trace")==0{
            applog.EnableDisplayTrace()
        }

    }


}