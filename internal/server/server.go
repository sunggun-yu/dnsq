package server

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/sunggun-yu/dnsq/internal/server/handlers"
)

type Server struct {
	Port        int
	StaticFiles fs.FS
}

// StripPrefixFS is a custom fs.FS that strips a given prefix from the file paths.
type StripPrefixFS struct {
	fs     fs.FS
	prefix string
}

// override the Open method to strip the prefix
func (s *StripPrefixFS) Open(name string) (fs.File, error) {
	return s.fs.Open(path.Join(s.prefix, name))
}

// NewServer creates a new server with the given port and static files.
func NewServer(port int, staticFiles embed.FS) *Server {

	// Create a new fs.FS that strips the "static" prefix
	strippedFS := &StripPrefixFS{
		fs:     staticFiles,
		prefix: "static",
	}

	return &Server{
		Port:        port,
		StaticFiles: strippedFS,
	}
}

// Run starts the server and listens on the given port with the given static files.
func (s *Server) Run() {

	// Create a new Gin router
	r := gin.Default()

	// Serve static files, such as index.html
	r.StaticFS("/static", http.FS(s.StaticFiles))

	// API endpoint for DNS lookup
	r.GET("/api/lookup", handlers.DNSLookupHandler)

	// Serve the index.html at the root. do not specify any path to serve the index.html. it will gets too many redirects if you specify the path.
	// not sure why lol
	r.NoRoute(func(c *gin.Context) {
		c.FileFromFS("", http.FS(s.StaticFiles))
	})

	// Run the server on the given port
	r.Run(fmt.Sprintf(":%d", s.Port))
}
