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
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"github.com/globaldce/globaldce-gateway/cli"
	"github.com/globaldce/globaldce-gateway/daemon"
)

//var selectednameregistration string
//var selectednameregistrationid int
//var selectedregistredname string

func registrationScreen(win fyne.Window) fyne.CanvasObject {
	/////////////////////////
	registrationloginstr := binding.NewString()

	////////////////////////
	//
	go func() {
		for {

			if daemon.GetActiveloginname() == "" {
				registrationloginstr.Set("Not logged in")
			} else {
				registrationloginstr.Set(daemon.GetActiveloginname())
			}
			time.Sleep(time.Second * 2)

		}
	}()
	////////////////////////
	registrationloginlabel := widget.NewLabelWithData(registrationloginstr)

	//registrationlogininputContainer :=container.New(  layout.NewGridWrapLayout(fyne.NewSize(150, 40)),registrationloginlabel)
	/*
		registrationloginentryContainer := container.NewHBox(registrationlogininputContainer,
			widget.NewButton("LOGIN AS", func() {
				fmt.Println("Entred registrationlogininput:")
				registrationLoginDialog(win)
			}),
			widget.NewButton("LOGOUT", func() {
				fmt.Println("LOGOUT:")
				daemon.PutActiveloginname("")
			}),

		)
	*/
	registrationloginContainer := container.NewHBox(
		widget.NewButton("LOGIN AS", func() {
			fmt.Println("Entred registrationlogininput:")
			registrationLoginDialog(win)
		}),
		widget.NewButton("LOGOUT", func() {
			fmt.Println("LOGOUT:")
			daemon.PutActiveloginname("")
		}),
	)

	registrationloginentryContainer := container.NewBorder(nil, nil, nil, registrationloginContainer, registrationloginlabel)
	//-////////////////////////////
	contentdirectorystr := binding.NewString()

	////////////////////////
	//
	go func() {
		for {
			activeloginname := daemon.GetActiveloginname()
			registerednamecontentpath, _ := daemon.GetCachedDirPathForRegistredName(activeloginname)
			if registerednamecontentpath != "" {
				//registerednamesdescriptionarray[i]+=" CONTENT DIRECTORY: "+registerednamepath
				contentdirectorystr.Set(" CONTENT DIRECTORY: " + registerednamecontentpath)
			} else if activeloginname != "" {
				contentdirectorystr.Set(" CONTENT DIRECTORY: ")
			} else {
				contentdirectorystr.Set("")
			}
			time.Sleep(time.Second * 2)

		}
	}()
	////////////////////////
	contentdirectorylabel := widget.NewLabelWithData(contentdirectorystr)

	//contentdirectoryContainer :=container.New(  layout.NewGridWrapLayout(fyne.NewSize(150, 40)),contentdirectorylabel)
	contentdirectoryContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(150, 40)), contentdirectorylabel)
	//
	//-////////////////////////////
	/*
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
				selectednameregistrationid=id

			}

			//
			go func() {
				for {

					//fmt.Println("*******!!!!!!!!",registerednames)
					registerednamesdescriptionarray:=daemon.Wlt.GetRegisteredNames()
					for i,registeredname:= range registerednamesdescriptionarray{
						registerednamepath,_:=daemon.GetCachedDirPathForRegistredName(registeredname)
						if registerednamepath!=""{
							registerednamesdescriptionarray[i]+=" CONTENT DIRECTORY: "+registerednamepath
						}

					}
					registerednames.Set(registerednamesdescriptionarray)
					time.Sleep(time.Second * 2)



				}
			}()
	*/

	nameregistrationbutton := widget.NewButton("NEW NAME REGISTRATION", func() {
		fmt.Println("creating a new name :")
		requestNameRegistrationDialog(win)
	})
	//nameregistrationbuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 40)),nameregistrationbutton)
	//
	//
	nameunregistrationbutton := widget.NewButton("NAME UNREGISTRATION", func() {

		fmt.Println("name unregistration")
		//
		//
		//registerednamesarray:=daemon.Wlt.GetRegisteredNames()

		err := cli.Sendnameunregistration(daemon.Wireswarm, daemon.Mn, daemon.Wlt, daemon.GetActiveloginname())
		if err != nil {
			dialog.ShowError(err, win)
		} else {
			dialog.ShowInformation("Name Unregistration", "Name unregistration is being broadcasted", win)

		}

	})
	setcontentdirbutton := widget.NewButton("SET CONTENT DIRECTORY", func() {
		fmt.Println("set content directory")
		selectContentDirDialog(win)
	})
	commitcontentdirbutton := widget.NewButton("COMMIT CONTENT DIRECTORY", func() {
		fmt.Println("commit content directory")
		commitContentDirDialog(win)
	})
	//nameunregistrationbuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(250, 40)),nameunregistrationbutton)
	//setcontentdirbuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(250, 40)),setcontentdirbutton)
	//commitcontentdirbuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(250, 40)),commitcontentdirbutton)

	//registrationbuttonsContainer:=container.NewHBox(nameunregistrationbuttoncontainer,setcontentdirbuttoncontainer,commitcontentdirbuttoncontainer)
	//registrationbuttonsContainer := container.NewBorder(nil, nil, nameunregistrationbutton,commitcontentdirbutton, setcontentdirbutton)
	registrationbuttonsContainer := container.NewVBox(nameunregistrationbutton, setcontentdirbutton, commitcontentdirbutton)
	layout := container.NewVBox(nameregistrationbutton, registrationloginentryContainer, contentdirectoryContainer /*,registrednameslistcontainer*/, registrationbuttonsContainer)
	return layout

}
func requestNameRegistrationDialog(win fyne.Window) {
	requestedname := widget.NewEntry()
	depositamount := widget.NewEntry()
	//contactname.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")
	depositamount.Validator = validation.NewRegexp(`^[0-9]+$`, "deposited amount can only contain numbers")
	items := []*widget.FormItem{
		widget.NewFormItem("Requested Name", requestedname),
		widget.NewFormItem("Deposit Amount", depositamount),
	}

	dialog.ShowForm("Inorder to proceed with the name registration, please provide the following:    ", "Okay  ", "Cancel", items, func(b bool) {
		if !b {
			fmt.Println("canceled")
			//nowalletFoundDialog(win,"")
			return
		}
		if b {

			fmt.Println("text")
			err := cli.Sendnameregistration(daemon.Wireswarm, daemon.Mn, daemon.Wlt, requestedname.Text, depositamount.Text)
			if err != nil {
				dialog.ShowError(err, win)
			} else {
				dialog.ShowInformation("Name Registration", "Name registration is being broadcasted", win)

			}
		}
	}, win)

}

////////////////////////////////
func selectContentDirDialog(win fyne.Window) {
	folderd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {

		folderpath := uri.Path()

		//registerednamesdescriptionarray:=daemon.Wlt.GetRegisteredNames()
		//fmt.Printf("Putting content folder path %s for selectednameregistrationid %d to registred name %s \n",folderpath,selectednameregistrationid,registerednamesdescriptionarray[selectednameregistrationid])
		daemon.PutCachedDirPathForRegistredName(daemon.GetActiveloginname(), folderpath)
		//os.Exit(0)

	}, win)
	folderd.Show()
}

////////////////////////////////
func commitContentDirDialog(win fyne.Window) {

	registrednamecommittxfeeEntry := widget.NewEntry()
	registrednamecommittxfeeEntry.Text = fmt.Sprintf("%d", daemon.Usersettings.Nameregistrationtxfee)
	registrednamecommittxfeeEntry.Validator = validation.NewRegexp(`^[0-9]+$`, "only numbers")
	//registrednamecommittxfeeSetDefaultButton:= widget.NewButton("Set default", func() {
	//	registrednamecommittxfeeEntry.SetText(fmt.Sprintf("%d",daemon.Usersettings.Nameregistrationtxfee))
	//})
	//nameregistrationtxfeeSetDefaultButtonContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(100, 40)),nameregistrationtxfeeSetDefaultButton)
	registrednamecommittxfeeEntryContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 40)), registrednamecommittxfeeEntry)
	registrednamecommittxfeeLabelContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(250, 40)), widget.NewLabel("Registred Name Commit Fee"))
	registrednamecommittxfeeContainer := container.NewHBox(registrednamecommittxfeeLabelContainer, registrednamecommittxfeeEntryContainer)
	/////////
	selectWalletFileCallback := func(response bool) {
		fmt.Println("Responded with", response)
		if response {
			//registerednamesdescriptionarray:=daemon.Wlt.GetRegisteredNames()
			//selectedregistredname:=registerednamesdescriptionarray[selectednameregistrationid]
			selectedregistredname := daemon.GetActiveloginname()
			path, _ := daemon.GetCachedDirPathForRegistredName(selectedregistredname)
			cerr := cli.Sendregistrednamecontentcommit(daemon.Wireswarm, daemon.Mn, daemon.Wlt, selectedregistredname, path, registrednamecommittxfeeEntry.Text)
			//Sendregistrednamecontentcommit(ws *wire.Swarm,mn *mainchain.Maincore,wlt *wallet.Wallet,namestring string,contentfolderstring string,amountfeestring string){
			if cerr != nil {
				dialog.ShowError(cerr, win)
			} else {
				dialog.ShowInformation("Registred Name Content", "Registred name content is being broadcasted", win)
			}

		}

	}
	/////////

	cnf := dialog.NewCustomConfirm("Would you like to send a commit transaction fee ?", "Yes ", "No  ", registrednamecommittxfeeContainer, selectWalletFileCallback, win)
	cnf.Show()
}
