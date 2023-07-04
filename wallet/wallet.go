package wallet

import (
	"github.com/globaldce/globaldce-gateway/utility"
	//"encoding/json"
	//"encoding/binary"
	//"bytes"
	"github.com/btcsuite/btcd/btcec/v2"
	//"github.com/globaldce/globaldce-gateway/applog"
	"fmt"
	//"os"
	//"time"
	"sync"
)

const (
	NB_INITIAL_HASHES int = 10000000

	WALLET_TYPE_SEQUENTIAL         uint32 = 1
	WALLET_TYPE_SEQUENTIAL_VERSION uint32 = 1
)

//
type Wallet struct {
	//HotWallet bool
	//Hotaddresses HotAddresses
	//Chain []byte
	Walletloaded        bool
	Walletstate         string
	Path                string
	Type                uint32
	Version             uint32
	Privatekeyarray     []*btcec.PrivateKey
	Assetarray          []Asset
	Lastknownblock      uint64
	Broadcastedtxarray  []Broadcastedtx
	Commcredentialarray []Commcredential
	Groupnamearray      []string
	Contactarray        []Contact

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

type Broadcastedtx struct {
	Tx                 utility.Transaction
	ConfirmationString string
}
type Asset struct {
	Hash            utility.Hash
	Index           uint32
	Value           uint64
	Privatekeyindex uint32
	StateString     string
	SpendingTxHash  utility.Hash
	SpendingIndex   uint32
}

func Newsequentialwallet(seedString string) *Wallet {
	wlt := new(Wallet)
	wlt.Walletloaded = false
	go Gensequentialwallet(wlt, seedString)
	return wlt
}
func Gensequentialwallet(wlt *Wallet, seedString string) {
	wlt.Type = WALLET_TYPE_SEQUENTIAL
	wlt.Version = WALLET_TYPE_SEQUENTIAL_VERSION
	InitialHashBytes := []byte(seedString)
	for i := 0; i < NB_INITIAL_HASHES; i++ {
		InitialHashBytes = utility.ComputeHashBytes(InitialHashBytes)
		wltgenprogress := int(i * 100 / NB_INITIAL_HASHES)
		if (i)%(NB_INITIAL_HASHES/10) == 0 {
			wlt.Walletstate = fmt.Sprintf("Wallet Generation Progress %d %%\n", wltgenprogress)
			fmt.Printf("Wallet Generation Progress %d %%\n", wltgenprogress)
		}

	}
	fmt.Printf("%x\n", InitialHashBytes)
	pk := utility.PrivKeyFromBytes(InitialHashBytes)
	wlt.Privatekeyarray = append(wlt.Privatekeyarray, &pk)
	wlt.Walletstate = ""
	wlt.Walletloaded = true
}

/*
func (wlt *Wallet) GetAssetInfo(int i) ([]byte,value){
	return len(wlt.Assetarray)
}*/
func (wlt *Wallet) GetNbAssets() int {
	return len(wlt.Assetarray)
}

func (wlt *Wallet) GetLastAddress() utility.Hash {
	//if wlt.HotWallet {
	//	return wlt.Hotaddresses.GetRandomAddress()
	//}
	//if len(wlt.Privatekeyarray)==0{
	//	return wlt.GenerateKeyPair()
	//}
	pk := wlt.Privatekeyarray[len(wlt.Privatekeyarray)-1]
	return utility.ComputeHash(pk.PubKey().SerializeCompressed())
}
func (wlt *Wallet) GetAddress(i uint) utility.Hash {
	if i > uint(len(wlt.Privatekeyarray)) {
		return *utility.NewHash(nil)
	}
	pk := wlt.Privatekeyarray[i]
	return utility.ComputeHash(pk.PubKey().SerializeCompressed())
}
func (wlt *Wallet) GetPrivatekeyindexFromAddress(addr utility.Hash) int {
	for privkeyindex, privkey := range wlt.Privatekeyarray {
		pubkey := privkey.PubKey().SerializeCompressed()
		if addr == utility.ComputeHash(pubkey) {
			return privkeyindex
		}
	}
	return -1
}

func (wlt *Wallet) GetNbAddresses() int {
	return len(wlt.Privatekeyarray)
}

/*


 */
func (wlt *Wallet) GenerateKeyPairs(nbkeypair int) {
	for i := 0; i < nbkeypair; i++ {
		wlt.GenerateKeyPair()

	}
}
func (wlt *Wallet) GenerateKeyPair() utility.Hash {

	//message:="message text"
	//messageHash:=ComputeHashBytes([]byte(message))
	//TODO support for other types of wallet
	prevpk := wlt.Privatekeyarray[len(wlt.Privatekeyarray)-1]
	tmppkbytes := utility.ComputeHashBytes(prevpk.Serialize())
	pk := utility.PrivKeyFromBytes(tmppkbytes)
	//if err != nil {
	//	//il, err//return (*PrivateKey)(key), nil
	//	applog.Trace("err: %x", err)
	//}
	// pk.Serialize() 							returns a 32 bytes private key
	// pk.PubKey().SerializeUncompressed() 		returns a 65 bytes public key
	// pk.PubKey().SerializeCompressed() 		returns a 33 bytes public key

	wlt.Privatekeyarray = append(wlt.Privatekeyarray, &pk)
	//applog.Trace("total length of keys %d",len(wlt.Privatekeyarray))
	//applog.Trace("private key: %d public key: %d OR %d ", len(pk.Serialize()), len(pk.PubKey().SerializeUncompressed()), len(pk.PubKey().SerializeCompressed()))

	//sig, err := pk.SignCompact(messageHash)
	//if err != nil {
	//	applog.Trace("error")
	//	//return 0
	//}
	//applog.Trace("signature %x", sig)

	// Verify the signature for the message using the public key.
	//verified := sig.Verify(messageHash, pk.PubKey())
	//applog.Trace("Signature Verified? %v", verified)

	return utility.ComputeHash(pk.PubKey().SerializeCompressed())
}

func (wlt *Wallet) AddAsset(txhash utility.Hash, index uint32, value uint64, privkeyindex uint32, assetstate string) {
	var emptytxhash utility.Hash
	tmpasset := Asset{
		Hash:            txhash,
		Index:           index,
		Value:           value,
		Privatekeyindex: privkeyindex,
		StateString:     assetstate, //"UNSPENT",
		SpendingTxHash:  emptytxhash,
		SpendingIndex:   999999,
	}
	wlt.Assetarray = append(wlt.Assetarray, tmpasset)
}
func (wlt *Wallet) GetAssetFromRegisteredName(name string) (*Asset, error) {

	for _, a := range wlt.Assetarray {
		if a.StateString == "NAMEREGISTERED_"+name {
			return &a, nil
		}
	}
	return nil, fmt.Errorf("No associated registered name was found in this wallet")
}
func (wlt *Wallet) AddBroadcastedtx(tx utility.Transaction) {

	btx := Broadcastedtx{
		Tx:                 tx,
		ConfirmationString: "",
	}
	wlt.Broadcastedtxarray = append(wlt.Broadcastedtxarray, btx)
}
