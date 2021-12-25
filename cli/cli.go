package cli

import (
    "github.com/globaldce/globaldce-toolbox/applog"
    "os"
    "github.com/globaldce/globaldce-toolbox/wallet"
    "time"
    "fmt"
    "strings"
    "path/filepath"
    "github.com/globaldce/globaldce-toolbox/daemon"
    "github.com/globaldce/globaldce-toolbox/utility"
)
const (
    appVersion="0.1.0"
)
var appName string
func Start(cliname string){
    applog.Init()
    //for i:=0;i<len(os.Args);i++{
    //    applog.Notice("arg %d %s",i,os.Args[i])
    //}
    settingserr:=daemon.LoadUsersettingsFile()
	if settingserr!=nil{
		//
		daemon.SetDefaultSettings()
	}
	daemon.ApplyUsersettings()

    appName=cliname

    if len(os.Args)<2{
        emptyCMD()
    }
    //appPath=os.Args[0]
    appPathOS, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
            applog.Fatal("%v",err)
    }
    daemon.AppPath=appPathOS
    daemon.AppPath=""//TODO include appPathOS
    // common options
    ///////////////////////
    
    switch strings.ToLower(os.Args[1]) {
        case "help":
            helpCMD()
        case "h":
            helpCMD()
            //
        case "sync":
            applog.Notice("User requested syncing")
            //
        case "mine":
            applog.Notice("User requested mining")
            daemon.Miningrequested=true
            //
        case "managewallet":
            applog.Notice("User requested managing wallet")

            daemon.Managingwalletrequested=true
            //
        case "createnewaddress":
            applog.Notice("User requested a new address")
            createnewaddressCMD()
        case "addbannedname":
            applog.Notice("User requested a bannedname")
            daemon.Usersettings.BannedNameArray=append(daemon.Usersettings.BannedNameArray,os.Args[2])
            _=daemon.SaveUsersettingsFile()
            os.Exit(0)
        default:
            emptyCMD()
    }

    for i := 2; i<len(os.Args); i++ {
        tmparg:= os.Args[i]
        if strings.Index(tmparg, "-path=")==0{
            daemon.AppPath=strings.TrimPrefix(tmparg, "-path=")
            applog.Notice("appPath set to: %s",daemon.AppPath)
        }
        if strings.Index(tmparg, "-port=")==0{
            daemon.AppLocalPort=strings.TrimPrefix(tmparg, "-port=")
            applog.Notice("appLocalPort set to: %s",daemon.AppLocalPort)
        }
        if strings.Index(tmparg, "-seed")==0{
            daemon.Seed=true
        }
        if strings.Index(tmparg, "-hotmining")==0{
            daemon.HotMining=true
            applog.EnableDisplayTrace()
        }
        if strings.Index(tmparg, "-trace")==0{
            applog.EnableDisplayTrace()
        }

    }

    applog.Notice("")
    
    daemon.Wlt=loadusermainwalletfile()
    //_=wlt
    //
    //go daemon.Mainloop()

    //
    go func() {
        if ( daemon.Managingwalletrequested){
            time.Sleep(3*time.Second)
            applog.LockDisplay()
            var useranswer string
            fmt.Printf("Before managing the wallet initial syncing must be completed\n")
            fmt.Printf("Do you want to wait until the initial syncing is completed yes/no (default: yes): ")
            
            fmt.Scanln(&useranswer)
            applog.UnlockDisplay()
            if strings.ToLower(useranswer)=="no"{
                daemon.Mn.SyncWallet(daemon.Wlt)
                managewallet(daemon.Wireswarm,daemon.Mn,daemon.Wlt)
                daemon.Managingwalletrequested=false
            }
        }
        //////////////////////////////
        for {
            time.Sleep(3*time.Second)
            if  daemon.Managingwalletrequested &&  daemon.Wireswarm.Syncingdone {
                managewallet(daemon.Wireswarm,daemon.Mn,daemon.Wlt)
                daemon.Managingwalletrequested=false
            }

        }

        //////////////////////////////
    }()
    //
    daemon.Mainloop()
}

////////////////////////////////////////////

func loadusermainwalletfile() *wallet.Wallet{
    if daemon.MainwalletFilePath==""{
        daemon.MainwalletFilePath= askuserwalletfilepath()
        applog.Notice("wallet file path set %s",daemon.MainwalletFilePath)
    }
    //
   
    wlt:=new(wallet.Wallet)

for (!daemon.Walletloaded)&&(!daemon.HotMining){    
    if _, err := os.Stat( daemon.MainwalletFilePath); !os.IsNotExist(err) {
        
        //if daemon.HotMining{
        //    // TODO need better error handling
        //    // 
        //    wlt.HotWallet=true
        //    wlt.Path= askuserwalletfilepath()
        //    applog.Notice("wallet file path set %s",wlt.Path)
        //    _= wlt.Hotaddresses.LoadJSONFile(wlt.Path)
        //    
        //} else{
            // TODO better error handdling
            daemon.MainwalletFileKey=askuserwalletfilekey()
            lerr:=wlt.LoadJSONFile(daemon.MainwalletFilePath,daemon.MainwalletFileKey)
            //daemon.Walletloaded=true
            if lerr==nil{
                daemon.Walletloaded=true
            } else {
                applog.Notice("wallet file %s not loaded",wlt.Path)
            }
       // }
        
	} else {
        applog.Notice("walletfile %s does not exist.",daemon.MainwalletFilePath)
        daemon.MainwalletFilePath= askuserwalletfilepath()
    }
}
for (!wlt.HotWallet)&&(daemon.HotMining){    
    if _, err := os.Stat( wlt.Path); !os.IsNotExist(err) {
        
        //f daemon.HotMining{
            // TODO need better error handling
            // 
            
            //wlt.Path= askuserwalletfilepath()
            applog.Notice("hotaddresses file path set %s",wlt.Path)
            lerr:= wlt.Hotaddresses.LoadJSONFile(wlt.Path)
            if lerr==nil{
                wlt.HotWallet=true
            } else {
                applog.Notice("hotaddresses file %s not loaded",wlt.Path)
                wlt.Path= askuserwalletfilepath()
            }
            
        //} else{
        //    // TODO better error handdling
        //    daemon.MainwalletFileKey=askuserwalletfilekey()
        //    wlt.LoadJSONFile(daemon.MainwalletFilePath,daemon.MainwalletFileKey)
        //    daemon.Walletloaded=true
        //}
        
	} else {
        applog.Notice("hotaddress file %s does not exist.",wlt.Path)
        wlt.Path= askuserwalletfilepath()
    }
}
    return wlt 
}
func askuserwalletfilekey() []byte{
    var key []byte
    for {
        
        fmt.Printf("\nPlease enter the walletfile key: ")
        fmt.Scanln(&key)
        if (len(key)>=8){
            break
        }
        //if (len(key)==0){
        //    fmt.Printf("\nNo key no encryption \n")
        //    break
        //}		
        fmt.Printf("\nEntered key must be at least 8 characters length ")
    }

    if (len(key)==32){
        // No hashing is needed
        return key
    } else {
        // If the key length is 32, the key is hashed first
        return utility.ComputeHashBytes(key)
    }

}
func askuserwalletfilepath() string{
    var entredpath string
    fmt.Printf("Please enter wallet file path (default: %s) :",filepath.Join(daemon.AppPath,daemon.MainwalletFilePathDefault))
    fmt.Scanln(&entredpath)
    if entredpath!=""{
        return entredpath
    } else{
        walletfilesdirpath:=filepath.Join("./","WalletFiles")
        if _, err := os.Stat(walletfilesdirpath); os.IsNotExist(err) {
            os.Mkdir(walletfilesdirpath, os.ModePerm)
            //TODO better error handling
        }
        return filepath.Join(daemon.AppPath,daemon.MainwalletFilePathDefault)
    }

}

