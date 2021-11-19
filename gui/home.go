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



func homeScreen(win fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Explore",  exploreScreen(win)),
		//container.NewTabItem("Search",  contactslistScreen()),
		container.NewTabItem("Share",  shareScreen(win)),

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

