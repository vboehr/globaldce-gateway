package daemon

import(
	"fmt"
	"encoding/json"
	"path/filepath"
	"github.com/globaldce/go-globaldce/utility"
	//"github.com/globaldce/go-globaldce/daemon"
)
var MainwalletFilePathDefault=filepath.Join("WalletFiles","Wallet.dat")
var (
	MaxDisplayedPostDefault=30
	NameregistrationtxfeeDefault=300
	//PublicposttxfeeDefault=500
	SendtoaddressarraytxfeeDefault=500
)

type UsersettingsType struct {
	//CloudWallet bool
	//Cloudaddresses CloudAddr
	//Chain []byte
	//Privatekeyarray [] *btcec.PrivateKey
	//Assetarray [] Asset
	MainwalletFilePath string
	//////////////////////////
	//Confirmationlayer uint32
	//////////////////////////
	MaxDisplayedPost int
	Nameregistrationtxfee int
	//Publicposttxfee int
	Sendtoaddressarraytxfee int
	//Lastknownblock uint32
	//Broadcastedtxarray [] Broadcastedtx
	//mu sync.Mutex
	//BannedNameArray []string
	Miningrequested bool
}
var Usersettings UsersettingsType
	
func SetDefaultSettings(){
	
	Usersettings=UsersettingsType{
		MainwalletFilePath:MainwalletFilePathDefault,
		/////////////////////////////////
		//Confirmationlayer: uint32(200),
		/////////////////////////////////
		MaxDisplayedPost:MaxDisplayedPostDefault,
		Nameregistrationtxfee:NameregistrationtxfeeDefault,
		//Publicposttxfee:PublicposttxfeeDefault,
		Sendtoaddressarraytxfee:SendtoaddressarraytxfeeDefault,
		//BannedNameArray:nil,
		Miningrequested:false,
	}
}

func ApplyUsersettings(){

	MainwalletFilePath=Usersettings.MainwalletFilePath//MainwalletFilePathDefault//"./WalletFiles/MainWallet.dat"
	Miningrequested=Usersettings.Miningrequested
}

//func GetMainwalletPath() string{
//	return Usersettings.MainwalletFilePath
//}
//func PutMainwalletPath(path string){
//	Usersettings.MainwalletFilePath=path
//}
func SaveUsersettingsFile() error{
	//if Mn!=nil{
	//	Usersettings.BannedNameArray=Mn.BannedNameArray
	//}
	
	//
	usersettingsfilebytes, err := json.Marshal(Usersettings)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	_=utility.SaveBytesFile(usersettingsfilebytes,"settings.ini")
	fmt.Println("Usersettings saved.")
	return nil
}
func LoadUsersettingsFile() error{
	usersettingsfilebytes,_:=utility.LoadBytesFile("settings.ini")
	uerr:=json.Unmarshal(*usersettingsfilebytes,&Usersettings)
	if uerr != nil {
		fmt.Println("error:", uerr)
		return uerr
	}

	return nil
}
