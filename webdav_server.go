package main

import (
	"golang.org/x/net/webdav"
	"net/http"
	"path"
	"strings"
)

var webDavLs webdav.LockSystem

func initWebDav() {
	webDavLs = webdav.NewMemLS()
}

func webdavHandler(w http.ResponseWriter, req *http.Request, user User) {
	fs := webdav.Dir(path.Join(cfg.Root_dir, "projects", user.project))
	if !strings.HasPrefix(req.URL.Path, cfg.Prefix_url+user.project) {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	handler := webdav.Handler{cfg.Prefix_url + user.project, fs, webDavLs, func(request *http.Request, err error) {
		checkErr(err)
	}}

	handler.ServeHTTP(w, req)
}
