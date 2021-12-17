package gui


import (
	//"log"
	//"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/theme"
	//"fyne.io/fyne/v2/dialog"
	//"fyne.io/fyne/v2/layout"
	//"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2"
	"github.com/globaldce/globaldce-toolbox/applog"
    "os"
	//"errors"
    "github.com/globaldce/globaldce-toolbox/wallet"
    //"time"
    "fmt"
    //"strings"
    //"path/filepath"
    "github.com/globaldce/globaldce-toolbox/daemon"
	//"github.com/globaldce/globaldce-toolbox/cli"
)

const appscreenWidth =1000
const appscreenHeight = 600
//var winredraw bool 
//var searchtext string

var guiApp fyne.App

func Start(appname string) {
	applog.Init()
	settingserr:=daemon.LoadUsersettingsFile()
	if settingserr!=nil{
		//
		daemon.SetDefaultSettings()
	}
	daemon.ApplyUsersettings()
	applog.Notice("Mainwalletpath %s",daemon.MainwalletFilePath)
	//daemon.Miningrequested=true
	daemon.Seed=true
	


	/////////////////////////////////////////
	guiApp = app.New()
	//guiApp := app.NewWithID("***.**") //TOBETESTED
	//guiApp.SetIcon(resourceDPng) //TOBETESTED can be used to set icon
	myWindow := guiApp.NewWindow(appname)
	guiApp.Settings().SetTheme(&myTheme{})
	myWindow.Resize(fyne.NewSize(appscreenWidth, appscreenHeight))
	//myWindow.SetFixedSize(true)//TODO FOR MOBILE
	myWindow.SetIcon(resourceLogoPng)


	//widget.ShowPopUp(widget.NewLabel("Home tab"), myWindow.Canvas())
	//passwordDialog(myWindow)
	/////////////////////////////////////
	wlt:=new(wallet.Wallet)
	//wlt.SetupDone=false
	
	//for (!daemon.Walletloaded) {

		if _, err := os.Stat( daemon.MainwalletFilePath); !os.IsNotExist(err) {
		
				passwordDialog(myWindow)

			//}
		} else {
			
			nowalletFoundDialog(myWindow,"walletfile "+daemon.MainwalletFilePath+" does not exist.")
		}
	//}
	daemon.Wlt=wlt
	fmt.Printf("Starting main loop\n")
	go daemon.Mainloop()
	/////////////////////////////////////
	//winredraw=true
	////////////////////////////////////
	//go func() {
	//	for {
	//		if winredraw {
			//hometab=homeScreen(myWindow)
			tabs := container.NewAppTabs(
				//container.NewTabItemWithIcon("Wallet",theme.FolderIcon(),  overviewScreen()),
				container.NewTabItem("Home",  homeScreen(myWindow)),
				container.NewTabItem("Balance",  balanceScreen()),
				container.NewTabItem("Registration",  registrationScreen(myWindow)),
				//container.NewTabItem("Send",  sendScreen()),//txbuilderScreen()
				container.NewTabItem("Send to",  txbuilderScreen()),
				container.NewTabItem("Settings",  settingsScreen(myWindow)),	
			)
		
			//tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")))
		
			tabs.SetTabLocation(container.TabLocationTop)
			
			myWindow.SetContent(tabs)
		
			//winredraw=false
			//time.Sleep(time.Minute * 2)
			//}
	//	}
	//}()
	/*
	go func() {
		for {
			tabs.Refresh()
			time.Sleep(time.Second * 2)
		}
	}()
	*/

	/*
	myWindow.SetCloseIntercept(func() { //SetOnClosed(func())
		dialog.ShowConfirm("Quitting", "Do you want to quit ?",
			func(response bool) {
				if response {
					fmt.Printf("Closing")
					myWindow.Close()
				}
			}, myWindow)
	})
	*/

	
	//myWindow.SetOnClosed(func(){
	//	fmt.Println("Closing")
	//	daemon.Wlt.SaveJSONFile(daemon.MainwalletFilePath,daemon.MainwalletFileKey)
	//})
	
	myWindow.ShowAndRun()
		fmt.Println("Closing")
		if daemon.Walletloaded {
			daemon.Wlt.SaveJSONFile(daemon.MainwalletFilePath,daemon.MainwalletFileKey)
		}
		_=daemon.SaveUsersettingsFile()
}

/*
func entryScreen() fyne.CanvasObject {
input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")
	

	content := container.NewVBox(input, widget.NewButton("Save", func() {
		fmt.Println("Content was:", input.Text)
	}))


return content

}
*/
