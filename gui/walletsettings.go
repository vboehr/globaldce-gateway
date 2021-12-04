package gui

import (
	//"log"
	"fmt"
	//"image/color"
	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/theme"
	//"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	//"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	//"net/url"
	//"net/url"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2"
	"github.com/globaldce/globaldce-toolbox/cli"
)

var walletsettingsDisplayedMainwalletFilePath binding.String

func walletSettingsScreen(win fyne.Window) fyne.CanvasObject {
/*
	entrymainwalletpath := widget.NewEntry()
	entrymainwalletpath.Text=cli.Usersettings.MainwalletFilePath
	//entry1.Validator = validation.NewRegexp(`^[0-9]+$`, "deposited amount can only contain numbers")
	//entry2 := widget.NewEntry()
	//entry2.Validator = validation.NewRegexp(`^[0-9]+$`, "deposited amount can only contain numbers")
	
	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "main wallet file path", Widget:entrymainwalletpath},
			//{Text: "Cool2", Widget: entry2},
		},
		OnSubmit: func() { // optional, handle form submission
			log.Println("Main wallet file path set to:", entrymainwalletpath.Text)
			//log.Println("entry 2:", entry2.Text)
			/////////////////////////////////////
			cli.Usersettings.MainwalletFilePath=entrymainwalletpath.Text
			/////////////////////////////////////
		},
	}

	return form
*/
walletsettingsDisplayedMainwalletFilePath=binding.NewString()
walletsettingsDisplayedMainwalletFilePath.Set(cli.Usersettings.MainwalletFilePath)
walletpathentry:=widget.NewEntryWithData(walletsettingsDisplayedMainwalletFilePath)
walletpathentry.Validator=nil
//walletpath := widget.NewLabelWithData(str)
//boundString := binding.NewString()
//s, _ := boundString.Get()
//log.Printf("Bound = '%s'", s)

//walletpathentry.Text=cli.Usersettings.MainwalletFilePath

//walletpathentry.Text=str
selectwalletpathButton:= widget.NewButton("Select wallet path", func() {
	
	//fmt.Println("********")
	selectWalletFileDialog(win)
	
})
selectwalletpathButtonContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(180, 40)),selectwalletpathButton)
walletpathentryContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(appscreenWidth-150, 40)),walletpathentry)
walletpathContainer:=container.NewHBox(walletpathentryContainer,selectwalletpathButtonContainer)
saveSettingsButton:= widget.NewButton("Save settings", func() {
	
	fmt.Println("Saving settings")
	_=cli.SaveUsersettingsFile()
})
saveSettingsButtonContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(150, 40)),saveSettingsButton)
return container.NewVBox(walletpathContainer,saveSettingsButtonContainer)
}