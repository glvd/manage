package datamodel

import (
	"github.com/glvd/go-admin/modules/db"
	"github.com/glvd/go-admin/plugins/admin/modules/table"
	"github.com/glvd/go-admin/template/types/form"
	editType "github.com/glvd/go-admin/template/types/table"
)

// GetPostsTable return the model of table posts.
func GetPostsTable() (postsTable table.Table) {

	postsTable = table.NewDefaultTable(table.DefaultConfig())

	info := postsTable.GetInfo()
	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("Title", "title", db.Varchar)
	info.AddField("Description", "description", db.Varchar)
	info.AddField("Content", "content", db.Varchar).FieldEditAble(editType.Textarea)
	info.AddField("Date", "data", db.Varchar)

	info.SetTable("posts").SetTitle("Posts").SetDescription("Posts")

	formList := postsTable.GetForm()
	formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowEdit().FieldNotAllowAdd()
	formList.AddField("Title", "title", db.Varchar, form.Text)
	formList.AddField("Description", "description", db.Varchar, form.Text)
	formList.AddField("Content", "content", db.Varchar, form.Text)
	formList.AddField("Date", "data", db.Varchar, form.Datetime)

	formList.SetTable("posts").SetTitle("Posts").SetDescription("Posts")

	return
}
