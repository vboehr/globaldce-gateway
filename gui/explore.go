package gui

import (
	"net/url"
	"fmt"
	"encoding/json"
	//"time"
	"fyne.io/fyne/v2"
	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	//"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/data/binding"
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

const iconSize = float32(100)
type post struct {
	Name string
	Content string
	//user    *user
}

func StringFromPost(p post) string{

	//json.Unmarshal([]byte(stringData), &data)
	b,_:=json.Marshal(p)
	return string(b)
	
}
func PostFromString(s string) post{
	var p post
	json.Unmarshal([]byte(s), &p)
	return p
}

type postRenderer struct {
	m         *postCell
	top, main *widget.Label
	pic       *widget.Icon
	link 	  *widget.Hyperlink
	sep       *widget.Separator
}
func (m *postRenderer) Destroy() {
}
func (m *postRenderer) Layout(s fyne.Size) {
	remainWidth := s.Width - iconSize - theme.Padding()*2
	remainStart := iconSize + theme.Padding()*2
	m.pic.Resize(fyne.NewSize(iconSize, iconSize))
	m.pic.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
	m.top.Move(fyne.NewPos(remainStart, -theme.Padding()))
	m.top.Resize(fyne.NewSize(remainWidth, m.top.MinSize().Height))

	m.link.Move(fyne.NewPos(remainStart,150 -theme.Padding()))//100 is the height of the cell
	m.link.Resize(fyne.NewSize(remainWidth, m.top.MinSize().Height))

	m.main.Move(fyne.NewPos(remainStart, m.top.MinSize().Height-theme.Padding()*4))
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
	return fyne.NewSize(500,200)
}
func (m *postRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{m.top, m.main, m.pic,m.link, m.sep}
}

func (m *postRenderer) Refresh() {
	m.top.SetText(m.m.msg.Name)
	m.pic.SetResource(theme.FyneLogo())
	m.main.SetText(m.m.msg.Content)
	
	m.link=widget.NewHyperlink("linktext", parseURL("https://www.github.com/globaldce"))
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
	emptylink:=widget.NewHyperlink("", parseURL(""))
	return &postRenderer{m: m,
		top:  name,
		main: body, pic: widget.NewIcon(nil),link:emptylink, sep: widget.NewSeparator()}
}

func (m *postCell) UpdatePost(s string)  {
	//m.Unbind()
	p:=PostFromString(s)
	m.msg.Name=p.Name
	m.msg.Content=p.Content

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
bindings.Set(getPosts(""))
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


	
	
	searchentry:=widget.NewEntry()
	searchentry.SetPlaceHolder("Seach for ...")
	searchentry.OnSubmitted=func(s string) {
		fmt.Println("Search submitted",s)
		
		bindings.Set(getPosts(s))
	}


	return container.NewBorder(searchentry, nil, nil, nil,list)

}

func getPosts(keywords string)[]string{

	_=keywords
	var sarray[]string
	//if keywords=="1"{
		sarray=append(sarray,StringFromPost(post{Name:"cool1",Content:"11111cool1 content text"}))
	//}
	//if keywords=="2"{
		sarray=append(sarray,StringFromPost(post{Name:"cool2",Content:"cool2222 content text"}))
	//}
	
	sarray=append(sarray,StringFromPost(post{Name:"cool33",Content:"cool33 content text"}))
	//bindings.Set(sarray)
	return sarray
}

