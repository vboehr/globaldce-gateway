package gui


import (
	//"log"
	//"strings"
	"fmt"
	"time"
	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	//"net/url"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/data/binding"
	"github.com/globaldce/globaldce/daemon"
	"github.com/globaldce/globaldce/cli"
)

//type RegistredNameInfo struct {
//    name string
//}


func registrationScreen(win fyne.Window) fyne.CanvasObject {
	/*
	tabs := container.NewAppTabs(
		container.NewTabItem("Send to",  txbuilderScreen()),
		//container.NewTabItem("List of Contacts",  contactslistScreen()),
		//container.NewTabItem("Add Contact",  addContactScreen()),

	)
	tabs.SetTabLocation(container.TabLocationTop)
	*/
	//text :=widget.NewLabel("Hello")
	//var registrednameslist * widget.List
	//var registerednames [] string//RegistredNameInfo

	registerednames := binding.BindStringList(
		&[]string{},
	)
	
	fmt.Printf("%v",registerednames)
	registrednameslist := widget.NewListWithData(registerednames,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
		go func() {
			for {
				//fmt.Println("*******!!!!!!!!",registerednames)
				
				registerednames.Set(daemon.Wlt.GetRegisteredNames())
				time.Sleep(time.Second * 2)
				//str.Set(fmt.Sprintf("WALLET BALANCE is %d", daemon.Wlt.ComputeBalance()))
				
			}
		}()
	

	nameregistrationbutton:= widget.NewButton("NEW NAME REGISTRATION", func() {
        fmt.Println("creating a new name :")
		requestNameRegistrationDialog(win)
    })
	nameregistrationbuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(350, 40)),nameregistrationbutton)
	//layout:=container.New(layout.NewPaddedLayout(),container.NewVBox(registrednameslist,nameregistrationcontainer))
	registrednameslistcontainer:=container.New(layout.NewGridWrapLayout(fyne.NewSize(appscreenWidth, appscreenHeight*3/4)),registrednameslist)
	layout:=container.NewVBox(registrednameslistcontainer,nameregistrationbuttoncontainer)
	return layout

}
func  requestNameRegistrationDialog(win fyne.Window){
	requestedname := widget.NewEntry()
	depositamount := widget.NewEntry()
	//contactname.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")
	depositamount.Validator = validation.NewRegexp(`^[0-9]+$`, "deposited amount can only contain numbers")
	items := []*widget.FormItem{
		widget.NewFormItem("Requested Name", requestedname),
		widget.NewFormItem("Deposit Amount", depositamount),
		//widget.NewFormItem("Password", password),
		//widget.NewFormItem("Remember me", widget.NewCheck("", func(checked bool) {
		//	remember = checked
		//})),
	}

	dialog.ShowForm("Inorder to proceed with the name registration, please provide the following:    ", "Okay  ", "Cancel", items, func(b bool) {
		if !b {
			fmt.Println("canceled")
			//nowalletFoundDialog(win,"")
			return
		}
		if b {
			
			fmt.Println("text",)
			err:=cli.Sendnameregistration(daemon.Wireswarm,daemon.Mn,daemon.Wlt,requestedname.Text,depositamount.Text)
			if err!=nil{
				dialog.ShowError(err,win)
			} else {
				dialog.ShowInformation("Name Registration", "Name registration is being broadcasted", win)

			}
		}
	}, win)

}

