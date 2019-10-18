package main

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/widget"
	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// ORM warp for xorm
type ORM struct {
	E *xorm.Engine
}

// Item describes anything
type Item struct {
	ID    string `xorm:"notnull unique"`
	Name  string `xorm:"notnull unique"`
	TagID []string
}

// Tag describes tag
type Tag struct {
	ID    string `xorm:"notnull unique"`
	Title string `xorm:"notnull unique"`
}

func pe(s error) {
	print("[error]:")
	println(s)
}

// InitDB generate new db file
func InitDB() string {
	file := "./data.db"
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		os.Create(file)
	}
	return file
}

// InitORM genreate orm handler
func InitORM(db string) (ORM, error) {
	orm, err := xorm.NewEngine("sqlite3", db)
	if err != nil {
		fmt.Println(err)
		return ORM{E: orm}, err
	}
	orm.ShowSQL(true)

	err = orm.CreateTables(&Item{})
	if err != nil {
		fmt.Println(err)
		return ORM{E: orm}, err
	}

	err = orm.CreateTables(&Tag{})
	if err != nil {
		fmt.Println(err)
		return ORM{E: orm}, err
	}

	return ORM{E: orm}, nil
}

func (o *ORM) checkItem(input *widget.Entry) {
	// println("check item:" + input.Text)
	items := []Item{}
	err := o.E.Sql("select * from item where item.name = '" + input.Text + "'").Find(&items)
	if err != nil {
		pe(err)
		return
	}
	if len(items) > 0 {
		ids := items[0].TagID
		setSQL := strings.Join(ids, "','")
		println(setSQL)
		checkText.SetText("Already added")
		tags := []Tag{}
		o.E.Sql("select * from tag where tag.i_d in ('" + setSQL + "')").Find(&tags)
		if len(tags) > 0 {
			str := "Already added Tags: "
			for _, v := range tags {
				str += v.Title + ","
			}
			checkText.SetText(str)
		}
	} else {
		checkText.SetText("Item is NOT exist")
	}
	tagText.SetText("")
}

func (o *ORM) submitItem(input *widget.Entry) {
	// println("submit item：" + input.Text)
	strs := strings.Split(tagText.Text, ",")
	ids := []string{}
	if len(strs) < 1 {
		return
	}
	for _, s := range strs {
		var id uuid.UUID
		var tag Tag
		s = strings.ReplaceAll(s, "[", "")
		s = strings.ReplaceAll(s, "]", "")
		println(s)
		has, _ := o.E.Where("title = ?", s).Get(&tag)
		println(id.String())
		if has {
			println(tag.ID)
			ids = append(ids, tag.ID)
		}
	}

	items := []Item{}
	o.E.Sql("select * from item where item.name = '" + input.Text + "'").Find(&items)
	if len(items) > 0 {
		aff, err := o.E.Exec(`UPDATE item SET tag_i_d = '["` + strings.Join(ids, `","`) + `"]' WHERE i_d = '` + items[0].ID + "' ")
		if err != nil {
			pe(err)
		}
		if num, _ := aff.RowsAffected(); num > 0 {
			checkText.SetText("Update success")
		} else {
			checkText.SetText("Update failed")
		}
	} else {
		num, _ := o.E.Insert(&Item{ID: uuid.New().String(), Name: input.Text, TagID: ids})
		if num > 0 {
			checkText.SetText("Insert success")
		} else {
			checkText.SetText("Insert failed")
		}
	}
	tagText.SetText("")
}

func (o *ORM) addTag(input *widget.Entry) (bool, error) {
	println("add tag：" + input.Text)
	has, err := o.E.Exist(&Tag{
		Title: input.Text,
	})
	if err != nil {
		pe(err)
		return false, err
	}
	if !has {
		_, err = o.E.Insert(&Tag{ID: uuid.New().String(), Title: input.Text})
		if err != nil {
			pe(err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (o *ORM) getTags() []string {
	s := []string{}
	tags := []Tag{}
	err := o.E.Sql("select * from tag").Find(&tags)
	if err != nil {
		pe(err)
		return s
	}
	for _, v := range tags {
		s = append(s, v.Title)
	}
	return s
}
