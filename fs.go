package staticserve

import (
	"io/fs"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type staticFs struct {
	fs      fs.FS        // filesystem
	stripFs string       // mountpoint
	htfs    http.Handler // handler for Filesystem
	stdh    string       // standard redirect target
}

func NewFs(fs fs.FS, strip string) *staticFs {
	htfs := http.FS(fs)

	return &staticFs{
		fs:      fs,
		stripFs: strip,
		htfs:    http.FileServer(htfs),
	}
}

func (s *staticFs) WithStandard(std string) *staticFs {
	s.stdh = std
	return s
}

func (s *staticFs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.CheckFile(r.URL) {
		s.htfs.ServeHTTP(w, r)
		return
	}
	r2 := new(http.Request)
	*r2 = *r
	r.URL.Path = s.stdh
	s.htfs.ServeHTTP(w, r2)
}

func checkFileFromString(fs fs.FS, name string) bool {
	_, err := fs.Open(name)
	return err == nil
}

func extractName(r *url.URL) string {
	return r.Path
}

func (s *staticFs) CheckFile(r *url.URL) bool {
	p := extractName(r)
	if s.stripFs != "" {
		p = filepath.Join(s.stripFs, p)
	}
	p = strings.TrimLeft(p, "/")
	return checkFileFromString(s.fs, p)
}
