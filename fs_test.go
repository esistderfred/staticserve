package staticserve

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func mustparse(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

func Test_extractName(t *testing.T) {
	type args struct {
		r *url.URL
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"test1",
			args{
				mustparse("http://blabla/"),
			},
			"/",
		},
		{
			"test1",
			args{
				mustparse("http://blabla/css/123.css"),
			},
			"/css/123.css",
		},
		{
			"test1",
			args{
				mustparse("http://blabla/index.html"),
			},
			"/index.html",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractName(tt.args.r); got != tt.want {
				t.Errorf("extractName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkFileFromString(t *testing.T) {
	files := os.DirFS("./testdata/fs1")

	type args struct {
		fs   fs.FS
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"test1",
			args{
				fs:   files,
				name: "index.html",
			},
			true,
		},
		{
			"test2",
			args{
				fs:   files,
				name: "index3.html",
			},
			false,
		},
		{
			"test3",
			args{
				fs:   files,
				name: "static/js1.js",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkFileFromString(tt.args.fs, tt.args.name); got != tt.want {
				t.Errorf("checkFileFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_staticFs_CheckFile(t *testing.T) {
	type args struct {
		r *url.URL
	}
	tests := []struct {
		name string
		s    *staticFs
		args args
		want bool
	}{
		{
			"test1",
			&staticFs{
				fs:      os.DirFS("./testdata/fs1"),
				stripFs: "",
			},
			args{
				mustparse("http://bla/index.html"),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.CheckFile(tt.args.r); got != tt.want {
				t.Errorf("staticFs.CheckFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_staticFs_ServeHTTP(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name string
		s    *staticFs
		args args
		want string
	}{
		{
			"test1",
			NewFs(os.DirFS("./testdata/fs1"), ""),
			args{
				httptest.NewRecorder(),
				httptest.NewRequest("GET", "http://test.de/asdf", nil),
			},
			"INDEX",
		},
		{
			"test2",
			NewFs(os.DirFS("./testdata/fs1"), ""),
			args{
				httptest.NewRecorder(),
				httptest.NewRequest("GET", "http://test.de/index.js", nil),
			},
			"JS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.ServeHTTP(tt.args.w, tt.args.r)
			if tt.args.w.Code != http.StatusOK {
				t.Error("return code was ", tt.args.w.Code, "instead of", http.StatusOK)
			}
			str := tt.args.w.Body.String()
			if str != tt.want {
				t.Error(str, "!=", tt.want)
			}
		})
	}
}
