package www

import (
	"embed"
	"io/fs"
	"os"
)

//go:embed *.htm *.js *.css
var content embed.FS
var Content fs.FS

func init() {
	if v, _ := os.LookupEnv("USE_LOCAL"); v == "1" {
		Content = os.DirFS("internal/server/www")
	} else {
		Content = content
	}
}
