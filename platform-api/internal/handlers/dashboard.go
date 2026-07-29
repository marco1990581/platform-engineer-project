package handlers

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web/*
var dashboardAssets embed.FS

func DashboardHandler() http.Handler {
	assets, err := fs.Sub(dashboardAssets, "web")
	if err != nil {
		panic(err)
	}

	return http.FileServer(http.FS(assets))
}
