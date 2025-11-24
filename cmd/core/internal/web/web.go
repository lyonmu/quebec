package web

import (
	"embed"
	"errors"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

//go:embed dist
var embeddedDist embed.FS

var (
	initOnce          sync.Once
	staticFS          http.FileSystem
	rawFS             fs.FS
	initializationErr error
)

// Register 将打包好的前端资源注册到 Gin，并在未命中任何后端路由时回退到前端路由。
// mountPath 用于控制前端挂载路径，空字符串时默认挂在根路径。
func Register(engine *gin.Engine, mountPath string) error {
	enginePointer := engine
	if enginePointer == nil {
		return errors.New("gin engine is nil")
	}

	if err := prepareFS(); err != nil {
		return err
	}

	prefix := normalizeMountPath(mountPath)
	handler, err := newSPAHandler(prefix, staticFS, rawFS)
	if err != nil {
		return err
	}

	enginePointer.NoRoute(func(ctx *gin.Context) {
		if ctx.Request == nil || ctx.Request.Method != http.MethodGet {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		if !handler.shouldHandle(ctx.Request.URL.Path) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		handler.ServeHTTP(ctx.Writer, ctx.Request)
		ctx.Abort()
	})

	return nil
}

func prepareFS() error {
	initOnce.Do(func() {
		var err error
		rawFS, err = fs.Sub(embeddedDist, "dist")
		if err != nil {
			initializationErr = err
			return
		}
		staticFS = http.FS(rawFS)
	})
	return initializationErr
}

func normalizeMountPath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return "/"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	if len(p) > 1 && strings.HasSuffix(p, "/") {
		p = strings.TrimSuffix(p, "/")
	}
	return path.Clean(p)
}

type spaHandler struct {
	prefix     string
	files      http.FileSystem
	indexBytes []byte
}

func newSPAHandler(prefix string, fsys http.FileSystem, raw fs.FS) (*spaHandler, error) {
	indexHTML, err := fs.ReadFile(raw, "index.html")
	if err != nil {
		return nil, err
	}

	return &spaHandler{
		prefix:     prefix,
		files:      fsys,
		indexBytes: indexHTML,
	}, nil
}

func (s *spaHandler) shouldHandle(requestPath string) bool {
	if s.prefix == "/" {
		return true
	}
	if requestPath == s.prefix {
		return true
	}
	return strings.HasPrefix(requestPath, s.prefix+"/")
}

func (s *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	relPath := s.relativePath(r.URL.Path)
	if relPath == "" {
		s.serveIndex(w, r)
		return
	}

	file, err := s.files.Open(relPath)
	if err != nil {
		s.serveIndex(w, r)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil || info.IsDir() {
		s.serveIndex(w, r)
		return
	}

	http.ServeContent(w, r, info.Name(), info.ModTime(), file)
}

func (s *spaHandler) relativePath(requestPath string) string {
	if s.prefix != "/" {
		requestPath = strings.TrimPrefix(requestPath, s.prefix)
	}
	requestPath = strings.TrimPrefix(requestPath, "/")
	clean := path.Clean(requestPath)
	if clean == "." {
		return ""
	}
	return clean
}

func (s *spaHandler) serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}
	_, _ = w.Write(s.indexBytes)
}
