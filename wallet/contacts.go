
package wallet
import
(
	"github.com/globaldce/go-globaldce/utility"
	//"encoding/json"
	//"encoding/binary"
	//"bytes"
	//"github.com/btcsuite/btcd/btcec"
	//"github.com/globaldce/go-globaldce/applog"
	//"fmt"
	//"os"
	//"sync"
)


type Contact struct{
	Name string
	Address utility.Hash
}
func (wlt *Wallet) AddContact(name string,hash utility.Hash) {
	//var emptytxhash utility.Hash
	tmpcontact:=Contact{
		Name:name,
		Address:hash,
	}
	wlt.Contactarray=append(wlt.Contactarray,tmpcontact)
}