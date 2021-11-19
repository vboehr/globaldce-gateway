package gui


import (
	//"log"
	//"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/theme"
	//"fyne.io/fyne/v2/dialog"
	//"fyne.io/fyne/v2/layout"
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
	"github.com/globaldce/globaldce-toolbox/cli"
)

const appscreenWidth =800
const appscreenHeight = 600


func Start(appname string) {
	applog.Init()
	settingserr:=cli.LoadUsersettingsFile()
	if settingserr!=nil{
		//
		cli.SetDefaultSettings()
	}
	cli.ApplyUsersettings()
	applog.Notice("Mainwalletpath %s",cli.Usersettings.MainwalletFilePath)
	daemon.Miningrequested=true
	daemon.Seed=true
	


	/////////////////////////////////////////
	myApp := app.New()
	//myApp := app.NewWithID("***.**") //TOBETESTED
	//myApp.SetIcon(resourceDPng) //TOBETESTED can be used to set icon
	myWindow := myApp.NewWindow(appname)
	myApp.Settings().SetTheme(&myTheme{})
	myWindow.Resize(fyne.NewSize(appscreenWidth, appscreenHeight))
	//myWindow.SetFixedSize(true)//TODO FOR MOBILE
	myWindow.SetIcon(resourceLogoPng)


	//widget.ShowPopUp(widget.NewLabel("Home tab"), myWindow.Canvas())
	//passwordDialog(myWindow)
	/////////////////////////////////////
	wlt:=new(wallet.Wallet)
	if _, err := os.Stat( daemon.MainwalletFilePath); !os.IsNotExist(err) {
		//for {
			passwordDialog(myWindow)

		//}
	} else {
		
		nowalletFoundDialog(myWindow,"walletfile "+daemon.MainwalletFilePath+" does not exist.")
	}
	daemon.Wlt=wlt
	go daemon.Mainloop()
	/////////////////////////////////////
	
	tabs := container.NewAppTabs(
		//container.NewTabItemWithIcon("Wallet",theme.FolderIcon(),  overviewScreen()),
		container.NewTabItem("Home",  homeScreen(myWindow)),
		container.NewTabItem("Balance",  balanceScreen()),
		container.NewTabItem("Registration",  registrationScreen(myWindow)),
		//container.NewTabItem("Send",  sendScreen()),//txbuilderScreen()
		container.NewTabItem("Send to",  txbuilderScreen()),
		container.NewTabItem("Settings",  settingsScreen()),	
	)

	//tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")))

	tabs.SetTabLocation(container.TabLocationTop)
	
	myWindow.SetContent(tabs)
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
		daemon.Wlt.SaveJSONFile(daemon.MainwalletFilePath,daemon.MainwalletFileKey)
		_=cli.SaveUsersettingsFile()
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
