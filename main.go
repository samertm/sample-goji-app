package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func serveIndex(c web.C, w http.ResponseWriter, r *http.Request) error {
	s := getSession(c)
	u := getUser(s)
	v := indexTemplateVars{}
	if u != nil {
		v.Name = u.Name
	}
	return indexTemplate.Execute(w, v)
}

func serveAuthCallback(c web.C, w http.ResponseWriter, r *http.Request) error {
	s := getSession(c)
	var name string
	if q := r.URL.Query(); len(q["name"]) > 0 {
		name = q["name"][0]
	}
	u, err := GetCreateUser(name)
	if err != nil {
		return err
	}
	s.Values[userIDSessionKey] = u.ID
	if err := s.Save(r, w); err != nil {
		log.Println(err)
	}
	return HTTPRedirect{To: "/", Code: http.StatusTemporaryRedirect}
}

func main() {
	// Serve static files.
	staticDirs := []string{"bower_components", "res"}
	for _, d := range staticDirs {
		static := web.New()
		pattern, prefix := fmt.Sprintf("/%s/*", d), fmt.Sprintf("/%s/", d)
		static.Get(pattern, http.StripPrefix(prefix, http.FileServer(http.Dir(d))))
		http.Handle(prefix, static)
	}

	goji.Use(applySessions)
	goji.Use(context.ClearHandler)

	goji.Get("/", handler(serveIndex))
	goji.Get("/auth_callback", handler(serveAuthCallback))
	goji.Serve()
}
