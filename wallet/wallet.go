package wallet
import
(
	"github.com/globaldce/globaldce-toolbox/utility"
	//"encoding/json"
	//"encoding/binary"
	//"bytes"
	"github.com/btcsuite/btcd/btcec"
	"github.com/globaldce/globaldce-toolbox/applog"
	"fmt"
	//"os"
	"sync"
)
// 
type Wallet struct {
	HotWallet bool
	Hotaddresses HotAddresses
	Chain []byte
	Privatekeyarray [] *btcec.PrivateKey
	Assetarray [] Asset
	Path string
	Lastknownblock uint64
	Broadcastedtxarray [] Broadcastedtx
	Contactarray [] Contact
	mu sync.Mutex
}
func (wlt *Wallet) Lock() {
	//return mn.confirmationlayer
	wlt.mu.Lock()
}
func (wlt *Wallet) Unlock() {
	//return mn.confirmationlayer
	wlt.mu.Unlock()
}
type Broadcastedtx struct{
	Tx utility.Transaction
	ConfirmationString string
}
type Asset struct {
	Hash utility.Hash
	Index uint32
	Value uint64
	Privatekeyindex uint32
	StateString string
	SpendingTxHash utility.Hash
	SpendingIndex uint32
}
/*
func (wlt *Wallet) GetAssetInfo(int i) ([]byte,value){
	return len(wlt.Assetarray)
}*/
func (wlt *Wallet) GetNbAssets() int{
	return len(wlt.Assetarray)
}


func (wlt *Wallet) GetLastAddress() utility.Hash{
	if wlt.HotWallet {
		return wlt.Hotaddresses.GetRandomAddress()
	}
	if len(wlt.Privatekeyarray)==0{
		return wlt.GenerateKeyPair()
	}
	pk:=wlt.Privatekeyarray[len(wlt.Privatekeyarray)-1]
	return utility.ComputeHash(pk.PubKey().SerializeCompressed())
}
func (wlt *Wallet) GetAddress(i uint) utility.Hash{
	if i>uint (len(wlt.Privatekeyarray)){
		return *utility.NewHash(nil)
	}
	pk:=wlt.Privatekeyarray[i]
	return utility.ComputeHash(pk.PubKey().SerializeCompressed())
}
func (wlt *Wallet) GetNbAddresses() int{
	return len(wlt.Privatekeyarray)
}

func (wlt *Wallet) GenerateKeyPair() utility.Hash{

	//message:="message text"
	//messageHash:=ComputeHashBytes([]byte(message))

	pk, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		//il, err//return (*PrivateKey)(key), nil
		applog.Trace("err: %x", err)
	}
	// pk.Serialize() 							returns a 32 bytes private key
	// pk.PubKey().SerializeUncompressed() 		returns a 65 bytes public key
	// pk.PubKey().SerializeCompressed() 		returns a 33 bytes public key

	wlt.Privatekeyarray=append(wlt.Privatekeyarray,pk)
	//applog.Trace("total length of keys %d",len(wlt.Privatekeyarray))
	//applog.Trace("private key: %d public key: %d OR %d ", len(pk.Serialize()), len(pk.PubKey().SerializeUncompressed()), len(pk.PubKey().SerializeCompressed()))

	/*
	sig, err := pk.Sign(messageHash)
	if err != nil {
		applog.Trace("error")
		//return 0
	}
	applog.Trace("signature %x", sig.Serialize())

	// Verify the signature for the message using the public key.
	verified := sig.Verify(messageHash, pk.PubKey())
	applog.Trace("Signature Verified? %v", verified)
	*/
	return utility.ComputeHash(pk.PubKey().SerializeCompressed())
}
func (wlt *Wallet) AddAsset(txhash utility.Hash,index uint32,value uint64,privkeyindex uint32,assetstate string) {
	var emptytxhash utility.Hash
	tmpasset:=Asset{
		Hash:txhash,
		Index:index,
		Value:value,
		Privatekeyindex:privkeyindex,
		StateString:assetstate,//"UNSPENT",
		SpendingTxHash:emptytxhash,
		SpendingIndex:999999,
	}
	wlt.Assetarray=append(wlt.Assetarray,tmpasset)
}
func (wlt *Wallet) GetAssetFromRegisteredName(name string)(*Asset,error){

	for _,a:=range wlt.Assetarray {
		if a.StateString=="NAMEREGISTERED_"+name {
			return &a,nil
		}
	}
	return nil,fmt.Errorf("No associated registered name was found in this wallet")
}
func (wlt *Wallet) AddBroadcastedtx(tx utility.Transaction) {

	btx:=Broadcastedtx{
		Tx:tx,
		ConfirmationString:"",
	}
	wlt.Broadcastedtxarray=append(wlt.Broadcastedtxarray,btx)
}

