package controllers

import (
	"github.com/newham/hamgo"
)

func SessionFilter(ctx hamgo.Context) bool {
	if ctx.GetSession() == nil || ctx.GetSession().Get("USER") == nil {
		ctx.Redirect("/login")
		return false
	}
	return true
}
