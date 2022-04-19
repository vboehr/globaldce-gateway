package gui

import (
	//"log"
	"fmt"
	"strconv"
	//"image/color"
	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	//"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	//"net/url"
	//"net/url"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2"
	//"github.com/globaldce/go-globaldce/cli"
	"github.com/globaldce/go-globaldce/daemon"
)

var walletsettingsDisplayedMainwalletFilePath binding.String

func settingsScreen(win fyne.Window) fyne.CanvasObject {
	// wallet path section
	walletsettingsDisplayedMainwalletFilePath=binding.NewString()
	walletsettingsDisplayedMainwalletFilePath.Set(daemon.Usersettings.MainwalletFilePath)
	walletpathentry:=widget.NewEntryWithData(walletsettingsDisplayedMainwalletFilePath)
	walletpathentry.Validator=nil
	selectwalletpathButton:= widget.NewButton("Select wallet path", func() {
		selectWalletFileDialog(win)
	})
	selectwalletpathButtonContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(180, 40)),selectwalletpathButton)
	walletpathentryContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(appscreenWidth-150, 40)),walletpathentry)
	walletpathContainer:=container.NewHBox(walletpathentryContainer,selectwalletpathButtonContainer)
	//
	//Nameregistrationtxfee:NameregistrationtxfeeDefault,
	nameregistrationtxfeeEntry:=widget.NewEntry()
	nameregistrationtxfeeEntry.Text=fmt.Sprintf("%d",daemon.Usersettings.Nameregistrationtxfee)
	nameregistrationtxfeeEntry.Validator=validation.NewRegexp(`^[0-9]+$`, "only numbers")
	nameregistrationtxfeeSetDefaultButton:= widget.NewButton("Set default", func() {
		nameregistrationtxfeeEntry.SetText(fmt.Sprintf("%d",daemon.Usersettings.Nameregistrationtxfee))
	})
	nameregistrationtxfeeSetDefaultButtonContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(100, 40)),nameregistrationtxfeeSetDefaultButton)
	nameregistrationtxfeeEntryContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(200, 40)),nameregistrationtxfeeEntry)
	nameregistrationtxfeeLabelContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(200, 40)),widget.NewLabel("Name registration fee"))
	nameregistrationtxfeeContainer:=container.NewHBox(nameregistrationtxfeeLabelContainer,nameregistrationtxfeeEntryContainer,nameregistrationtxfeeSetDefaultButtonContainer)

	//Publicposttxfee:PublicposttxfeeDefault,
	/*
	publicposttxfeeEntry:=widget.NewEntry()
	publicposttxfeeEntry.Text=fmt.Sprintf("%d",daemon.Usersettings.Publicposttxfee)
	publicposttxfeeEntry.Validator=validation.NewRegexp(`^[0-9]+$`, "only numbers")
	publicposttxfeeSetDefaultButton:= widget.NewButton("Set default", func() {
		publicposttxfeeEntry.SetText(fmt.Sprintf("%d",daemon.Usersettings.Publicposttxfee))
	})
	publicposttxfeeSetDefaultButtonContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(100, 40)),publicposttxfeeSetDefaultButton)
	publicposttxfeeEntryContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(200, 40)),publicposttxfeeEntry)
	publicposttxfeeLabelContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(200, 40)),widget.NewLabel("Public post fee"))
	publicposttxfeeContainer:=container.NewHBox(publicposttxfeeLabelContainer,publicposttxfeeEntryContainer,publicposttxfeeSetDefaultButtonContainer)
	*/
	//Sendtoaddressarraytxfee:SendtoaddressarraytxfeeDefault,
	sendtoaddressarraytxfeeEntry:=widget.NewEntry()
	sendtoaddressarraytxfeeEntry.Text=fmt.Sprintf("%d",daemon.Usersettings.Sendtoaddressarraytxfee)
	sendtoaddressarraytxfeeEntry.Validator=validation.NewRegexp(`^[0-9]+$`, "only numbers")
	sendtoaddressarraytxfeeSetDefaultButton:= widget.NewButton("Set default", func() {
		sendtoaddressarraytxfeeEntry.SetText(fmt.Sprintf("%d",daemon.Usersettings.Sendtoaddressarraytxfee))
	})
	sendtoaddressarraytxfeeSetDefaultButtonContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(100, 40)),sendtoaddressarraytxfeeSetDefaultButton)
	sendtoaddressarraytxfeeEntryContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(200, 40)),sendtoaddressarraytxfeeEntry)
	sendtoaddressarraytxfeeLabelContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(200, 40)),widget.NewLabel("Send to addresses fee"))
	sendtoaddressarraytxfeeContainer:=container.NewHBox(sendtoaddressarraytxfeeLabelContainer,sendtoaddressarraytxfeeEntryContainer,sendtoaddressarraytxfeeSetDefaultButtonContainer)
	//
	MiningrequestedCheck:=widget.NewCheck("Mining", func(checkstatus bool) { daemon.Usersettings.Miningrequested=checkstatus })
	MiningrequestedCheck.SetChecked(daemon.Usersettings.Miningrequested)
	// save settings button section
	saveSettingsButton:= widget.NewButton("Save settings", func() {
		fmt.Println("Saving settings")
		nameregistrationtxfeeNewValue, nameregistrationtxfeeNewValueErr := strconv.ParseInt(nameregistrationtxfeeEntry.Text, 10, 64)
		if nameregistrationtxfeeNewValueErr!=nil{
			dialog.ShowInformation("Error", "Entred value for name registration fee is inappropriete ", win)
			return
		}
		daemon.Usersettings.Nameregistrationtxfee=int(nameregistrationtxfeeNewValue)
		/*
		publicposttxfeeNewValue, publicposttxfeeNewValueErr := strconv.ParseInt(publicposttxfeeEntry.Text, 10, 64)
		if publicposttxfeeNewValueErr!=nil{
			dialog.ShowInformation("Error", "Entred value for public post fee is inappropriete ", win)
			return
		}
		daemon.Usersettings.Publicposttxfee=int(publicposttxfeeNewValue)
		*/

		sendtoaddressarraytxfeeNewValue, sendtoaddressarraytxfeeNewValueErr := strconv.ParseInt(sendtoaddressarraytxfeeEntry.Text, 10, 64)
		if sendtoaddressarraytxfeeNewValueErr!=nil{
			dialog.ShowInformation("Error", "Entred value for send to addresses fee is inappropriete ", win)
			return
		}
		daemon.Usersettings.Sendtoaddressarraytxfee=int(sendtoaddressarraytxfeeNewValue)
		_=daemon.SaveUsersettingsFile()
		dialog.ShowInformation("Saved settings", "Saved settings will take effect once the app starts", win)
	})
	saveSettingsButtonContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(150, 40)),saveSettingsButton)
	//
	return container.NewVBox(
		walletpathContainer,
		nameregistrationtxfeeContainer,
		//publicposttxfeeContainer,
		sendtoaddressarraytxfeeContainer,
		MiningrequestedCheck,
		saveSettingsButtonContainer,
	)
}