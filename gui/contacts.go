package gui


import (
	//"log"
	"strings"
	"fmt"
	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	//"net/url"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2"
)
/*
	names := make([]string, 5)
	
	names[0] = "Alice"
	names[1] = "Bob"
	names[2]="Cara"
	names[3] = "Linda"
	names[4]="Loren"
*/

type ContactInfo struct {
    name string
    address  string
}
var list * widget.List
var data []ContactInfo

var contacts []ContactInfo

/*
func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("cool app")
myWindow.Resize(fyne.Size{800, 600})


	myWindow.SetContent(contactsScreen())
	myWindow.ShowAndRun()
}
*/



func contactsScreen() fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("List of Contacts",  contactslistScreen()),
		container.NewTabItem("Add Contact",  addContactScreen()),

	)

//tabsisTransitioning
	tabs.OnChanged = func(_ *container.TabItem) {
		//list.Refresh()
		fmt.Println("!!!!!!",list)
	}
	tabs.SetTabLocation(container.TabLocationTop)



	return tabs

}

func addContactScreen() fyne.CanvasObject {
	inputname := widget.NewEntry()
	inputname.SetPlaceHolder("Enter name...")
	inputaddress := widget.NewEntry()
	inputaddress.SetPlaceHolder("Enter address...")

	content := container.New(  layout.NewGridWrapLayout(fyne.NewSize(350, 500)),container.NewVBox(inputname, widget.NewButton("Save", func() {
		//fmt.Println("Content was:", input.Text)
		addcontact(inputname.Text,inputaddress.Text)
		data=contacts
		//list.Refresh()
		//fmt.Println("contacts",contacts)
		
	})))
	scr := container.New(layout.NewCenterLayout(), content)
	return scr

}

func contactslistScreen() fyne.CanvasObject {
	fmt.Println("**********contacts",contacts)
	data=contacts

	//icon := widget.NewIcon(nil)
	input := widget.NewEntry()
	input.SetPlaceHolder("Search Contact ...")

	
	saveButton:= widget.NewButton("Copy", func() {
		fmt.Println("Content was:", input.Text)
	})

	label := container.NewVBox(
		widget.NewLabel("The top row of the VBox                                           "),
			input,)
	label.Add(saveButton)
	hbox := container.NewVBox( widget.NewIcon(nil), label)

	list = widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id].name)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		//label.SetText(data[id])
		//icon.SetResource(theme.DocumentIcon())
	//input := widget.NewEntry()
	input.Text=data[id].name

	input.Refresh()
	

	//label = container.NewVBox(input, widget.NewButton("Save", func() {
	//	fmt.Println("Content was:", input.Text)
	//}))


	}
	//list.OnUnselected = func(id widget.ListItemID) {
	//	label.SetText("Select An Item From The List")
	//	icon.SetResource(nil)
	//}
	//list.Select(1)
	//return list

	input.OnChanged = func(input string) {
	//optionSelectEntry.SetOptions(filteredListForSelect(&input))
		fmt.Println("Content was:", input)
		data=filter(contacts,input)
		fmt.Println("*****",list)
	
		list.Refresh()
		
	}
	//listpointer=list
	return container.NewHSplit( container.NewCenter(hbox),list)
}

func filter(s []ContactInfo,input string) []ContactInfo {
	input =strings.ToUpper(input)
	names := make([]ContactInfo, 0)
	for _, selement := range s {
		if strings.Index(strings.ToUpper(selement.name), input)  ==0{
			names=append(names,selement)
		}
    }
    return names
}

func addcontact(name string,address string) {

	//nbcontacts:=len(contacts)
	addedcontact:=ContactInfo{name,address}

	for i, selement := range contacts {
		if  0 < strings.Compare(selement.name,name)  {
			newcontacts := make([]ContactInfo, 0)
			newcontacts= append(newcontacts,contacts[:i]...)
			newcontacts=append(newcontacts,addedcontact)

			//fmt.Println("****",i,newcontacts,contacts[i:])
			newcontacts=append(newcontacts,contacts[i:]...)
			contacts=newcontacts
			//fmt.Println("contacts",name,address,contacts)
			return

		}
    }

	//if len(contacts)==0{
		contacts=append(contacts,addedcontact)
		//fmt.Println("contacts",name,address,contacts)
		return	
	//}
    
}
