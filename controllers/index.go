package controllers

import (
	"net/url"
	"strings"

	"github.com/newham/gofs/api"
	"github.com/newham/hamgo"
)

func FolderController(ctx hamgo.Context) {
	path := getPath(ctx, "/folder/")
	folders := api.GetFolder(path)
	ctx.JSONHTML(folders, "public/index.html")
}

func getPath(ctx hamgo.Context, prefix string) string {
	path, _ := url.QueryUnescape(strings.TrimPrefix(ctx.R().URL.String(), prefix))
	return path
}

func FileController(ctx hamgo.Context) {
	ctx.File(api.ROOT_PATH + getPath(ctx, "/file/"))
}

func IndexController(ctx hamgo.Context) {
	ctx.Redirect("/folder/")
}

func VideoController(ctx hamgo.Context) {
	video := getPath(ctx, "/video/")
	ctx.PutData("video", video)
	ctx.HTML("public/player.html")
}
