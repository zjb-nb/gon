package file

import (
	"gon/gonweb"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func getFileExt(file string) string {
	pos := strings.Index(file, ".")
	return file[pos+1:]
}

// /home/static/a.png ->  /file/static.png
type StaticFile struct {
	dir    string
	prefix string
	ext    map[string]string
}

func (s StaticFile) ServerStaticResource(c *gonweb.GonContext) {
	path := strings.TrimPrefix(c.R.URL.Path, s.prefix)
	path = filepath.Join(s.dir, path)
	f, err := os.Open(path)
	if err != nil {
		c.ResponStatus = http.StatusNotFound
		c.ResponsData = []byte(err.Error())
		return
	}
	defer f.Close()
	ext := getFileExt(f.Name())
	ext, ok := s.ext[ext]
	if !ok {
		c.ResponStatus = http.StatusNotFound
		c.ResponsData = []byte(err.Error())
		return
	}
	data, err := io.ReadAll(f)
	if err != nil {
		c.ResponStatus = http.StatusNotFound
		c.ResponsData = []byte(err.Error())
		return
	}
	c.ResponsData = data
	c.ResponStatus = http.StatusOK
}

func NewStaticFile(dir, prefix string, ext map[string]string) *StaticFile {
	return &StaticFile{
		dir:    dir,
		prefix: prefix,
		ext:    ext,
	}
}
func NewStaticIMGFile(dir, prefix string) *StaticFile {
	return &StaticFile{
		dir:    dir,
		prefix: prefix,
		ext: map[string]string{
			"jpeg": "image/jpeg",
			"jpg":  "image/jpeg",
			"jpe":  "image/jpeg",
			"pdf":  "image/pdf",
			"png":  "image/png",
		},
	}
}
