package daemon

import (
	"github.com/globaldce/globaldce-gateway/content"
	"github.com/globaldce/globaldce-gateway/mainchain"
	"github.com/globaldce/globaldce-gateway/wallet"
	"github.com/globaldce/globaldce-gateway/wire"
	//"path/filepath"
	//"fmt"
	//"os"
	//"github.com/globaldce/globaldce-gateway/applog"
)

var (
	MainwalletFilePath = ""
	MainwalletFileKey  []byte
	AppPath            = ""
	AppName            = ""
	AppLocalPort       = ""
	//Walletinstantiated=false
	//Walletstate=""
	Miningrequested            = false
	Miningrunning              = false
	Miningaddrressesfileloaded = false
	//Miningaddressesloaded=false
	Managingwalletrequested = false
	Seed                    = false
	SyncingMinNbPeers       = 0
	AppIsClosing            = false
	Mn                      *mainchain.Maincore
	Mncc                    *content.ContentClient
	Wireswarm               *wire.Swarm
	Wlt                     *wallet.Wallet
	MAddresses              *MiningAddresses
)

func Walletinstantiated() bool {
	if Wlt == nil {
		return false
	}
	return Wlt.Walletloaded
}
