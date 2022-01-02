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
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/globaldce/globaldce-toolbox/daemon"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"github.com/globaldce/globaldce-toolbox/cli"
)

var publicsharename string
var attachmentpathArray []string//
var selectedattachmentid=-1

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
	
	linkinput := widget.NewEntry()
	linkinput.SetPlaceHolder("Enter public post link...")

	textinput := widget.NewMultiLineEntry()
	textinput.SetPlaceHolder("Enter public post text...")
	

	//sharebuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(350, 40)),sharebutton)

	//layout:=container.New(layout.NewPaddedLayout(),container.NewVBox(registrednameslist,nameregistrationcontainer))
	///////////registrednameslistcontainer:=container.New(layout.NewGridWrapLayout(fyne.NewSize(appscreenWidth, appscreenHeight/4)),registrednameslist)
	//attachmentcontainer:=container.New(layout.NewGridWrapLayout(fyne.NewSize(appscreenWidth*2/3, appscreenHeight*2/3)),attachmentbuilderScreen(win))

	//sharelayout:=container.NewVBox(input,attachmentcontainer,sharebuttoncontainer)

	//
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
	
	//

	sharebutton:= widget.NewButton("SHARE PUBLIC POST", func() {
		if publicsharename!=""{
			publicsharetext:=textinput.Text
			publicsharelink:=linkinput.Text
			fmt.Printf("creating a new public post for %s : %s\n",publicsharename,publicsharelink,publicsharetext)
			cli.Sendpublicpost(daemon.Wireswarm,daemon.Mn,daemon.Wlt,publicsharename,publicsharelink,publicsharetext,attachmentpathArray,publicposttxfeeEntry.Text)
			dialog.ShowInformation("Public Post", "Public post is being broadcasted", win)
			textinput.SetText("")
			linkinput.SetText("")
			attachmentpathArray=make([]string,0)
			 
			componentsTree.Refresh()
		}
    })
	sharebuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(350, 40)),sharebutton)
	//
	basicshare:=container.NewVBox(linkinput,textinput,publicposttxfeeContainer,sharebuttoncontainer)
	attachmentshare:=attachmentbuilderScreen(win)
	sharelayout:=container.NewVSplit(attachmentshare,basicshare)
	return container.NewHSplit(registrednameslist, sharelayout)//layout
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
var componentsTree *widget.List

func attachmentbuilderScreen(win fyne.Window) fyne.CanvasObject{

	componentsTree = widget.NewList(
		 func() int {
			 return len(attachmentpathArray)
		 },
		 func() fyne.CanvasObject {
 
 
			 return widget.NewLabel("template")
		 },
		 func(i widget.ListItemID, o fyne.CanvasObject) {
 
 
 
			 o.(*widget.Label).SetText("Attachment path "+attachmentpathArray[i]) // i need to update this when attachmentpathArray was updated
		 })
	 componentsTree.OnSelected = func(id widget.ListItemID) {
		 selectedattachmentid=id
	 }
 
 
	 addbutton:= widget.NewButton("Add attachment file", func() {
 
		 //
		 selectImageFileDialog(win,componentsTree)
		 //componentsTree.Refresh()
	 })
  
 
	 rmvbutton:= widget.NewButton("Remove selection", func() {
		 if selectedattachmentid!=-1 && 0<len(attachmentpathArray) {
			 
			 attachmentpathArray=removeattachment(attachmentpathArray,selectedattachmentid)
			 
			 componentsTree.Refresh()
		 }
		 
	 })
	 
	/*
	 completebutton:= widget.NewButton("SEND", func() {
		 fmt.Println("got :",attachmentpathArray)
 
	 })
 */
 
	 //completebuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(float32(appscreenWidth/4), 40)),completebutton)
 
	 buttonscontainer:=container.NewHBox(addbutton,rmvbutton)
	 label:= container.NewBorder(buttonscontainer, nil, nil,nil,componentsTree)
	 
	 //completiontext:=widget.NewLabel("  ")// TODO Add balance information
	 //formlayout:=container.New(layout.NewPaddedLayout(),container.NewVBox(completiontext,completebuttoncontainer))
 
   
 
	 //return container.NewVSplit(label,formlayout)
	 return label
 }
 
 func removeattachment(s []string, i int) []string {
	 //s[len(s)-1], s[i] = s[i], s[len(s)-1]
	 //return s[:len(s)-1]
	 return append(s[:i], s[i+1:]...)
 }
 
 func selectImageFileDialog(win fyne.Window,componentsTree fyne.Widget) {
	 fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		 if err != nil {
			 //dialog.ShowError(err, win)
			 //noimageFoundDialog(win,"")
			 fmt.Printf("selectImageFile error %v",err)
			 return
		 }
		 if reader == nil {
			 fmt.Println("Cancelled")
			 //noimageFoundDialog(win,"")
			 return
		 }
 
		 fileuri:= reader.URI()
		 reader.Close()
		 
		 fmt.Println("Wallet file path",fileuri.Path())
		 attachmentpathArray = append(attachmentpathArray,fileuri.Path())
		 componentsTree.Refresh()
	 }, win)
	 fd.SetFilter(storage.NewExtensionFileFilter([]string{".jpg"}))
	 fd.Show()
 }