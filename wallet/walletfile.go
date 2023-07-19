package wallet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/globaldce/globaldce-gateway/utility"
	//"github.com/btcsuite/btcd/btcec/v2"
	"fmt"
	"github.com/globaldce/globaldce-gateway/applog"
	"os"
	//"sync"
)

type PrivateKeyBytes []byte
type Walletfile struct {
	//Path string
	//Chain []byte
	//
	Privatekeyarray []PrivateKeyBytes

	Assetarray         []Asset
	Lastknownblock     uint64
	Broadcastedtxarray []Broadcastedtx
	Groupnamearray     []string
	Contactarray       []Contact
}

func (wlt *Wallet) SaveJSONWalletFile(walletfilepath string, key []byte) {

	walletfile := new(Walletfile)
	//walletfile.Chain=wlt.Chain

	for i := 0; i < len(wlt.Privatekeyarray); i++ {
		walletfile.Privatekeyarray = append(walletfile.Privatekeyarray, wlt.Privatekeyarray[i].Serialize())
	}

	walletfile.Assetarray = wlt.Assetarray
	walletfile.Lastknownblock = wlt.Lastknownblock
	walletfile.Broadcastedtxarray = wlt.Broadcastedtxarray
	walletfile.Groupnamearray = wlt.Groupnamearray
	walletfile.Contactarray = wlt.Contactarray
	walletfilerawstring, err := json.Marshal(walletfile)
	if err != nil {
		fmt.Println("error:", err)
	}
	var walletfilestring []byte
	if len(key) != 0 {
		walletfilestring, _ = utility.Encrypt(key, walletfilerawstring)
	} else {
		fmt.Println("No encryption key set")
		walletfilestring = walletfilerawstring
	}
	bufferWalletfiletype := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufferWalletfiletype, uint32(WALLET_TYPE_SEQUENTIAL)) // type 1

	//
	bufferWalletfileversion := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufferWalletfileversion, uint32(WALLET_TYPE_SEQUENTIAL_VERSION)) // type 1
	//

	bufferWalletfileseize := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufferWalletfileseize, uint32(len(walletfilestring)))
	//applog.Trace("  walletfile: %s", walletfilestring)

	f, err := os.OpenFile(walletfilepath, os.O_WRONLY|os.O_CREATE, 0755)

	if err != nil {
		//log.Fatal(err)
		fmt.Println("error:", err)
	}

	_, wterr := f.Write(bufferWalletfiletype)
	if wterr != nil {
		//log.Fatal(err)
		fmt.Println("error:", wterr)
	}
	//
	_, wverr := f.Write(bufferWalletfileversion)
	if wverr != nil {
		//log.Fatal(err)
		fmt.Println("error:", wverr)
	}
	//
	_, werr := f.Write(bufferWalletfileseize)
	if werr != nil {
		//log.Fatal(err)
		fmt.Println("error:", werr)
	}
	_, wserr := f.Write(walletfilestring)
	if wserr != nil {
		//log.Fatal(err)
		fmt.Println("error:", wserr)
	}
	applog.Notice("Wallet saved.")

}
func (wlt *Wallet) LoadJSONWalletFile(path string, key []byte) error {
	//walletfilerawstring:=*LoadJSONFile(path)
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		//log.Fatal(err)
		fmt.Println("error:", err)
	}
	/////////////////////////////
	bufferWalletfiletype := make([]byte, 4)
	_, rterr := f.Read(bufferWalletfiletype)
	if rterr != nil {
		//log.Fatal(err)
		fmt.Println("error:", rterr)
	}
	var walletfiletype uint32
	readerWalletfiletype := bytes.NewReader(bufferWalletfiletype)

	binary.Read(readerWalletfiletype, binary.LittleEndian, &walletfiletype)
	fmt.Println("type:", walletfiletype)
	////////////////////////////
	bufferWalletfileversion := make([]byte, 4)
	_, rverr := f.Read(bufferWalletfiletype)
	if rverr != nil {
		//log.Fatal(err)
		fmt.Println("error:", rverr)
	}
	var walletfileversion uint32
	readerWalletfileversion := bytes.NewReader(bufferWalletfileversion)

	binary.Read(readerWalletfileversion, binary.LittleEndian, &walletfileversion)
	fmt.Println("version:", walletfileversion)
	////////////////////////////

	bufferWalletfileseize := make([]byte, 4)
	_, rserr := f.Read(bufferWalletfileseize)
	if rserr != nil {
		//log.Fatal(err)
		fmt.Println("error:", rserr)
	}
	var walletfileseize uint32
	readerWalletfileseize := bytes.NewReader(bufferWalletfileseize)

	binary.Read(readerWalletfileseize, binary.LittleEndian, &walletfileseize)
	fmt.Println("seize:", walletfileseize)
	walletfilerawstring := make([]byte, walletfileseize)
	_, rerr := f.Read(walletfilerawstring)
	if rerr != nil {
		//log.Fatal(err)
		fmt.Println("reading error:", rerr)
	}

	var walletfilestring []byte
	fmt.Println("key length", len(key))
	if len(key) != 0 {
		walletfilestring, _ = utility.Decrypt(key, walletfilerawstring)
	} else {
		fmt.Println("No encryption key set")
		fmt.Println("data:", walletfilerawstring)
		walletfilestring = walletfilerawstring
	}

	//applog.Trace("read JSONFILE CONTENT: %s", walletfilestring)
	walletfile := new(Walletfile)
	uerr := json.Unmarshal(walletfilestring, walletfile)
	if uerr != nil {
		fmt.Println("unmarshal error:", uerr)
		return uerr
	}
	//applog.Trace("read JSONFILE CONTENT: %d %d", len (walletfile.Keypairarray),len(walletfile.Assetarray))
	for i := 0; i < len(walletfile.Privatekeyarray); i++ {

		privKey := utility.PrivKeyFromBytes([]byte(walletfile.Privatekeyarray[i]))
		wlt.Privatekeyarray = append(wlt.Privatekeyarray, &privKey)
	}
	wlt.Assetarray = walletfile.Assetarray
	wlt.Lastknownblock = walletfile.Lastknownblock
	wlt.Path = path
	wlt.Broadcastedtxarray = walletfile.Broadcastedtxarray
	wlt.Contactarray = walletfile.Contactarray
	wlt.Walletloaded = true

	return nil
}
