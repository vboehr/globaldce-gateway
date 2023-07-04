package gui

import (
	//"log"
	"fmt"
	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	///////////////
	///////////////
	"fyne.io/fyne/v2"
	"time"

	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	//"fyne.io/fyne/v2/widget"
	"github.com/globaldce/globaldce-gateway/daemon"
)

func receiveScreen() fyne.CanvasObject {
	/*data := make([]string, 1000)
	for i := range data {
		data[i] = fmt.Sprintf("Adr %d", i)
	}*/
	//////////////////////
	//icon := widget.NewIcon(nil)
	input := widget.NewEntry()
	input.SetPlaceHolder("Selected address ...")

	input.Disable()
	saveButton := widget.NewButton("Copy", func() {
		fmt.Println("Content was:", input.Text)
		clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
		clipboard.SetContent(input.Text)
	})

	label := container.NewVBox(
		widget.NewLabel("Selected address                                           "),
		input)
	label.Add(saveButton)
	genAddrButton := widget.NewButton("Generate New Address", func() {
		fmt.Println("New address generation ... ")
		if daemon.Walletinstantiated() {
			if daemon.Wlt.Walletloaded {
				daemon.Wlt.GenerateKeyPair()
			} else {
				fmt.Println("New address generation CANCELED - wallet not loaded")
			}
		}
	})
	hbox := container.NewVBox(widget.NewIcon(nil), label, genAddrButton)

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
			updatedassestsdestails := daemon.Wlt.GetAddressesDetails()
			if updatedassestsdestails != nil {
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
		textvalue, _ := wltaddresses.GetValue(id)
		input.Text = textvalue
		input.Disable()
		input.Refresh()
		//label = container.NewVBox(input, widget.NewButton("Save", func() {
		//	fmt.Println("Content was:", input.Text)
		//}))

	}

	return container.NewVSplit(hbox, wltaddresseslist)
}
