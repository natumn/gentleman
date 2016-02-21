package cookies

import (
	"github.com/nbio/st"
	"gopkg.in/h2non/gentleman.v0/context"
	"net/http"
	"testing"
)

func TestCookieSet(t *testing.T) {
	ctx := context.New()
	fn := newHandler()

	Set("foo", "bar").Request(ctx, fn.fn)
	st.Expect(t, fn.called, true)
	st.Expect(t, ctx.Request.Header.Get("Cookie"), "foo=bar")
}

func TestCookieAdd(t *testing.T) {
	ctx := context.New()
	ctx.Request.Header.Set("foo", "foo")
	fn := newHandler()

	Add(&http.Cookie{Name: "foo", Value: "bar"}).Request(ctx, fn.fn)
	st.Expect(t, fn.called, true)
	st.Expect(t, ctx.Request.Header.Get("Cookie"), "foo=bar")
}

func TestCookieDelAll(t *testing.T) {
	ctx := context.New()
	ctx.Request.Header.Set("Cookie", "foo=foo")
	fn := newHandler()

	DelAll().Request(ctx, fn.fn)
	st.Expect(t, fn.called, true)
	st.Expect(t, ctx.Request.Header.Get("Cookie"), "")
}

func TestCookieSetMap(t *testing.T) {
	ctx := context.New()
	fn := newHandler()
	cookies := map[string]string{"foo": "bar"}
	SetMap(cookies).Request(ctx, fn.fn)
	st.Expect(t, fn.called, true)
	st.Expect(t, ctx.Request.Header.Get("Cookie"), "foo=bar")
}

func TestCookieJar(t *testing.T) {
	ctx := context.New()
	fn := newHandler()
	jar := ctx.Client.Jar
	Jar().Request(ctx, fn.fn)
	st.Expect(t, fn.called, true)
	st.Reject(t, ctx.Client.Jar, jar)
}

type handler struct {
	fn     context.Handler
	called bool
}

func newHandler() *handler {
	h := &handler{}
	h.fn = context.NewHandler(func(c *context.Context) {
		h.called = true
	})
	return h
}