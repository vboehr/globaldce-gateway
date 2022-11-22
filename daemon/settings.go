package daemon

import(
	"fmt"
	"encoding/json"
	"path/filepath"
	"github.com/globaldce/globaldce-gateway/utility"
	//"github.com/globaldce/globaldce-gateway/mainchain"
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
	CachedDirInfoArray []CachedDirInfo
	Recentdappnamesarray []string
	Activeloginname string
}
var Usersettings UsersettingsType

func GetActiveloginname() string{
	registerednames:=Wlt.GetRegisteredNames()
	for _,regname:=range registerednames{
		if regname==Usersettings.Activeloginname{
			return Usersettings.Activeloginname
		}
	}
	return ""
}

func PutActiveloginname(newlogin string) {
	Usersettings.Activeloginname=newlogin
}

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
		CachedDirInfoArray:nil,
		Activeloginname:"",
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
	usersettingsfilebytes,lerr:=utility.LoadBytesFile("settings.ini")
	if lerr != nil {
		fmt.Println("error:", lerr)
		return lerr
	}
	uerr:=json.Unmarshal(*usersettingsfilebytes,&Usersettings)
	if uerr != nil {
		fmt.Println("error:",uerr)
		return uerr
	}

	return nil
}
type CachedDirInfo struct {
	RegistredName string
	Path string
}
func GetCachedDirPathForRegistredName(name string) (string,int){
		for i , cacheddirinfo:= range Usersettings.CachedDirInfoArray{
			if cacheddirinfo.RegistredName==name{
				//fmt.Println("********************", cacheddirinfo.Path,i)
				return cacheddirinfo.Path,i
			}
		}
		return "",-1
}
func PutCachedDirPathForRegistredName(name string,path string) {
	p,i:=GetCachedDirPathForRegistredName(name)
	if p==""{
		newcacheddirinfo:=new(CachedDirInfo)
		newcacheddirinfo.RegistredName=name
		newcacheddirinfo.Path=path
		Usersettings.CachedDirInfoArray=append(Usersettings.CachedDirInfoArray,*newcacheddirinfo)
	} else{
		Usersettings.CachedDirInfoArray[i].Path=path
	}

}
func AddToRecentDappNames(dappnameinputText string) {
	for _ , tmpname:= range Usersettings.Recentdappnamesarray{
		if tmpname==dappnameinputText{
			return 
		}
	}
	Usersettings.Recentdappnamesarray=append([]string{dappnameinputText}, Usersettings.Recentdappnamesarray ...)
	return
}
func GetRecentDappNames() ([]string){
	return Usersettings.Recentdappnamesarray
}
func ClearRecentDappNameWithId(dappnameselectedid int){
	if (dappnameselectedid<0)||(dappnameselectedid>=len(Usersettings.Recentdappnamesarray)){
		return
	}
	Usersettings.Recentdappnamesarray=append(Usersettings.Recentdappnamesarray[:dappnameselectedid], Usersettings.Recentdappnamesarray[dappnameselectedid+1:] ...)
}
func ClearAllRecentDappNames(){
	Usersettings.Recentdappnamesarray=make([]string, 0)
}
