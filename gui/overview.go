package gui


import (
	//"log"
	"fmt"
	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	///////////////
	"net/url"
	"fyne.io/fyne/v2/layout"
	///////////////
	"fyne.io/fyne/v2"
	"time"

	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	//"fyne.io/fyne/v2/widget"
	//"github.com/globaldce/globaldce-gateway/mainchain"
	"github.com/globaldce/globaldce-gateway/daemon"
)
/*
func balanceScreen(win fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		//container.NewTabItem("Balance",  welcomeScreen()),
		//container.NewTabItemWithIcon("Balance",theme.DocumentSaveIcon(),  listTxsScreen()),
		container.NewTabItem("Overview", overviewScreen(win)),
		container.NewTabItem("Addresses",  addressesScreen()),

	)
	tabs.SetTabLocation(container.TabLocationTop)
	return tabs
}
*/

//var recentdappnamesarray []string
var dappnameselectedid int =-1
//var daemon.Usersettings.Activeloginname string=""

func overviewScreen(win fyne.Window) fyne.CanvasObject {
	
	balancestr := binding.NewString()
	balancestr.Set("")
	balancetext := widget.NewLabelWithData(balancestr)

	//syncingstr := binding.NewString()
	//syncingstr.Set("")
	//syncingtext := widget.NewLabelWithData(syncingstr)
	//input:=widget.NewLabel("globaldce gateway "+cli.AppVersion)
	
	/////////////////////////
	registrationloginstr:= binding.NewString()

	////////////////////////
			//
			go func() {
				for {
					
					if daemon.GetActiveloginname()==""{
						registrationloginstr.Set("Not logged in")
					} else {
						registrationloginstr.Set(daemon.GetActiveloginname())
					}
					time.Sleep(time.Second * 2)

					
				}
			}()
	////////////////////////
	
	registrationloginlabel := widget.NewLabelWithData(registrationloginstr)

	registrationlogininputContainer :=container.New(  layout.NewGridWrapLayout(fyne.NewSize(300, 40)),registrationloginlabel)
	registrationloginentryContainer := container.NewHBox(registrationlogininputContainer, 
	widget.NewButton("LOGIN AS", func() {
		fmt.Println("Entred registrationlogininput:")
		registrationLoginDialog(win)
	}),
	widget.NewButton("LOGOUT", func() {
		fmt.Println("LOGOUT:")
		//daemon.Usersettings.Activeloginname=""
		daemon.PutActiveloginname("")
	}),

	)
	
	//dappnameinputstr := binding.NewString()
	
	dappnameinput := widget.NewEntry()
	//dappnameinput.Bind(dappnameinputstr)
	dappnameinput.SetPlaceHolder("Dapp name to be loaded ...")
	dappnameinputContainer:=container.New(  layout.NewGridWrapLayout(fyne.NewSize(400, 40)),dappnameinput)
	dappentryContainer := container.NewHBox(dappnameinputContainer, widget.NewButton("LOAD", func() {
		fmt.Println("Entred dapp name:", dappnameinput.Text)
		dappnameinputText:=dappnameinput.Text
		serr:=daemon.Mn.ServeContent(dappnameinputText)
		if serr==nil{
			//dappnameinputstr.Set(dappnameinputText)
			//mainchain.HandleWebSocket(dappnameinputText)
			u, err := url.Parse("http://localhost:8080/"+dappnameinputText+"/index.html")//("http://localhost:8080/"+dappnameinputText+"/index.html")//("./Cache/Content/dapptest/index.html")//("https://fyne.io/")
			_=err
			guiApp.OpenURL(u)
			////////////////
			daemon.AddToRecentDappNames(dappnameinputText)
		}

		//recentdappnamesarray=append([]string{dappnameinputText}, recentdappnamesarray ...)
	}))
	/////////////////////////

	/////////////////////////////////////////////
	/////////////////////////////////////////////
	recentdappnames := binding.BindStringList(
		&[]string{},
	)
	
	fmt.Printf("%v",recentdappnames)
	recentdappnameslist := widget.NewListWithData(recentdappnames,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
		//
		recentdappnameslist.OnSelected = func(id widget.ListItemID) {
			//label.SetText(data[id])
			//icon.SetResource(theme.DocumentIcon())
		//input := widget.NewEntry()
		textvalue,_:=recentdappnames.GetValue(id)
		//dappnameselected:=textvalue
		dappnameinput.Text=textvalue
		fmt.Printf("dappnameselected %d\n",int(id))
		dappnameselectedid=int(id)
		dappnameinput.Refresh()
		}
		//
		go func() {
			for {
				//fmt.Println("*******!!!!!!!!",recentdappnames)
				
				//recentdappnames.Set(recentdappnamesarray)
				recentdappnames.Set(daemon.GetRecentDappNames())

				time.Sleep(time.Second * 2)
				//str.Set(fmt.Sprintf("WALLET BALANCE is %d", daemon.Wlt.ComputeBalance()))
				
			}
		}()
	

	recentdappnamesclearbutton:= widget.NewButton("CLEAR", func() {
		daemon.ClearRecentDappNameWithId(dappnameselectedid)
		//if (dappnameselectedid<0)||(dappnameselectedid>=len(recentdappnamesarray)){
		//	return
		//}
		//recentdappnamesarray=append(recentdappnamesarray[:dappnameselectedid], recentdappnamesarray[dappnameselectedid+1:] ...)
    })
	recentdappnamesclearallbutton:= widget.NewButton("CLEAR ALL", func() {
		daemon.ClearAllRecentDappNames()
		//recentdappnamesarray=make([]string, 0)
    })
	recentdappnamesclearbuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(100, 40)),recentdappnamesclearbutton)
	//
	//
	recentdappnamesclearallbuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(100, 40)),recentdappnamesclearallbutton)
	//nameunregistrationbuttoncontainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(350, 40)),recentdappnamesclearbuttoncontainer)
	//layout:=container.New(layout.NewPaddedLayout(),container.NewVBox(recentdappnameslist,nameregistrationcontainer))
	recentdappnameslistcontainer:=container.New(layout.NewGridWrapLayout(fyne.NewSize(300, 200)),recentdappnameslist)
	recentdappnamesclearbuttonscontainer := container.NewHBox(recentdappnamesclearbuttoncontainer,recentdappnamesclearallbuttoncontainer)
	recentdappnamesContainer:=container.NewVBox(recentdappnameslistcontainer,recentdappnamesclearbuttonscontainer)
	/////////////////////////////////////////////
	/////////////////////////////////////////////
	label := container.NewVBox(
		//widget.NewLabel(""),
			balancetext,
			registrationloginentryContainer,
			dappentryContainer,
			recentdappnamesContainer,
			widget.NewLabel(" "),
		)
	hbox:=label
	//hbox := container.NewVBox( widget.NewIcon(nil), label)	

	go func() {
		for {	
			//fmt.Println("*******",daemon.Wlt.ComputeBalance())
			walletpathstr:=fmt.Sprintf("Wallet path: %s",daemon.MainwalletFilePath)
			var walletbalancestr string
			if daemon.Walletinstantiated{
				if daemon.Wlt.Walletstate=="" {
					walletbalancestr=fmt.Sprintf("Wallet balance is %f", float64(daemon.Wlt.ComputeBalance()/1000000.0))
				} else {
					walletbalancestr=daemon.Wlt.Walletstate
				}
			}

			
			syncingstr:=""
			if daemon.Wireswarm!=nil{
				if daemon.Wireswarm.Syncingdone {
					syncingstr="SYNCING DONE"
				}
			}

			if daemon.Miningrunning{
				syncingstr="CPU MINING RUNNING"
			}
			balancestr.Set(syncingstr+"\n"+walletpathstr+"\n"+walletbalancestr+"")
			time.Sleep(time.Second * 2)
		}
	}()

		///////////////////////
		///////////////////////
		assestsdestails := binding.BindStringList(
			&[]string{},
		)
		
		//fmt.Printf("%v",registerednames)
		assestsdestailslist := widget.NewListWithData(assestsdestails,
			func() fyne.CanvasObject {
				return widget.NewLabel("template")
			},
			func(i binding.DataItem, o fyne.CanvasObject) {
				o.(*widget.Label).Bind(i.(binding.String))
			})
			go func() {
				for {
					//fmt.Println("*******!!!!!!!!",registerednames)
					time.Sleep(time.Second * 3)
					updatedassestsdestails:=daemon.Wlt.GetAssetsDetails()
					if updatedassestsdestails!=nil{
						assestsdestails.Set(updatedassestsdestails)
					}
					time.Sleep(time.Second * 5)
					//str.Set(fmt.Sprintf("WALLET BALANCE is %d", daemon.Wlt.ComputeBalance()))
					
				}
			}()
			/*
			assestsdestailslist.OnSelected = func(id widget.ListItemID) {
				//label.SetText(data[id])
				//icon.SetResource(theme.DocumentIcon())
			//input := widget.NewEntry(
			textvalue,_:=assestsdestails.GetValue(id)
			input.Text=textvalue
			//input.Disable()
			input.Refresh()
			//label = container.NewVBox(input, widget.NewButton("Save", func() {
			//	fmt.Println("Content was:", input.Text)
			//}))
		
			}
			*/
		///////////////////////
		
	/*
	data := make([]string, 1000)
	for i := range data {
		data[i] = fmt.Sprintf("Tx %d", i)
	}



	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id])
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
	
		input.Text=data[id]

		input.Refresh()
	
	}
	*/
	//
	//containerhbox:=container.New(layout.NewBorderLayout(nil, nil, hbox, nil))
	//return container.NewHSplit(containerhbox,assestsdestailslist)

	return container.NewVSplit( container.NewCenter(hbox),assestsdestailslist)
}
