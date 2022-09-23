package daemon
import
(
	//"crypto/sha256"
	//"math/big"
	"github.com/globaldce/globaldce-gateway/utility"
	"encoding/json"
	"encoding/binary"
	"bytes"
	"math/rand"
	
	//"github.com/globaldce/globaldce-gateway/applog"
	"fmt"
	"os"
)
//TODO integrate address book in the wallet
type MiningAddresses struct {
	Addressesarray [] utility.Hash
}

func (maddresses* MiningAddresses) AddAddress(addr utility.Hash){
	maddresses.Addressesarray=append(maddresses.Addressesarray,addr)
}//(rand.Intn(113)
func (maddresses* MiningAddresses) GetRandomAddress() utility.Hash{
	maxnb:=len(maddresses.Addressesarray)
	i:=rand.Intn( maxnb )
	//fmt.Printf("****%d",i)
	return maddresses.Addressesarray[i]
}//(rand.Intn(113)

func (maddresses *MiningAddresses) LoadJSONFile(addrbookfilepath string) error{
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
	uerr:=json.Unmarshal(addrbookfilerawbytes,maddresses)
	if uerr != nil {
		fmt.Println("error:", uerr)
		return uerr
	}
	return nil
}

func (maddresses *MiningAddresses) SaveJSONFile(addrbookfilepath string){
	addrbookfilerawbytes, err := json.Marshal(*maddresses)
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