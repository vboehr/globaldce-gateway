package daemon

import (
	//"crypto/sha256"
	//"math/big"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/globaldce/globaldce-gateway/utility"
	"math/rand"

	//"github.com/globaldce/globaldce-gateway/applog"
	"fmt"
	"os"
)

//TODO integrate address  in the wallet
type MiningAddresses struct {
	Addressesarray []utility.Hash
}

func (maddresses *MiningAddresses) AddAddress(addr utility.Hash) {
	maddresses.Addressesarray = append(maddresses.Addressesarray, addr)
} //(rand.Intn(113)
func (maddresses *MiningAddresses) GetAddress(i int) utility.Hash {
	return maddresses.Addressesarray[i]
}
func (maddresses *MiningAddresses) GetRandomAddress() utility.Hash {
	maxnb := len(maddresses.Addressesarray)
	i := rand.Intn(maxnb)
	//fmt.Printf("****%d",i)
	return maddresses.Addressesarray[i]
} //(rand.Intn(113)

func (maddresses *MiningAddresses) LoadJSONMiningAddressesFile(filepath string) error {
	f, err := os.OpenFile(filepath, os.O_RDONLY, 0755)
	if err != nil {
		//log.Fatal(err)
		fmt.Println("error:", err)
	}
	defer f.Close()
	bufferMiningAddrfileseize := make([]byte, 4)
	_, rserr := f.Read(bufferMiningAddrfileseize)
	if rserr != nil {
		//log.Fatal(err)
		fmt.Println("error:", rserr)
	}
	var miningaddrfileseize uint32
	readerMiningAddrfileseize := bytes.NewReader(bufferMiningAddrfileseize)
	binary.Read(readerMiningAddrfileseize, binary.LittleEndian, &miningaddrfileseize)

	addrfilerawbytes := make([]byte, miningaddrfileseize)
	_, rerr := f.Read(addrfilerawbytes)
	if rerr != nil {
		//log.Fatal(err)
		fmt.Println("error:", rerr)
	}
	uerr := json.Unmarshal(addrfilerawbytes, maddresses)
	if uerr != nil {
		fmt.Println("error:", uerr)
		return uerr
	}
	return nil
}

func (maddresses *MiningAddresses) SaveJSONMiningAddressesFile(addrfilepath string) {
	addrfilerawbytes, err := json.Marshal(*maddresses)
	if err != nil {
		fmt.Println("error:", err)
	}

	bufferMiningAddrfileseize := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufferMiningAddrfileseize, uint32(len(addrfilerawbytes)))

	f, err := os.OpenFile(addrfilepath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		//log.Fatal(err)
		fmt.Println("error:", err)
	}
	defer f.Close()
	_, wserr := f.Write(bufferMiningAddrfileseize)
	if wserr != nil {
		//log.Fatal(err)
		fmt.Println("error:", wserr)
	}
	_, werr := f.Write(addrfilerawbytes)
	if werr != nil {
		//log.Fatal(err)
		fmt.Println("error:", werr)
	}
}

/*
   address,addrerr:=hex.DecodeString(addrstring)
   if addrerr!=nil{
       applog.Warning("\nError: inappropriate address provided - %v",err)
       os.Exit(0)
   }
*/
