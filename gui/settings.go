package gui


import (
	"log"
	//"fmt"
	//"image/color"
	//"fyne.io/fyne/v2/app"
	//"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/theme"
	//"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	//"fyne.io/fyne/v2/canvas"
	"net/url"
	//"net/url"
	//"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2"
	"github.com/globaldce/globaldce-toolbox/cli"
)





func settingsScreen() fyne.CanvasObject {

	entrymainwalletpath := widget.NewEntry()
	entrymainwalletpath.Text=cli.Usersettings.MainwalletFilePath
	//entry1.Validator = validation.NewRegexp(`^[0-9]+$`, "deposited amount can only contain numbers")
	//entry2 := widget.NewEntry()
	//entry2.Validator = validation.NewRegexp(`^[0-9]+$`, "deposited amount can only contain numbers")

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "main wallet file path", Widget:entrymainwalletpath},
			//{Text: "Cool2", Widget: entry2},
		},
		OnSubmit: func() { // optional, handle form submission
			log.Println("Main wallet file path set to:", entrymainwalletpath.Text)
			//log.Println("entry 2:", entry2.Text)
			cli.Usersettings.MainwalletFilePath=entrymainwalletpath.Text
		},
	}

	return form

}

//func settingsScreen() fyne.CanvasObject {
//	selectEntry := widget.NewSelectEntry([]string{"Option A", "Option B", "Option C"})
//	selectEntry.PlaceHolder = "Type or select"
//	disabledCheck := widget.NewCheck("Disabled check", func(bool) {})
//	disabledCheck.Disable()
//	radio := widget.NewRadioGroup([]string{"Radio Item 1", "Radio Item 2"}, func(s string) { fmt.Println("selected", s) })
//	radio.Horizontal = true
//	disabledRadio := widget.NewRadioGroup([]string{"Disabled radio"}, func(string) {})
//	disabledRadio.Disable()
//
//	return container.NewVBox(
//		widget.NewSelect([]string{"Option 1", "Option 2", "Option 3"}, func(s string) { fmt.Println("selected", s) }),
//		selectEntry,
//		widget.NewCheck("Check", func(on bool) { fmt.Println("checked", on) }),
//		disabledCheck,
//		radio,
//		disabledRadio,
//		widget.NewSlider(0, 100),
//	)
//}

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}






//func welcomeScreen() fyne.CanvasObject {
	//logo := canvas.NewImageFromResource(data.FyneScene)//
	//logo.FillMode = canvas.ImageFillContain
	//if fyne.CurrentDevice().IsMobile() {
	//	logo.SetMinSize(fyne.NewSize(171, 125))
	//} else {
	//	logo.SetMinSize(fyne.NewSize(228, 167))
	//}

//	return container.NewCenter(container.NewVBox(
//		widget.NewLabelWithStyle("Welcome to the Fyne toolkit demo app", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		//logo,
//		container.NewHBox(
//			widget.NewHyperlink("fyne.io", parseURL("https://fyne.io/")),
//			widget.NewLabel("-"),
//			widget.NewHyperlink("documentation", parseURL("https://fyne.io/develop/")),
//			widget.NewLabel("-"),
//			widget.NewHyperlink("sponsor", parseURL("https://github.com/sponsors/fyne-io")),
//		),
//	))
//}

