package gui


import (
	//"log"
	"fmt"
	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	//"net/url"
	//"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2"
	"time"

	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	//"fyne.io/fyne/v2/widget"
	"github.com/globaldce/globaldce-gateway/cli"
	"github.com/globaldce/globaldce-gateway/daemon"
)

func balanceScreen() fyne.CanvasObject {
	tabs := container.NewAppTabs(
		//container.NewTabItem("Balance",  welcomeScreen()),
		//container.NewTabItemWithIcon("Balance",theme.DocumentSaveIcon(),  listTxsScreen()),
		container.NewTabItem("Overview", overviewScreen()),
		container.NewTabItem("Addresses",  addressesScreen()),

	)
	tabs.SetTabLocation(container.TabLocationTop)
	return tabs
}

func overviewScreen() fyne.CanvasObject {

	balancestr := binding.NewString()
	balancestr.Set("")
	balancetext := widget.NewLabelWithData(balancestr)

	//syncingstr := binding.NewString()
	//syncingstr.Set("")
	//syncingtext := widget.NewLabelWithData(syncingstr)
	input:=widget.NewLabel("globaldce gateway "+cli.AppVersion)
	
	label := container.NewVBox(
		widget.NewLabel(""),
			input,balancetext,)

	hbox := container.NewVBox( widget.NewIcon(nil), label)	

	go func() {
		for {	
			//fmt.Println("*******",daemon.Wlt.ComputeBalance())
			walletpathstr:=fmt.Sprintf("Wallet path: %s",daemon.MainwalletFilePath)
			walletbalancestr:=fmt.Sprintf("Wallet balance is %f", float64(daemon.Wlt.ComputeBalance()/1000000.0))
			syncingstr:=""
			if daemon.Wireswarm.Syncingdone{
				syncingstr="SYNCING DONE"
			}
			if daemon.Miningrunning{
				syncingstr="CPU MINING RUNNING"
			}
			balancestr.Set(syncingstr+"\n"+walletpathstr+"\n"+walletbalancestr+"\n\n\n\n\n\n\n\n\n")
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
	return container.NewHSplit( container.NewCenter(hbox),assestsdestailslist)
}

func addressesScreen() fyne.CanvasObject {
	/*data := make([]string, 1000)
	for i := range data {
		data[i] = fmt.Sprintf("Adr %d", i)
	}*/
	//////////////////////
	//icon := widget.NewIcon(nil)
	input := widget.NewEntry()
	input.SetPlaceHolder("Selected address ...")
	
	input.Disable()
	saveButton:= widget.NewButton("Copy", func() {
		fmt.Println("Content was:", input.Text)
		clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
		clipboard.SetContent(input.Text)
	})

	label := container.NewVBox(
			widget.NewLabel("Selected address                                           "),
				input,)
	label.Add(saveButton)
	hbox := container.NewVBox( widget.NewIcon(nil), label)

	///////////////////////
	wltaddresses := binding.BindStringList(
		&[]string{},
	)
	
	//fmt.Printf("%v",registerednames)
	wltaddresseslist := widget.NewListWithData(wltaddresses,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
		go func() {
			for {
				//fmt.Println("*******!!!!!!!!",registerednames)
				time.Sleep(time.Second * 2)
				updatedassestsdestails:=daemon.Wlt.GetAddressesDetails()
				if updatedassestsdestails!=nil{
					wltaddresses.Set(updatedassestsdestails)
				}
				
				//str.Set(fmt.Sprintf("WALLET BALANCE is %d", daemon.Wlt.ComputeBalance()))
				
			}
		}()
	///////////////////////
	
	/*
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
	*/
	
	wltaddresseslist.OnSelected = func(id widget.ListItemID) {
		//label.SetText(data[id])
		//icon.SetResource(theme.DocumentIcon())
	//input := widget.NewEntry(
	textvalue,_:=wltaddresses.GetValue(id)
	input.Text=textvalue
	input.Disable()
	input.Refresh()
	//label = container.NewVBox(input, widget.NewButton("Save", func() {
	//	fmt.Println("Content was:", input.Text)
	//}))

	}
	

	return container.NewHSplit( container.NewCenter(hbox),wltaddresseslist)
}


