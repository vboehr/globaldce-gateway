package gui


import (
	//"os"
	//"log"
	"fmt"
	"errors"
	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	//"fyne.io/fyne/v2/data/validation"
	//"net/url"
	//"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2"
	"github.com/globaldce/globaldce-gateway/daemon"
	"github.com/globaldce/globaldce-gateway/utility"
	"github.com/globaldce/globaldce-gateway/cli"
	"github.com/globaldce/globaldce-gateway/wallet"
)


func passwordValidation() fyne.StringValidator {

	return func(text string) error {
		if len(text) <8 { //TODO set minimum length password
			return errors.New("Password length")
		}
		/////////////////////////////
		/////////////////////////////
		return nil // Nothing to validate with, same as having no validator.
	}
}

func  passwordDialog(win fyne.Window){
	//username := widget.NewEntry()
	//username.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")
	password := widget.NewPasswordEntry()
	//password.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "password can only contain letters, numbers, '_', and '-'")
	password.Validator=passwordValidation()
	//remember := false
	items := []*widget.FormItem{
		//widget.NewFormItem("Username", username),
		widget.NewFormItem("Password", password),
		//widget.NewFormItem("Remember me", widget.NewCheck("", func(checked bool) {
		//	remember = checked
		//})),
	}

	dialog.ShowForm("Password required to load wallet at: "+daemon.MainwalletFilePath, "Okay  ", "Cancel", items, func(b bool) {
		if !b {
			fmt.Println("canceled")
			nowalletFoundDialog(win,"")
			return
		}
		if b {
			text:=password.Text
				key:=[]byte(text)
				if (len(key)==32){
					// No hashing is needed
					daemon.MainwalletFileKey= key
				} else {
					// If the key length is 32, the key is hashed first
					daemon.MainwalletFileKey=utility.ComputeHashBytes(key)
				}
				lerr:=daemon.Wlt.LoadJSONWalletFile(daemon.MainwalletFilePath,daemon.MainwalletFileKey)
				if lerr!=nil{
					//nowalletFoundDialog(win,"Could not decrypt walletfile "+daemon.MainwalletFilePath)
					passwordDecryptionFailedDialog(win,"Could not decrypt walletfile "+daemon.MainwalletFilePath)
				} else {
					daemon.Wlt.Walletloaded=true
					daemon.Mn.SyncWallet(daemon.Wlt)
					//
				}
		}
		//var rememberText string
		//if remember {
		//	rememberText = "and remember this login"
		//}

		//fmt.Println("Entred password", password.Text)
	}, win)

}
func newWalletCreationDialog(win fyne.Window){
	selectWalletFileCallback := func(response bool){
		fmt.Println("Responded with", response)
		if response {
			newSequentialWalletSeedCreationDialog(win)
		} else {
			//
			nowalletFoundDialog(win,"")
		}
	}
	//TODO support other types of wallets
	//combo := widget.NewSelect([]string{"Sequential Wallet", "Option 2"}, func(value string) {
	combo := widget.NewSelect([]string{"Sequential Wallet"}, func(value string) {
		fmt.Println("Select set to", value)
	})
	combo.SetSelected("Sequential Wallet")
	content := container.NewVBox(widget.NewLabel("Wallet Type :"),combo)

	cnf:=dialog.NewCustomConfirm("Would you like to create a new wallet file ?", "Yes ", "No  ", content,selectWalletFileCallback, win)
	cnf.Show()
}


//
func newSequentialWalletPasswordCreationDialog(win fyne.Window,wltseedString string) {
	//
	wltpassword := widget.NewPasswordEntry()//NewMultiLineEntry()//
	wltpassword.Validator=passwordValidation()
	firstpassword := widget.NewPasswordEntry()
	//password.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "password can only contain letters, numbers, '_', and '-'")
	firstpassword.Validator=passwordValidation()
	secondpassword := widget.NewPasswordEntry()
	secondpassword.Validator=passwordValidation()

	remember:=widget.NewLabel("THIS PASSWORD WILL BE USED TO DECRYPT YOUR WALLET")


	items := []*widget.FormItem{

		widget.NewFormItem("Password", firstpassword),
		widget.NewFormItem("Password", secondpassword),
		widget.NewFormItem("", remember),
	}


	dialog.ShowForm("A password is required for the new wallet          ", "Okay  ", "Cancel", items, func(b bool) {
		if !b {
			fmt.Println("canceled")
			nowalletFoundDialog(win,"")
			return
		}
		//fmt.Println("Password",password.Text)

		if firstpassword.Text!=secondpassword.Text  {
			nowalletFoundDialog(win,"Wallet creation abroted entred passwords do not match")
		}
		if  b&&(firstpassword.Text==secondpassword.Text) {
			text:=firstpassword.Text
				key:=[]byte(text)
				if (len(key)==32){
					// No hashing is needed
					daemon.MainwalletFileKey= key
				} else {
					// If the key length is 32, the key is hashed first
					daemon.MainwalletFileKey=utility.ComputeHashBytes(key)
				}
				walletsettingsDisplayedMainwalletFilePath.Set(daemon.MainwalletFilePath)
				//fmt.Println("*** wltseed.Text",wltseedString,text)
				//os.Exit(0)
				newSequentialWalletCreationProgressDialog(win,wltseedString)
		} 
		//var rememberText string
		//if remember {
		//	rememberText = "and remember this login"
		//}
		//fmt.Println("Entred password", password.Text)
	}, win)


}
//
func newSequentialWalletSeedCreationDialog(win fyne.Window) {
	//

	fmt.Printf("\nCreating new sequential wallet \n")


	//
	//username := widget.NewEntry()
	//username.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")
	cli.GenerateNewMainwalletFilePath()

	wltseed := widget.NewPasswordEntry()//NewMultiLineEntry()//
	wltseed.MultiLine=true
	wltseed.Wrapping=fyne.TextWrapBreak
	wltseed.Validator=passwordValidation()

	wltseed.Text,_=wallet.GenerateMnemonicSeedString()
    fmt.Printf("Random Mnemonic Seed String :%s\n",wltseed.Text)
	//wltseed.SetMinSize(fyne.NewSize(100, 20))
	//wltseed.Validator=wltseedValidation()


	remember:=widget.NewLabel("CAUNTION: IF YOU LOSE YOUR SEED YOU CAN NOT RECOVER YOUR WALLET")
	generateSeedButton:=widget.NewButton("Generate New Seed", func() {
		fmt.Println("tapped Generate New Seed")
		wltseed.Text,_=wallet.GenerateMnemonicSeedString()
		wltseed.Refresh()
		})

	items := []*widget.FormItem{
		//widget.NewFormItem("Username", username),
		widget.NewFormItem("Seed", wltseed),
		widget.NewFormItem("", generateSeedButton),
		widget.NewFormItem("", remember),
	}


	//usersinfocontent = widget.NewForm(widget.NewFormItem("New label", scroll))
	//dialog.ShowForm("Password required to create a new wallet          ", "Okay  ", "Cancel", items, func(b bool) {
	dialog.ShowForm("A seed and a password are required to create a new wallet          ", "Okay  ", "Cancel", items, func(b bool) {
		if !b {
			fmt.Println("canceled")
			nowalletFoundDialog(win,"")
			return
		}
		//fmt.Println("Password",password.Text)

		if  b {

				walletsettingsDisplayedMainwalletFilePath.Set(daemon.MainwalletFilePath)
				fmt.Println("*** wltseed.Text",wltseed.Text)
				//os.Exit(0)
				newSequentialWalletPasswordCreationDialog(win,wltseed.Text)
		} 

	}, win)


}
///////////////////////////////////////////////////////////////////

func newSequentialWalletCreationProgressDialog(win fyne.Window,seedString string) {
	//wlt:=new(wallet.Wallet)        
    //wlt.Type=wallet.WALLET_TYPE_SEQUENTIAL
	//creationProgressDialog:=dialog.NewProgress("New Sequential Wallet", "Please stand by", win)
	wlt:=wallet.Newsequentialwallet(seedString)
	/*
    InitialHashBytes:=[]byte(seedString)
    for i:=0;i<wallet.NB_INITIAL_HASHES;i++{
        InitialHashBytes = utility.ComputeHashBytes(InitialHashBytes)
        wltgenprogress:=int(i*100/wallet.NB_INITIAL_HASHES)
        if (i)%(wallet.NB_INITIAL_HASHES/10)==0{
            fmt.Printf("Wallet Generation Progress %d %%\n",wltgenprogress)
			//creationProgressDialog.SetValue(float64(wltgenprogress))
			daemon.Walletstate=fmt.Sprintf("Wallet Generation Progress %d %%",wltgenprogress)
        }
        
    }
    
    fmt.Printf("%x\n",InitialHashBytes)

	pk := utility.PrivKeyFromBytes(InitialHashBytes)
    wlt.Privatekeyarray=append(wlt.Privatekeyarray,&pk)
	*/
	daemon.Wlt=wlt
	   
}

///////////////////////////////////////////////////////////////////
func passwordDecryptionFailedDialog(win fyne.Window,err string){
	selectWalletFileCallback := func(response bool){
		fmt.Println("Responded with", response)
		if response {
			//selectWalletFileDialog(win)
			newWalletCreationDialog(win)
		} else {
			//s
			//passwordCreationDialog(win)//
			nowalletFoundDialog(win,err)
		}
		
	}
	cnf := dialog.NewConfirm(err, "Would you like to re enter the wallet password", selectWalletFileCallback, win)
	cnf.SetDismissText("No  ")
	cnf.SetConfirmText("Yes ")
	cnf.Show()
}

func nowalletFoundDialog(win fyne.Window,err string){
	selectWalletFileCallback := func(response bool){
		fmt.Println("Responded with", response)
		if response {
			selectWalletFileDialog(win)
		} else {
			//
			newWalletCreationDialog(win)

		}
		
	}
	cnf := dialog.NewConfirm(err, "Would you like to select an existing wallet file", selectWalletFileCallback, win)
	cnf.SetDismissText("No  ")
	cnf.SetConfirmText("Yes ")
	cnf.Show()
}
func selectWalletFileDialog(win fyne.Window) {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			//dialog.ShowError(err, win)
			nowalletFoundDialog(win,"")
			return
		}
		if reader == nil {
			fmt.Println("Cancelled")
			nowalletFoundDialog(win,"")
			return
		}

		//if reader == nil {
		//	fmt.Println("Cancelled")
		//	return
		//}
		filepath:= reader.URI()
		reader.Close()
		//
		fmt.Println("Wallet file path",filepath)
		walletsettingsDisplayedMainwalletFilePath.Set(filepath.Path())
		daemon.MainwalletFilePath=filepath.Path()
		daemon.Usersettings.MainwalletFilePath=filepath.Path()
		if !daemon.Walletinstantiated(){
			daemon.MainwalletFilePath=filepath.Path()
			passwordDialog(win)
		} 

	}, win)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".dat"}))
	fd.Show()
}
//////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////

func registrationLoginDialog(win fyne.Window){

	//TODO support other types of wallets
	//combo := widget.NewSelect([]string{"Sequential Wallet", "Option 2"}, func(value string) {
	registerednames:=daemon.Wlt.GetRegisteredNames()
	combo := widget.NewSelect(registerednames, func(value string) {
		fmt.Println("Select set to", value)

		//daemon.Usersettings.Activeloginname=value
		daemon.PutActiveloginname(value)
	})
	//registerednamesarray:=daemon.Wlt.GetRegisteredNames()
	combo.SetSelected("")
	content := container.NewVBox(widget.NewLabel("LOGIN AS:"),combo)
	registrationLoginCallback := func(response bool){
		fmt.Println("Responded with", response)
		if response {
			
			//newSequentialWalletCreationDialog(win)
		}// else {
			//
			
			//nowalletFoundDialog(win,"")

		//}
		
	}
	cnf:=dialog.NewCustomConfirm("Would you like to login using the selected registration ?", "Yes ", "No  ", content,registrationLoginCallback, win)
	cnf.Show()
}
