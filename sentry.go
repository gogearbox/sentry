package sentry

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gogearbox/gearbox"
	"github.com/valyala/fasthttp"
)

type handler struct {
	repanic         bool
	waitForDelivery bool
	timeout         time.Duration
}

// Options struct holds sentry middleware settings
type Options struct {
	// Repanic configures whether Sentry should repanic after recovery
	Repanic bool
	// WaitForDelivery configures whether you want to block the request before moving forward with the response
	WaitForDelivery bool
	// Timeout for the event delivery requests
	Timeout time.Duration
}

// New returns middleware handler
func New(options ...Options) func(ctx gearbox.Context) {
	var op Options
	if len(options) > 0 {
		op = options[0]
	}

	timeout := op.Timeout
	if timeout == 0 {
		timeout = 2 * time.Second
	}
	handler := &handler{
		repanic:         op.Repanic,
		timeout:         timeout,
		waitForDelivery: op.WaitForDelivery,
	}

	return handler.handle
}

func (h *handler) handle(ctx gearbox.Context) {
	hub := sentry.CurrentHub().Clone()
	scope := hub.Scope()
	scope.SetRequest(convert(ctx.Context()))
	scope.SetRequestBody(ctx.Context().Request.Body())

	defer h.recoverWithSentry(hub, ctx)

	ctx.Next()
}

func (h *handler) recoverWithSentry(hub *sentry.Hub, ctx gearbox.Context) {
	if err := recover(); err != nil {
		eventID := hub.RecoverWithContext(
			context.WithValue(context.Background(), sentry.RequestContextKey, ctx),
			err,
		)
		if eventID != nil && h.waitForDelivery {
			hub.Flush(h.timeout)
		}
		if h.repanic {
			panic(err)
		}
	}
}

func convert(ctx *fasthttp.RequestCtx) *http.Request {
	r := new(http.Request)

	r.Method = string(ctx.Method())
	uri := ctx.URI()
	// Ignore error.
	r.URL, _ = url.Parse(fmt.Sprintf("%s://%s%s", uri.Scheme(), uri.Host(), uri.Path()))

	// Headers
	r.Header = make(http.Header)
	r.Header.Add("Host", string(ctx.Host()))
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		r.Header.Add(string(key), string(value))
	})
	r.Host = string(ctx.Host())

	// Cookies
	ctx.Request.Header.VisitAllCookie(func(key, value []byte) {
		r.AddCookie(&http.Cookie{Name: string(key), Value: string(value)})
	})

	// Env
	r.RemoteAddr = ctx.RemoteAddr().String()

	// QueryString
	r.URL.RawQuery = string(ctx.URI().QueryString())

	// Body
	r.Body = ioutil.NopCloser(bytes.NewReader(ctx.Request.Body()))

	return r
}
