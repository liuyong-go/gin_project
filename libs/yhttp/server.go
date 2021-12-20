package yhttp

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/liuyong-go/gin_project/config"
	"github.com/liuyong-go/gin_project/libs/logger"
	"github.com/liuyong-go/gin_project/routes"
)

var srv = &http.Server{
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20,
}

//开启http
func Start() {
	c := config.Config
	gin.SetMode(c.HTTP.Mode)
	r := gin.New()
	r.Use(gin.Recovery())
	store := cookie.NewStore([]byte(c.HTTP.CookieSecret))
	store.Options(sessions.Options{
		Domain:   c.HTTP.CookieDomain,
		MaxAge:   c.HTTP.CookieMaxAge,
		Secure:   c.HTTP.CookieSecure,
		HttpOnly: c.HTTP.CookieHttpOnly,
		Path:     "/",
	})
	session := sessions.Sessions(c.HTTP.CookieName, store)
	r.Use(session)
	routes.SetRoutes(r)
	srv.Addr = c.HTTP.Listen
	srv.Handler = r
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			os.Exit(3)
		}
	}()
}

//关闭http
func Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Info(context.TODO(), "cannot shutdown http server:", err)
		os.Exit(2)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		logger.Info(context.TODO(), "shutdown http server timeout of 5 seconds.")
	default:
		logger.Info(context.TODO(), "http server stopped")
	}
}
