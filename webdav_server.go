package main

import (
	"fmt"
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
	final_path_string := fmt.Sprintf("%s/%s", cfg.Prefix_url, user.project)
	if !strings.HasPrefix(req.URL.Path, final_path_string) {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	handler := webdav.Handler{final_path_string, fs, webDavLs, func(request *http.Request, err error) {
		checkErr(err)
	}}

	handler.ServeHTTP(w, req)
}
