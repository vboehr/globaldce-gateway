package gui

import (
	"fyne.io/fyne/v2"
	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/layout"
    //"log"
    "fmt"
    "github.com/globaldce/globaldce-toolbox/cli"
    "github.com/globaldce/globaldce-toolbox/daemon"
    "fyne.io/fyne/v2/data/validation"
)

 	

type TxInfo struct {
    address string
    amount  string
}


var componentsList []TxInfo//{"test1: test1"}
//var componentsList = []string{"test1: test1"}
var selecteditemid=-1
/*
func main() {
	
    app := app.New()
    window := app.NewWindow("cool app")

 
    window.SetContent(txbuilderScreen())

    window.Resize(fyne.NewSize(800, 600))
    window.ShowAndRun()
}*/

func txbuilderScreen() fyne.CanvasObject{

   componentsTree := widget.NewList(
        func() int {
            return len(componentsList)
        },
        func() fyne.CanvasObject {


            return widget.NewLabel("template")
        },
        func(i widget.ListItemID, o fyne.CanvasObject) {



            o.(*widget.Label).SetText("Pay to "+componentsList[i].address +" an amount "+componentsList[i].amount) // i need to update this when componentsList was updated
        })
	componentsTree.OnSelected = func(id widget.ListItemID) {
		selecteditemid=id
	}

    //nameEntry := widget.NewEntry()
    //typeEntry := widget.NewEntry()
/*
!!!!!!!!!!!
searchButton:=widget.NewButton("Fixed size window", func() {
			w := fyne.CurrentApp().NewWindow("Fixed")
			w.SetContent(fyne.NewContainerWithLayout(layout.NewCenterLayout(), widget.NewLabel("Hello World!")))

			w.Resize(fyne.NewSize(240, 180))
			w.SetFixedSize(true)
			w.Show()
		})
!!!!
*/
    /*
    form := &widget.Form{
        Items: []*widget.FormItem{
            {Text: "Name", Widget: nameEntry},
            {Text: "Type", Widget: typeEntry}},
        OnSubmit: func() {
            componentsList = append(componentsList, TxInfo{nameEntry.Text, typeEntry.Text}) // append an item to componentsList array
	componentsTree.Refresh()
        },
    }
    */

    addressinput := widget.NewEntry()// TODO validation of address
	addressinput.SetPlaceHolder("Enter recipient address ...")

	amountinput := widget.NewEntry()// TODO validation of amount as float
	amountinput.SetPlaceHolder("Enter amount ...")

	outputform := container.NewVBox(addressinput,amountinput, widget.NewButton("Add", func() {
		fmt.Printf("Address was: %s - Amount was: %s", addressinput.Text,amountinput.Text)
        componentsList = append(componentsList, TxInfo{addressinput.Text, amountinput.Text}) // append an item to componentsList array
        componentsTree.Refresh()
		addressinput.Text=""
		addressinput.Refresh()
		amountinput.Text=""
		amountinput.Refresh()
	}))




    rmvbutton:= widget.NewButton("Remove Selection", func() {
		if selecteditemid!=-1 && 0<len(componentsList) {
            
			componentsList=remove(componentsList,selecteditemid)
            
			componentsTree.Refresh()
		}
		
	})
	
    //layout := container.New(layout.NewGridWrapLayout(fyne.NewSize(350, 500)),label , form)


    label:= container.NewBorder(rmvbutton, nil, nil,nil,componentsTree)
    
    completiontext:=widget.NewLabel("  ")// TODO Add balance information
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
    //
    completebutton:= widget.NewButton("SEND", func() {
        fmt.Println("got :",componentsList)
        var addressarray []string
        var amountarray []string
        for i,_:=range componentsList{
            addressarray=append(addressarray,componentsList[i].address)
            amountarray=append(amountarray,componentsList[i].amount)
        }
        cli.Sendtoaddressarray(daemon.Wireswarm, daemon.Wlt,addressarray,amountarray,sendtoaddressarraytxfeeEntry.Text)
        var emptycomponentsList []TxInfo
        componentsList=emptycomponentsList
            
        componentsTree.Refresh()

    })


    completebuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(appscreenWidth/4, 40)),completebutton)
    //
    formlayout:=container.New(layout.NewPaddedLayout(),container.NewVBox(outputform,completiontext,sendtoaddressarraytxfeeContainer,completebuttoncontainer))
    //layout:= container.NewHSplit(label,formlayout)

    return container.NewHSplit(label,formlayout)//layout
}

func remove(s []TxInfo, i int) []TxInfo {
    //s[len(s)-1], s[i] = s[i], s[len(s)-1]
    //return s[:len(s)-1]
    return append(s[:i], s[i+1:]...)
}
