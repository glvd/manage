package datamodel

import (
	"github.com/glvd/go-admin/modules/db"
	"github.com/glvd/go-admin/plugins/admin/modules/table"
	"github.com/glvd/go-admin/template/types"
	"github.com/glvd/go-admin/template/types/form"
	editType "github.com/glvd/go-admin/template/types/table"
)

// VideoListTable ...
func VideoListTable() (vlTable table.Table) {
	vlTable = table.NewDefaultTable(table.Config{
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
	info := vlTable.GetInfo()
	info.AddField("ID", "id", db.Varchar).FieldSortable()
	info.AddField("path", "video_id", db.Varchar)
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

	info.SetTable("videos").SetTitle("Videos").SetDescription("Videos")

	//vlTable.GetInfo().SetTabGroups(types.
	//	NewTabGroups("video_no", "intro", "created_at").
	//	AddGroup("source_path", "tags", "actors")).
	//	SetTabHeaders("profile1", "profile2")
	//edit/add form
	formList := vlTable.GetForm()
	formList.SetBeforeInsert(VideoInsert)

	formList.AddField("VideoNo", "video_no", db.Varchar, form.Text).FieldNotAllowEdit()
	//formList.AddField("PosterPath", "poster_path", db.Varchar, form.Text)
	//formList.AddField("ThumbPath", "thumb_path", db.Varchar, form.Text)
	formList.AddField("SourcePath", "source_path", db.Varchar, form.Text)
	//vlTable.GetInfo().SetTabGroups(types.
	//	NewTabGroups("video_no", "intro", "created_at").
	//	AddGroup("source_path", "tags", "actors"))
	//SetTabHeaders("profile1", "profile2")
	formList.SetTable("videos").SetTitle("Videos").SetDescription("Videos")
	return
}