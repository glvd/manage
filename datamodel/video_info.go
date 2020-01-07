package datamodel

import (
	"github.com/glvd/backmanage/model"
	"github.com/glvd/backmanage/modules/scrape"
	"github.com/glvd/go-admin/modules/db"
	form2 "github.com/glvd/go-admin/plugins/admin/modules/form"
	"github.com/glvd/go-admin/plugins/admin/modules/table"
	"github.com/glvd/go-admin/template/types"
	"github.com/glvd/go-admin/template/types/form"
	editType "github.com/glvd/go-admin/template/types/table"
	"github.com/goextension/log"
	"os"
	"path/filepath"
	"strings"
)

// VideoInfoTable ...
func VideoInfoTable() (t table.Table) {
	t = table.NewDefaultTable(table.Config{
		Driver:     db.DriverMysql,
		CanAdd:     true,
		Editable:   false,
		Deletable:  true,
		Exportable: true,
		Connection: table.DefaultConnectionName,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: table.DefaultPrimaryKeyName,
		},
	})
	info := t.GetInfo().HideFilterArea()
	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("Poster", "poster_path", db.Text).FieldDisplay(func(value types.FieldModel) interface{} {
		if value.Value == "" {
			return ""
		}

		img := ImageLoad(value.Value)
		if img != "" {
			img = "data:image/jpg;base64," + img
		}
		return "<img height=\"120px\" src=\"" + img + "\"/>"
	})
	info.AddField("VideoID", "video_id", db.Varchar)
	info.AddField("VideoNo", "video_no", db.Varchar).FieldWidth(120).FieldSortable().FieldEditAble(editType.Text).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Intro", "intro", db.Varchar).FieldWidth(640).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Actors", "actors", db.Varchar).FieldWidth(120).FieldSortable().FieldDisplay(func(value types.FieldModel) interface{} {
		return SplitArguments(value)
	}).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Tags", "tags", db.Varchar).FieldWidth(180).FieldDisplay(func(value types.FieldModel) interface{} {
		return SplitArguments(value)
	}).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Source", "source_path", db.Varchar).FieldEditAble(editType.Text).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("CreateTime", "created_at", db.Timestamp).FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField("UpdateTime", "updated_at", db.Timestamp)

	info.SetTable("dhash_video_infos").SetTitle("Videos").SetDescription("Videos")

	//t.GetInfo().SetTabGroups(types.
	//	NewTabGroups("video_no", "intro", "created_at").
	//	AddGroup("source_path", "tags", "actors")).
	//	SetTabHeaders("profile1", "profile2")
	//edit/add form
	formList := t.GetForm()
	formList.SetBeforeInsert(VideoInsert)

	formList.AddField("VideoNo", "video_no", db.Varchar, form.Text).FieldNotAllowEdit()
	//formList.AddField("PosterPath", "poster_path", db.Varchar, form.Text)
	//formList.AddField("ThumbPath", "thumb_path", db.Varchar, form.Text)
	formList.AddField("SourcePath", "source_path", db.Varchar, form.Text)
	//t.GetInfo().SetTabGroups(types.
	//	NewTabGroups("video_no", "intro", "created_at").
	//	AddGroup("source_path", "tags", "actors"))
	//SetTabHeaders("profile1", "profile2")
	formList.SetTable("dhash_video_infos").SetTitle("Videos").SetDescription("Videos")
	return
}

// VideoInsert ...
func VideoInsert(values form2.Values) error {
	no := strings.ToUpper(strings.TrimSpace(values.Get("video_no")))
	if no != "" {
		err := os.RemoveAll(filepath.Join("data", "info", no))
		if err != nil {
			log.Error("remove", err)
		}
		values.Add("video_no", no)
		c, err := scrape.FindContent(no)
		if err != nil {
			return err
		}
		var v model.VideoDetail
		err = v.CopyInfo(c)
		if err != nil {
			return err
		}

		err = model.InsertVideo(&v)
		if err != nil {
			return err
		}
		log.Infow("insert", "values", values)
		values.Add("poster_path", filepath.Join("data", "info", v.No, "image.jpg"))
		values.Add("thumb_path", filepath.Join("data", "info", v.No, "thumb.jpg"))
		values.Add("info_path", filepath.Join("data", "info", v.No, ".info"))
		values.Add("internal_id", v.InternalID)
		values.Add("intro", v.Intro)
		values.Add("actors", v.RoleString())
		values.Add("tags", v.TagString())
	}
	return nil
}

// SplitArguments ...
func SplitArguments(value types.FieldModel) interface{} {
	ss := strings.Split(value.Value, ",")
	var rlt []string
	for _, s := range ss {
		rlt = append(rlt, "<a href=\"#\">"+s+"</a>")
	}
	return strings.Join(rlt, "<br>")
}