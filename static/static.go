package static

import (
	"embed"
	"net/http"
)

//go:embed css js
var static embed.FS

func Handler() http.Handler {
	return http.FileServer(http.FS(static))
}
