package gui


import (
	//"log"
	//"strings"
	//"os"
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
	"github.com/globaldce/globaldce-gateway/daemon"
	"github.com/globaldce/globaldce-gateway/cli"
)

//type RegistredNameInfo struct {
//    name string
//}

//var selectednameregistration string
var selectednameregistrationid int

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
		//

		registrednameslist.OnSelected = func(id widget.ListItemID) {
			textvalue,_:=registerednames.GetValue(id)
			//selectednameregistration=textvalue
			_=textvalue

		}
		//
		go func() {
			for {
				//fmt.Println("*******!!!!!!!!",registerednames)
				registerednamesdescriptionarray:=daemon.Wlt.GetRegisteredNames()
				for i,registeredname:= range registerednamesdescriptionarray{
					registerednamepath,_:=daemon.GetCachedDirPathForRegistredName(registeredname)
					if registerednamepath!=""{
						registerednamesdescriptionarray[i]+=registeredname+" CONTENT DIRECTORY: "+registerednamepath
					}
					
				}
				registerednames.Set(registerednamesdescriptionarray)
				time.Sleep(time.Second * 2)
				//str.Set(fmt.Sprintf("WALLET BALANCE is %d", daemon.Wlt.ComputeBalance()))
				
			}
		}()
	

	nameregistrationbutton:= widget.NewButton("NEW NAME REGISTRATION", func() {
        fmt.Println("creating a new name :")
		requestNameRegistrationDialog(win)
    })
	nameregistrationbuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(350, 40)),nameregistrationbutton)
	//
	//
	nameunregistrationbutton:= widget.NewButton("NAME UNREGISTRATION", func() {
		
        fmt.Println("name unregistration")
		//requestNameUnregistrationDialog(win)
		//
		registerednamesarray:=daemon.Wlt.GetRegisteredNames()
		err:=cli.Sendnameunregistration(daemon.Wireswarm,daemon.Mn,daemon.Wlt,registerednamesarray[selectednameregistrationid])
		if err!=nil{
			dialog.ShowError(err,win)
		} else {
			dialog.ShowInformation("Name Unregistration", "Name unregistration is being broadcasted", win)

		}

    })
	setcontentdirbutton:= widget.NewButton("SET CONTENT DIRECTORY", func() {
        fmt.Println("set content directory")
		selectContentDirDialog(win)
    })
	commitcontentdirbutton:= widget.NewButton("COMMIT CONTENT DIRECTORY", func() {
        fmt.Println("commit content directory")
		commitContentDirDialog(win)
    })
	nameunregistrationbuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(350, 40)),nameunregistrationbutton)
	setcontentdirbuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(350, 40)),setcontentdirbutton)
	commitcontentdirbuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(350, 40)),commitcontentdirbutton)
	//layout:=container.New(layout.NewPaddedLayout(),container.NewVBox(registrednameslist,nameregistrationcontainer))
	registrednameslistcontainer:=container.New(layout.NewGridWrapLayout(fyne.NewSize(appscreenWidth, appscreenHeight*3/4)),registrednameslist)
	layout:=container.NewVBox(nameregistrationbuttoncontainer,registrednameslistcontainer,nameunregistrationbuttoncontainer,setcontentdirbuttoncontainer,commitcontentdirbuttoncontainer)
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

////////////////////////////////
func selectContentDirDialog(win fyne.Window) {
	folderd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {

		folderpath:= uri.Path()
		
		fmt.Println("content folder path",folderpath)
		registerednamesdescriptionarray:=daemon.Wlt.GetRegisteredNames()
		daemon.PutCachedDirPathForRegistredName(registerednamesdescriptionarray[selectednameregistrationid],folderpath)


	}, win)
	folderd.Show()
}
////////////////////////////////
func commitContentDirDialog(win fyne.Window) {
//Are you sure
	//var registrednamecommittxfeevalue string

	//TODO support other types of wallets
	//combo := widget.NewSelect([]string{"Sequential Wallet", "Option 2"}, func(value string) {
	//combo := widget.NewSelect([]string{"Sequential Wallet"}, func(value string) {
	//	fmt.Println("Select set to", value)
	//})
	//combo.SetSelected("Sequential Wallet")
	//content := container.NewVBox(widget.NewLabel("Wallet Type :"),combo)
	registrednamecommittxfeeEntry:=widget.NewEntry()
	registrednamecommittxfeeEntry.Text=fmt.Sprintf("%d",daemon.Usersettings.Nameregistrationtxfee)
	registrednamecommittxfeeEntry.Validator=validation.NewRegexp(`^[0-9]+$`, "only numbers")
	//registrednamecommittxfeeSetDefaultButton:= widget.NewButton("Set default", func() {
	//	registrednamecommittxfeeEntry.SetText(fmt.Sprintf("%d",daemon.Usersettings.Nameregistrationtxfee))
	//})
	//nameregistrationtxfeeSetDefaultButtonContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(100, 40)),nameregistrationtxfeeSetDefaultButton)
	registrednamecommittxfeeEntryContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(200, 40)),registrednamecommittxfeeEntry)
	registrednamecommittxfeeLabelContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(250, 40)),widget.NewLabel("Registred Name Commit Fee"))
	registrednamecommittxfeeContainer:=container.NewHBox(registrednamecommittxfeeLabelContainer,registrednamecommittxfeeEntryContainer)
	/////////
	selectWalletFileCallback := func(response bool){
		fmt.Println("Responded with", response)
		if response {
			registerednamesdescriptionarray:=daemon.Wlt.GetRegisteredNames()
			selectednameregistration:=registerednamesdescriptionarray[selectednameregistrationid]
			path,_:=daemon.GetCachedDirPathForRegistredName(selectednameregistration)
			cerr:=cli.Sendregistrednamecontentcommit(daemon.Wireswarm,daemon.Mn,daemon.Wlt,selectednameregistration,path,registrednamecommittxfeeEntry.Text)
			//Sendregistrednamecontentcommit(ws *wire.Swarm,mn *mainchain.Maincore,wlt *wallet.Wallet,namestring string,contentfolderstring string,amountfeestring string){
			if cerr!=nil{
				dialog.ShowError(cerr,win)
			} else {
				dialog.ShowInformation("Registred Name Content", "Registred name content is being broadcasted", win)
			}		
			
		}// else {
		//}
		
	}
	/////////

	cnf:=dialog.NewCustomConfirm("Would you like to send a commit transaction fee ?", "Yes ", "No  ", registrednamecommittxfeeContainer,selectWalletFileCallback, win)
	cnf.Show()
}