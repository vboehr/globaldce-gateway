package gui


import (
	//"log"
	//"strings"
	//"fmt"
	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/theme"
	//"fyne.io/fyne/v2/widget"
	//"net/url"
	//"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2"
)



func sendScreen() fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Send to",  txbuilderScreen()),
		//container.NewTabItem("List of Contacts",  contactslistScreen()),
		//container.NewTabItem("Add Contact",  addContactScreen()),

	)
	/*
	//tabsisTransitioning
	tabs.OnChanged = func(_ *container.TabItem) {
		//list.Refresh()
		fmt.Println("!!!!!!",list)
	}
	*/
	tabs.SetTabLocation(container.TabLocationTop)



	return tabs

}

