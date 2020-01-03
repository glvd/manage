package guard

import (
	"github.com/glvd/go-admin/context"
	"github.com/glvd/go-admin/modules/auth"
	"github.com/glvd/go-admin/modules/service"
	"html/template"
	"strconv"
)

// MenuEditParam ...
type MenuEditParam struct {
	Id       string
	Title    string
	Header   string
	ParentId int64
	Icon     string
	Uri      string
	Roles    []string
	Alert    template.HTML
}

// HasAlert ...
func (e MenuEditParam) HasAlert() bool {
	return e.Alert != template.HTML("")
}

// MenuEdit ...
func MenuEdit(srv service.List) context.Handler {
	return func(ctx *context.Context) {

		parentId := ctx.FormValue("parent_id")
		if parentId == "" {
			parentId = "0"
		}

		var (
			parentIdInt, _ = strconv.Atoi(parentId)
			token          = ctx.FormValue("_t")
			alert          template.HTML
		)

		if !auth.GetService(srv.Get("auth")).CheckToken(token) {
			alert = getAlert("edit fail, wrong token")
		}

		if alert == "" {
			alert = checkEmpty(ctx, "id", "title", "icon")
		}

		// TODO: check the user permission

		ctx.SetUserValue("edit_menu_param", &MenuEditParam{
			Id:       ctx.FormValue("id"),
			Title:    ctx.FormValue("title"),
			Header:   ctx.FormValue("header"),
			ParentId: int64(parentIdInt),
			Icon:     ctx.FormValue("icon"),
			Uri:      ctx.FormValue("uri"),
			Roles:    ctx.Request.Form["roles[]"],
			Alert:    alert,
		})
		ctx.Next()
	}
}

// GetMenuEditParam ...
func GetMenuEditParam(ctx *context.Context) *MenuEditParam {
	return ctx.UserValue["edit_menu_param"].(*MenuEditParam)
}

func checkEmpty(ctx *context.Context, key ...string) template.HTML {
	for _, k := range key {
		if ctx.FormValue(k) == "" {
			return getAlert("wrong " + k)
		}
	}
	return template.HTML("")
}
