package wallet

import (
	//"crypto/sha256"
	//"math/big"
	"github.com/globaldce/globaldce-gateway/utility"
	//"encoding/json"
	//"encoding/binary"
	//"bytes"
	//"math/rand"
	//"github.com/btcsuite/btcd/btcec/v2"
	//"github.com/globaldce/globaldce-gateway/applog"
	"fmt"
	"strings"
)

func (wlt *Wallet) GetRegisteredNames() []string {
	var registerednames []string
	for _, asset := range wlt.Assetarray {

		if strings.Index(asset.StateString, "NAMEREGISTERED_") == 0 {
			r := strings.NewReplacer("NAMEREGISTERED_", "")
			tmpstatestring := asset.StateString
			tmpregisteredname := r.Replace(tmpstatestring)
			registerednames = append(registerednames, tmpregisteredname)
			fmt.Printf("Registred Name %s\n", tmpregisteredname)
		}

	}
	return registerednames
}
func (wlt *Wallet) GetAddressesDetails() []string {
	var addresses []string
	//fmt.Printf("\nNumber of addresses %d\n",len(wlt.Privatekeyarray))
	for _, prvkey := range wlt.Privatekeyarray {
		address := utility.ComputeHash(prvkey.PubKey().SerializeCompressed())
		addresses = append(addresses, fmt.Sprintf("%x", address))
	}
	return addresses

}
func (wlt *Wallet) GetAssetsDetails() []string {
	var assestsdestails []string
	//fmt.Printf("\nNumber of addresses %d\n",len(wlt.Privatekeyarray))
	for _, tmpasset := range wlt.Assetarray {

		assestsdestails = append(assestsdestails, fmt.Sprintf(" %f ", (float64(tmpasset.Value)/1000000.0))+tmpasset.StateString)
	}
	return assestsdestails

}
