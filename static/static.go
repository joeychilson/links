package static

import (
	"embed"
	"net/http"
)

//go:embed css
var static embed.FS

func Handler() http.Handler {
	return http.FileServer(http.FS(static))
}
