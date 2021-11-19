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
	"fyne.io/fyne/v2/data/binding"
	"github.com/globaldce/globaldce-toolbox/daemon"
	//"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"github.com/globaldce/globaldce-toolbox/cli"
)

var publicsharename string

func shareScreen(win fyne.Window) fyne.CanvasObject {
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
			return widget.NewLabel("Registred Names")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
		registrednameslist.OnSelected = func(id widget.ListItemID) {
			//label.SetText(data[id])
			//icon.SetResource(theme.DocumentIcon())
		//input := widget.NewEntry()
		textvalue,_:=registerednames.GetValue(id)
		publicsharename=textvalue

		}
		go func() {
			for {
				//fmt.Println("*******!!!!!!!!",registerednames)
				
				registerednames.Set(daemon.Wlt.GetRegisteredNames())
				time.Sleep(time.Second * 2)
				//str.Set(fmt.Sprintf("WALLET BALANCE is %d", daemon.Wlt.ComputeBalance()))
				
			}
		}()
	
	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("Enter public post text...")
	

	nameregistrationbutton:= widget.NewButton("SHARE PUBLIC POST", func() {
		if publicsharename!=""{
			publicsharetext:=input.Text
			fmt.Printf("creating a new public post for %s : %s\n",publicsharename,publicsharetext)
			cli.Sendpublicpost(daemon.Wireswarm,daemon.Mn,daemon.Wlt,publicsharename,publicsharetext)
			dialog.ShowInformation("Public Post", "Public post is being broadcasted", win)
			input.SetText("")
		}
        


		
    })
	nameregistrationbuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(350, 40)),nameregistrationbutton)
	//nameregistrationbuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(350, 40)),nameregistrationbutton)

	//layout:=container.New(layout.NewPaddedLayout(),container.NewVBox(registrednameslist,nameregistrationcontainer))
	///////////registrednameslistcontainer:=container.New(layout.NewGridWrapLayout(fyne.NewSize(appscreenWidth, appscreenHeight/4)),registrednameslist)
	sharelayout:=container.NewVBox(input,nameregistrationbuttoncontainer)
	return container.NewHBox(registrednameslist, sharelayout)//layout
	//container.New(layout.NewFormLayout(), label1, value1, label2, value2)
}
/*
func  requestNameRegistrationDialog(win fyne.Window){
	requestedname := widget.NewEntry()
	depositamount := widget.NewEntry()
	//contactname.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")
	depositamount.Validator = validation.NewRegexp(`^[0-9]+$`, "username can only contain numbers")
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
			
			fmt.Println("text",requestedname.Text,depositamount.Text)
		}
	}, win)

}

*/