package wallet
import
(
	//"crypto/sha256"
	//"math/big"
	"github.com/globaldce/go-globaldce/utility"
	"encoding/json"
	"encoding/binary"
	"bytes"
	"math/rand"
	//"github.com/btcsuite/btcd/btcec/v2"
	//"github.com/globaldce/go-globaldce/applog"
	"fmt"
	"os"
)
//TODO integrate address book in the wallet
type HotAddresses struct {
	Addrarray [] utility.Hash
}

func (hota* HotAddresses) AddAddress(addr utility.Hash){
	hota.Addrarray=append(hota.Addrarray,addr)
}//(rand.Intn(113)
func (hota* HotAddresses) GetRandomAddress() utility.Hash{
	maxnb:=len(hota.Addrarray)
	i:=rand.Intn( maxnb )
	//fmt.Printf("****%d",i)
	return hota.Addrarray[i]
}//(rand.Intn(113)

func (hota *HotAddresses) LoadJSONFile(addrbookfilepath string) error{
	f, err := os.OpenFile(addrbookfilepath, os.O_RDONLY, 0755)
	if err != nil {
		//log.Fatal(err)
		fmt.Println("error:", err)
	}
	bufferAddrbookfileseize := make([]byte, 4)
	_, rserr := f.Read(bufferAddrbookfileseize)
	if rserr != nil {
		//log.Fatal(err)
		fmt.Println("error:", rserr)
	}
	var addrbookfileseize uint32
	readerAddrbookfileseize := bytes.NewReader(bufferAddrbookfileseize)
	binary.Read(readerAddrbookfileseize, binary.LittleEndian, &addrbookfileseize)

	addrbookfilerawbytes := make([]byte, addrbookfileseize)
	_, rerr := f.Read(addrbookfilerawbytes)
	if rerr != nil {
		//log.Fatal(err)
		fmt.Println("error:", rerr)
	}
	uerr:=json.Unmarshal(addrbookfilerawbytes,hota)
	if uerr != nil {
		fmt.Println("error:", uerr)
		return uerr
	}
	return nil
}

func (hota *HotAddresses) SaveJSONFile(addrbookfilepath string){
	addrbookfilerawbytes, err := json.Marshal(*hota)
	if err != nil {
		fmt.Println("error:", err)
	}
	
	bufferAddrbookfileseize := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufferAddrbookfileseize, uint32(len(addrbookfilerawbytes)))

	f, err := os.OpenFile(addrbookfilepath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		//log.Fatal(err)
		fmt.Println("error:", err)
	}
	_, wserr := f.Write(bufferAddrbookfileseize)
	if wserr != nil {
		//log.Fatal(err)
		fmt.Println("error:", wserr)
	}
	_, werr := f.Write(addrbookfilerawbytes)
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