package access_log

import (
	"encoding/json"
	"github.com/soluble1/mweb"
	"log"
)

type MiddlewareBuilder struct {
	logFunc func(accLog string)
}

func NewAccessLog() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(accLog string) {
			log.Println(accLog)
		},
	}
}

func (m *MiddlewareBuilder) LogFunc(logFunc func(string)) *MiddlewareBuilder {
	m.logFunc = logFunc

	return m
}

type accessLog struct {
	Host       string
	Route      string
	HTTPMethod string
	Path       string
}

func (m *MiddlewareBuilder) Build() mweb.Middleware {
	return func(next mweb.HandlerFunc) mweb.HandlerFunc {
		return func(ctx *mweb.Context) {
			defer func() {
				accLog := &accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchedRoute,
					HTTPMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
				}
				val, _ := json.Marshal(accLog)
				m.logFunc(string(val))
			}()

			next(ctx)
		}
	}
}
