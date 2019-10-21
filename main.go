package main

import (
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

var tagText = widget.NewLabel("")
var checkText = widget.NewLabel("")

func clickTag(tag string) {
	str := tagText.String()
	tagWithBracket := "[" + tag + "]"
	if strings.Contains(str, tagWithBracket) {
		str = strings.ReplaceAll(str, tagWithBracket+",", "")
	} else {
		str = str + "[" + tag + "]" + ","
	}
	tagText.SetText(str)
}

func showTags(box *widget.Box, tags []string) {
	rvVbox := widget.NewHBox()
	vv := tags[len(tags)-1]
	button := widget.NewButton(vv, func() {
		clickTag(vv)
	})
	rvVbox.Append(button)
	box.Append(rvVbox)
}

func main() {
	db := InitDB()
	orm, err := InitORM(db)
	if err != nil {
		return
	}
	app := app.New()

	w := app.NewWindow("Add tags to item and save to sqlite")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(900, 600))
	w.Canvas().SetScale(2)
	lHbox := widget.NewVBox()
	entry := widget.NewEntry()
	check := widget.NewButton("Check Item", func() { orm.checkItem(entry) })
	submit := widget.NewButton("Submit Item", func() { orm.submitItem(entry) })
	delete := widget.NewButton("Delete Item", func() { orm.deleteItem(entry) })
	lHbox.Append(entry)
	lHbox.Append(widget.NewHBox(check, submit, delete))
	lHbox.Append(checkText)
	lHbox.Append(tagText)

	tags := orm.getTags()
	rHbox := widget.NewVBox()
	itemEntry := widget.NewEntry()
	rHbox.Append(itemEntry)
	rHbox.Append(widget.NewButton("Add Tag", func() {
		hasAdded, _ := orm.addTag(itemEntry)
		if hasAdded {
			showTags(rHbox, orm.getTags())
		}
	}))
	rHbox.Append(widget.NewLabel("All Tags:"))
	scroller := widget.NewScrollContainer(rHbox)
	var rvHboxs []*widget.Box
	lineNumber := 6
	for index := 0; index < len(tags)/lineNumber+1; index++ {
		rvVbox := widget.NewHBox()
		rvHboxs = append(rvHboxs, rvVbox)
		rHbox.Append(rvVbox)
	}
	for i, v := range tags {
		vv := v
		rvHboxs[i/lineNumber].Append(
			widget.NewButton(v, func() {
				clickTag(vv)
			}),
		)
	}

	w.SetContent(widget.NewHBox(
		scroller,
		lHbox,
	))

	w.ShowAndRun()
}
