package daemon
import (
	"github.com/globaldce/go-globaldce/mainchain"
	"github.com/globaldce/go-globaldce/wire"
	"github.com/globaldce/go-globaldce/wallet"
	//"path/filepath"
	//"fmt"
	//"os"
	//"github.com/globaldce/go-globaldce/applog"
)
	var (

		MainwalletFilePath=""
		MainwalletFileKey []byte
		AppPath=""
		AppLocalPort=""
		Walletloaded=false
		Miningrequested=false
		Miningrunning=false
		HotMining=false
		Managingwalletrequested=false
		Seed=false
		SyncingMinNbPeers=0
		AppIsClosing=false
		Mn *mainchain.Maincore
		Wireswarm *wire.Swarm
		Wlt *wallet.Wallet
	)

