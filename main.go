package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"sort"
	"strings"
)

var (
	text         string
	inTE, outTE  *walk.TextEdit
	lineA, lineB *walk.LineEdit
	mw           *walk.MainWindow
)

type Res struct {
	Key  string
	Time int
	Rate float64
	Len  int
}
type ResSlice []Res

func (a ResSlice) Len() int {
	return len(a)
}
func (a ResSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ResSlice) Less(i, j int) bool {
	return a[j].Rate < a[i].Rate
}
func main() {

	err := MainWindow{
		AssignTo: &mw,
		Title:    "词频统计分析",
		Font: Font{
			PointSize: 12,
		},
		MinSize: Size{Width: 600, Height: 400},
		//Size:     Size{Width: 600, Height: 400},
		Layout: VBox{},

		Children: []Widget{

			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						TextAlignment: AlignCenter,
						Text:          "输入",
					},
					Label{
						TextAlignment: AlignCenter,
						Text:          "",
					},
					Label{
						TextAlignment: AlignCenter,
						Text:          "输出",
					},
				},
			},
			Composite{
				Layout: HBox{},
				//MinSize: Size{Width: 400, Height: 370},
				Children: []Widget{

					TextEdit{
						AssignTo:      &inTE,
						VScroll:       true,
						OnTextChanged: SaveText,
					},
					Composite{
						Layout: VBox{},
						Children: []Widget{

							PushButton{
								Text:      "统计频",
								OnClicked: GetRate,
							},

							Label{
								TextAlignment: AlignCenter,
								Text:          "替换",
							},

							LineEdit{

								AssignTo: &lineA,
							},
							Label{
								TextAlignment: AlignCenter,
								Text:          "为",
							},
							LineEdit{

								AssignTo: &lineB,
							},
							Label{
								TextAlignment: AlignCenter,
								Text:          "-->",
							},
							PushButton{
								Text:      "替换字",
								OnClicked: ChangeText,
							},
						},
					},
					TextEdit{AssignTo: &outTE, ReadOnly: true, VScroll: true},
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					Label{
						Text:          "参考",
						TextAlignment: AlignCenter,
					},
					Label{
						Text:       "常用1字母组合：e  t  a",
						TextColor:  walk.RGB(0xff, 0, 0),
						Background: SolidColorBrush{Color: walk.RGB(0xf0, 0xf0, 0xf0)},
					},
					Label{
						Text:       "常用2字母组合：am  is  us  oh  hi  my  p.m  a.m  do  no  to  an  we  up  so  me",
						TextColor:  walk.RGB(0xff, 0, 0),
						Background: SolidColorBrush{Color: walk.RGB(0xf0, 0xf0, 0xf0)},
					},
					Label{
						Text:       "常用3字母组合：the  ing  and  her  ere",
						TextColor:  walk.RGB(0xff, 0, 0),
						Background: SolidColorBrush{Color: walk.RGB(0xf0, 0xf0, 0xf0)},
					},
				},
			},
		},
	}.Create()
	if err != nil {
		log.Fatal(err)
	}
	mw.Run()
}
func GetRate() {
	n := len(text)
	result := ""
	listOne := make(map[string]*Res, 1024)
	listTwo := make(map[string]*Res, 1024)
	listThree := make(map[string]*Res, 1024)
	totalOne := 0
	totalTwo := 0
	totalThree := 0
	for _, v := range text {
		if v >= 'A' && v <= 'Z' {
			_, ok := listOne[string(v)]
			if ok {
				listOne[string(v)].Time++
			} else {
				listOne[string(v)] = &Res{Len: 1, Time: 1}
			}
			totalOne++
		}
	}

	for i := 0; i < n-1; i++ {
		if text[i] >= 'A' && text[i] <= 'Z' && text[i+1] >= 'A' && text[i+1] < 'Z' {
			_, ok := listTwo[string(text[i])+string(text[i+1])]
			if ok {
				listTwo[string(text[i])+string(text[i+1])].Time++
			} else {
				listTwo[string(text[i])+string(text[i+1])] = &Res{Len: 1, Time: 1}
			}
			totalTwo++
		}
	}
	for i := 0; i < n-2; i++ {
		if text[i] >= 'A' && text[i] <= 'Z' && text[i+1] >= 'A' && text[i+1] < 'Z' && text[i+2] >= 'A' && text[i+2] < 'Z' {
			_, ok := listThree[string(text[i])+string(text[i+1])+string(text[i+2])]
			if ok {
				listThree[string(text[i])+string(text[i+1])+string(text[i+2])].Time++
			} else {
				listThree[string(text[i])+string(text[i+1])+string(text[i+2])] = &Res{Len: 1, Time: 1}
			}
			totalThree++
		}
	}
	for i, _ := range listOne {
		listOne[i].Rate = float64(listOne[i].Time) * 100 / float64(totalOne)
	}
	for i, _ := range listTwo {
		listTwo[i].Rate = float64(listTwo[i].Time) * 100 / float64(totalTwo)
	}
	for i, _ := range listThree {
		listThree[i].Rate = float64(listThree[i].Time) * 100 / float64(totalThree)
	}
	result += "连续1个字符:\r\n"
	result += mySort(listOne)

	result += "\r\n连续2个字符:\r\n"
	result += mySort(listTwo)

	result += "\r\n连续3个字符:\r\n"
	result += mySort(listThree)

	outTE.SetText(result)
}

func ChangeText() {
	a := lineA.Text()
	b := lineB.Text()
	if len(a) != len(b) {
		walk.MsgBox(mw, "Error", "替换字符个数不一致，请重新输入", walk.MsgBoxIconInformation)
		return
	}
	for i := range a {
		text = strings.Replace(text, string(a[i]), string(b[i]), -1)
	}
	outTE.SetText(text)
}

func SaveText() {
	text = inTE.Text()
}

func mySort(mp map[string]*Res) (result string) {
	var list []Res
	for i, v := range mp {
		list = append(list, Res{Key: i, Time: v.Time, Rate: v.Rate})
	}
	sort.Sort(ResSlice(list))
	for _, v := range list {
		result += fmt.Sprintf("%v:\t%v\t%.2f%% \r\n", v.Key, v.Time, v.Rate)
	}
	return result
}
