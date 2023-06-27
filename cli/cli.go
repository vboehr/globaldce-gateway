package cli

import (
    "github.com/globaldce/globaldce-gateway/applog"
    "os"
    "github.com/globaldce/globaldce-gateway/wallet"
    "time"
    "fmt"
    "strings"
    "path/filepath"
    "github.com/globaldce/globaldce-gateway/daemon"
    "github.com/globaldce/globaldce-gateway/utility"
)

//var appName string
func Start(tmpappname string,tmpappversion string){
    

    fmt.Printf("%s %s \n",tmpappname,tmpappversion)
    //
    daemon.AppName=tmpappname
    //daemon.AppPath=tmpapppath
    applog.Init(daemon.AppPath)
    //for i:=0;i<len(os.Args);i++{
    //    applog.Notice("arg %d %s",i,os.Args[i])
    //}
    settingserr:=daemon.LoadUsersettingsFile()
	if settingserr!=nil{
		//
		daemon.SetDefaultSettings()
	}
	daemon.ApplyUsersettings()

    //appName=cliname

    if len(os.Args)<2{
        emptyCMD()
    }

 
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
        /*
        case "addbannedname":
            applog.Notice("User requested a bannedname")
            daemon.Usersettings.BannedNameArray=append(daemon.Usersettings.BannedNameArray,os.Args[2])
            _=daemon.SaveUsersettingsFile()
            os.Exit(0)
        */
        default:
            emptyCMD()
    }


    applog.Notice("")
    InterpretOptions()
    daemon.Wlt=Loadusermainwalletfile()
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
    //daemon.MainInit()
    go daemon.Mainloop()
    for {
        time.Sleep(3*time.Second)
    }
}

////////////////////////////////////////////

func Loadusermainwalletfile() *wallet.Wallet{
    if daemon.MainwalletFilePath==""{
        daemon.MainwalletFilePath= askuserwalletfilepath()
        applog.Notice("wallet file path set %s",daemon.MainwalletFilePath)
    }
    //
   
    

for (!daemon.Walletinstantiated()) && (!daemon.Miningaddrressesfileloaded){    
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
            wlt:=new(wallet.Wallet)
            lerr:=wlt.LoadJSONWalletFile(daemon.MainwalletFilePath,daemon.MainwalletFileKey)
            //daemon.Walletinstantiated=true
            if lerr==nil{
                //daemon.Walletinstantiated=true
                wlt.Walletloaded=true
            } else {
                applog.Notice("wallet file %s not loaded",wlt.Path)
            }
       // }
        
	} else {
        applog.Notice("\nwalletfile %s does not exist.",daemon.MainwalletFilePath)
        applog.LockDisplay()
        fmt.Printf("Do you want to create a new wallet at %s (default: no) :",daemon.MainwalletFilePath)
        var useranswer string
      
        fmt.Scanln(&useranswer)
        applog.UnlockDisplay()
        if strings.ToLower(useranswer)=="yes"{
            //TODO support for othe wallet types
            //
            applog.LockDisplay()
            wlt:=Newsequentialwallet()

            //
            daemon.MainwalletFileKey=createuserwalletfilekey()
    
            applog.UnlockDisplay()

            return wlt
        }
        daemon.MainwalletFilePath= askuserwalletfilepath()
    }
}
/*
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
        //    wlt.LoadJSONWalletFile(daemon.MainwalletFilePath,daemon.MainwalletFileKey)
        //    daemon.Walletinstantiated=true
        //}
        
	} else {
        applog.Notice("hotaddress file %s does not exist.",wlt.Path)
        wlt.Path= askuserwalletfilepath()
    }
}
*/
    //return wlt 
    return nil
}
func Newsequentialwallet()  *wallet.Wallet{

    
    
    fmt.Printf("\nCreating new sequential wallet \n")

    randomseedString,_:=wallet.GenerateMnemonicSeedString()
    fmt.Printf("Random Mnemonic Seed String :%s\n",randomseedString)
    seedStringCorrectlyEntred:=false
    var seedString string
    var reentredseedString string
    for (!seedStringCorrectlyEntred){
        fmt.Printf("Enter Seed String (default: random seed string will be used) :")
        
        
        fmt.Scanln(&seedString)
        fmt.Printf("Please Reenter Seed String (default: random seed string will be used) :")
        
        fmt.Scanln(&reentredseedString)
        if seedString==reentredseedString{
            seedStringCorrectlyEntred=true
        } else{
            fmt.Printf("Entred Seed Strings Do not match\n")
        }
    }


    if seedString==""{
        seedString=randomseedString
    }
    ///////////////////////////////////
    wlt:=wallet.Newsequentialwallet(seedString)

    ////////////////////////////////////
    return wlt    
}
func createuserwalletfilekey() []byte{
    var keystr string
    var reentredkeystr string
    for {
        
        fmt.Printf("\nPlease enter the walletfile key: ")
        fmt.Scanln(&keystr)
        if (len(keystr)<8){
            fmt.Printf("\nEntered key must be at least 8 characters length ")
        } else{
            
            fmt.Printf("\nPlease reenter the walletfile key: ")
            fmt.Scanln(&reentredkeystr)
            if reentredkeystr==keystr{
                break
            } else{
                fmt.Printf("\nEntered keys do not match ")
            }
        }

        

    }

    //if (len(key)==32){
        // No hashing is needed
    //    return key
    //} else {
        // If the key length is 32, the key is hashed first
        return utility.ComputeHashBytes([]byte(keystr))
    //}

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

    //if (len(key)==32){
        // No hashing is needed
    //    return key
    //} else {
        // If the key length is 32, the key is hashed first
        return utility.ComputeHashBytes(key)
    //}

}
func askuserwalletfilepath() string{
    var entredpath string
    fmt.Printf("\nPlease enter wallet file path (default: %s) :",filepath.Join(daemon.AppPath,daemon.MainwalletFilePathDefault))
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

