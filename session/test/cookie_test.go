package test

import (
	"github.com/google/uuid"
	"github.com/soluble1/mweb"
	"github.com/soluble1/mweb/session"
	"github.com/soluble1/mweb/session/cookie"
	"github.com/soluble1/mweb/session/memory"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestCookie(t *testing.T) {
	server := mweb.NewHTTPServer()

	manager := session.Manager{
		SessCtxKey: "my_session",
		Store:      memory.NewStore(30 * time.Minute),
		Propagator: cookie.NewPropagator("sessId", func(p *cookie.Propagator) {
			cookie.WithCookieOpt(func(cookie *http.Cookie) {
				// cookie.HttpOnly = true, js脚本将无法读取到cookie信息
				cookie.HttpOnly = true
			})
		}),
	}

	server.Post("/login", func(ctx *mweb.Context) {
		ctx.RespData = []byte("this is /login")

		id := uuid.New()
		sess, err := manager.InitSession(ctx, id.String())
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}

		err = sess.Set(ctx.Req.Context(), "myKey", "xiaoLongRen")
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
	})

	server.Get("/resource", func(ctx *mweb.Context) {
		sess, err := manager.GetSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}

		val, err := sess.Get(ctx.Req.Context(), "myKey")
		ctx.RespData = []byte(val)
		ctx.RespStatusCode = 200
	})

	server.Post("/logout", func(ctx *mweb.Context) {
		_ = manager.RemoveSession(ctx)
	})

	server.Use(func(next mweb.HandlerFunc) mweb.HandlerFunc {
		return func(ctx *mweb.Context) {
			if ctx.Req.URL.Path != "/login" {
				sess, err := manager.GetSession(ctx)
				if err != nil {
					ctx.RespStatusCode = 401
					log.Println(err)
					return
				}
				_ = manager.Refresh(ctx.Req.Context(), sess.ID())
			}
			next(ctx)
		}
	})

	server.Start(":8081")
}
