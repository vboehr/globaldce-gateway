package gui

import (
	"path/filepath"
	"net/url"
	"fmt"
	"encoding/json"
	"time"
	//"strings"
	"fyne.io/fyne/v2"
	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/data/binding"
	"github.com/globaldce/globaldce-toolbox/mainchain"
	"github.com/globaldce/globaldce-toolbox/daemon"
	"github.com/globaldce/globaldce-toolbox/utility"
)

//////////////////////////////////////////////////
//////////////////////////////////////////////////
func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}
func displayPostDetails(p *post) func() {
	return func() {
		if p!=nil{
			fmt.Println("post",p)
			win := guiApp.NewWindow("Public post details")
			win.SetContent(widget.NewLabel(p.Content))
			win.Resize(fyne.NewSize(200, 200))
			win.Show()

		}
	  
	}
  }


const iconSize = float32(100)
type post struct {
	Name string
	Link string
	Content string
	AttachmentSizeArray []int
	AttachmentHashArray []utility.Hash
	//user    *user
}

func PostInfoFromString(s string) mainchain.PostInfo{
	var p mainchain.PostInfo
	json.Unmarshal([]byte(s), &p)
	return p
}
func DataFilePathFromHash(h utility.Hash) string{
	var s string//
	s=filepath.Join("Data","DataFiles",fmt.Sprintf("%x",h))
	//s="./Data/DataFiles/77d730d337fdb932c267d602343c0f2cb271c3572a5a4fd91f083cf66aae83a6"
	return s
}

type postRenderer struct {
	m         *postCell
	top, main *widget.Label
	pic       *canvas.Image
	link 	  *widget.Hyperlink
	details   *widget.Button
	sep       *widget.Separator
}
func (m *postRenderer) Destroy() {
}
func (m *postRenderer) Layout(s fyne.Size) {
	remainWidth := s.Width - iconSize - theme.Padding()*2
	remainStart := iconSize + theme.Padding()*2
	m.pic.Resize(fyne.NewSize(iconSize, iconSize))
	m.pic.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
	m.details.Resize(fyne.NewSize(iconSize, iconSize/3))
	m.details.Move(fyne.NewPos(theme.Padding(), theme.Padding()+iconSize))
	m.top.Move(fyne.NewPos(remainStart, -theme.Padding()))
	m.top.Resize(fyne.NewSize(remainWidth, m.top.MinSize().Height))

	m.link.Move(fyne.NewPos(remainStart, m.top.MinSize().Height-theme.Padding()*4))//100 is the height of the cell
	if m.m.msg.Link!=""{
		m.link.Resize(fyne.NewSize(remainWidth, m.top.MinSize().Height))
	}
	

	m.main.Move(fyne.NewPos(remainStart, m.top.MinSize().Height+2*theme.Padding()))
	m.main.Resize(fyne.NewSize(remainWidth, m.main.MinSize().Height))
	
	m.sep.Move(fyne.NewPos(0, s.Height-theme.SeparatorThicknessSize()))
	m.sep.Resize(fyne.NewSize(s.Width, theme.SeparatorThicknessSize()))
}
func (m *postRenderer) MinSize() fyne.Size {
	s1 := m.top.MinSize()
	s2 := m.main.MinSize()
	w := fyne.Max(s1.Width, s2.Width)
	//return fyne.NewSize(w+iconSize+theme.Padding()*2,
	//	s1.Height+s2.Height-theme.Padding()*4)
	_=w
	return fyne.NewSize(1000,200)
}
func (m *postRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{m.top, m.main, m.pic,m.link,m.details, m.sep}
}

func (m *postRenderer) Refresh() {
	m.top.SetText(m.m.msg.Name)
	m.details= widget.NewButton("Click me", displayPostDetails(m.m.msg))
	///////////////////////////////////
	//m.pic.SetResource(theme.FyneLogo())
	if m.m.msg.AttachmentHashArray!=nil{
		m.pic=canvas.NewImageFromFile(DataFilePathFromHash(m.m.msg.AttachmentHashArray[0]))
	}
	
	///////////////////////////////////
	//m.pic.SetResource(nil)
	m.main.SetText(m.m.msg.Content)
	
	m.link=widget.NewHyperlink(m.m.msg.Link, parseURL(m.m.msg.Link))
	//fmt.Printf("link is %s",m.m.msg.Link)
	/*
	if m.m.msg.user.name != "" {
		m.top.SetText(m.m.msg.user.name)
	} else {
		m.top.SetText(m.m.msg.user.username)
	}
	m.main.SetText(m.m.msg.Content)
	go m.pic.SetResource(m.m.avatarResource())
	*/
	
}
type postCell struct {
	widget.BaseWidget
	msg *post
}

func (m *postCell) CreateRenderer() fyne.WidgetRenderer {
	name := widget.NewLabelWithStyle(m.msg.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	name.Wrapping = fyne.TextTruncate
	body := widget.NewLabel(m.msg.Content)
	body.Wrapping = fyne.TextWrapWord
	emptybutton:=widget.NewButton("", displayPostDetails(nil))
	emptylink:=widget.NewHyperlink(m.msg.Link, parseURL(m.msg.Link))
	return &postRenderer{m: m,
		top:  name,//canvas.NewImageFromFile("./rawtest/unnamed.jpg")
		//main: body, pic: widget.NewIcon(nil),link:emptylink, sep: widget.NewSeparator()}
		main: body, pic:canvas.NewImageFromFile(""),link:emptylink,details:emptybutton, sep: widget.NewSeparator()}
}

func (m *postCell) UpdatePost(s string)  {
	//m.Unbind()
	p:=PostInfoFromString(s)
	m.msg.Name=p.Name
	m.msg.Link=p.Link
	m.msg.Content=p.Content
	m.msg.AttachmentSizeArray=p.AttachmentSizeArray
	m.msg.AttachmentHashArray=p.AttachmentHashArray

}
func newPostCell(m *post) *postCell {
	ret := &postCell{msg: m}
	ret.ExtendBaseWidget(ret)
	return ret
}

//////////////////////////////////////////////////
//////////////////////////////////////////////////
//var bindings []binding.DataMap

var MaxDisplayedPost =50
var searchtext string =""
func exploreScreen(w fyne.Window)  fyne.CanvasObject{
	bindings := binding.BindStringList(
		&[]string{},
	)
	//var newbindings []binding.DataMap
	//bindings=&newbindings

//for _, todo := range data {
//  bindings = append(bindings, binding.BindStruct(&todo))
//}
//getPosts("")



list := widget.NewTable(
  func() (int, int) {
    //return len(bindings), 1
	l,_:=bindings.Get()
	return len(l), 1
  },
  func() fyne.CanvasObject {
    //return widget.NewLabel("wide Content")
	m2:=&post{Content:"TMPCONTENT"}
	return newPostCell(m2)
  },
  func(i widget.TableCellID, o fyne.CanvasObject) {
    //title, _ := (bindings)[i.Row].GetItem("Name")

    //log.Println(title)
    //o.(*widget.Label).Bind(title.(binding.String))
	///////////////////////////////////////////////
	//o.(*postCell.Label).Bind(title.(binding.String))
	//_=title
	//o.(*postCell).Cool("xxxx")
	bs,_:=bindings.GetValue(i.Row)
	o.(*postCell).UpdatePost(bs)
	//fmt.Printf("*********************bs",bs)
	o.(*postCell).Refresh()

	//////////////////////////////////////
  })



go func() {
	for {
		//fmt.Println("*******!!!!!!!!",registerednames)
		time.Sleep(time.Second * 2)
		//assestsdestails.Set(daemon.Wlt.GetAssetsDetails())
		bindings.Set(getPosts(searchtext))
		list.Refresh()
		//str.Set(fmt.Sprintf("WALLET BALANCE is %d", daemon.Wlt.ComputeBalance()))
		
	}
}()
	
	
	searchentry:=widget.NewEntry()
	searchentry.SetPlaceHolder("Seach for ...")
	searchentry.OnSubmitted=func(s string) {
		fmt.Println("Search submitted",s)
		searchtext=s
		bindings.Set(getPosts(s))
		list.Refresh()
	}


	return container.NewBorder(searchentry, nil, nil, nil,list)

}

func getPosts(keywords string)[]string{
	_=keywords
	//GetPostInfoStringArray []string
	/*
	var sarray []string
	
	s1:=mainchain.StringFromPostInfo(mainchain.PostInfo{Name:"cool1",Link:"",Content:"11111cool1 content text"})
	//if strings.Index(s1,keywords)>=0{
		sarray=append(sarray,s1)
	//}
	s2:=mainchain.StringFromPostInfo(mainchain.PostInfo{Name:"cool2",Link:"https://www.google.com",Content:"cool2222 content text"})
	//if strings.Index(s2,keywords)>=0{
		sarray=append(sarray,s2)
	//}
	
	sarray=append(sarray,mainchain.StringFromPostInfo(mainchain.PostInfo{Name:"cool33",Link:"",Content:"cool33 content text"}))
	//bindings.Set(sarray)
	return sarray
	*/
	
	if daemon.Mn==nil{
		return nil
	}
	sarray:=daemon.Mn.GetPostInfoStringArray(30) 
	fmt.Printf("********%v",sarray)
	return sarray
	
	//return nil
}

