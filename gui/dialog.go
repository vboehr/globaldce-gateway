package gui


import (
	//"log"
	"fmt"
	"errors"
	//"fyne.io/fyne/v2/app"
	//"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	//"fyne.io/fyne/v2/data/validation"
	//"net/url"
	//"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2"
	"github.com/globaldce/globaldce-toolbox/daemon"
	"github.com/globaldce/globaldce-toolbox/utility"
	"github.com/globaldce/globaldce-toolbox/cli"
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
				lerr:=daemon.Wlt.LoadJSONFile(daemon.MainwalletFilePath,daemon.MainwalletFileKey)
				if lerr!=nil{
					//nowalletFoundDialog(win,"Could not decrypt walletfile "+daemon.MainwalletFilePath)
					passwordDecryptionFailedDialog(win,"Could not decrypt walletfile "+daemon.MainwalletFilePath)
				} else {
					daemon.Walletloaded=true
				}
		}
		//var rememberText string
		//if remember {
		//	rememberText = "and remember this login"
		//}

		//fmt.Println("Entred password", password.Text)
	}, win)

}
func  passwordCreationDialog(win fyne.Window) {
	//username := widget.NewEntry()
	//username.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")
	cli.GenerateNewMainwalletFilePath()
	firstpassword := widget.NewPasswordEntry()
	//password.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "password can only contain letters, numbers, '_', and '-'")
	firstpassword.Validator=passwordValidation()
	secondpassword := widget.NewPasswordEntry()
	secondpassword.Validator=passwordValidation()
	//remember := false
	items := []*widget.FormItem{
		//widget.NewFormItem("Username", username),
		widget.NewFormItem("Password", firstpassword),
		widget.NewFormItem("Password", secondpassword),
		//widget.NewFormItem("Remember me", widget.NewCheck("", func(checked bool) {
		//	remember = checked
		//})),
	}
	dialog.ShowForm("Password required to create a new wallet          ", "Okay  ", "Cancel", items, func(b bool) {
		if !b {
			fmt.Println("canceled")
			nowalletFoundDialog(win,"")
			return
		}
		//fmt.Println("Password",password.Text)

		if firstpassword.Text==secondpassword.Text && b {
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
				daemon.Walletloaded=true	
		} else {
			nowalletFoundDialog(win,"Wallet creation abroted entred password do not match")
		}
		//var rememberText string
		//if remember {
		//	rememberText = "and remember this login"
		//}
		//fmt.Println("Entred password", password.Text)
	}, win)


}

func passwordDecryptionFailedDialog(win fyne.Window,err string){
	selectWalletFileCallback := func(response bool){
		fmt.Println("Responded with", response)
		if response {
			//selectWalletFileDialog(win)
			passwordDialog(win)
		} else {
			//s
			//passwordCreationDialog(win)
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
			passwordCreationDialog(win)

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
		if !daemon.Walletloaded{
			daemon.MainwalletFilePath=filepath.Path()
			passwordDialog(win)
		} 

	}, win)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".dat"}))
	fd.Show()
}
