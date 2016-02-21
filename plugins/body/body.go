package body

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	c "gopkg.in/h2non/gentleman.v0/context"
	p "gopkg.in/h2non/gentleman.v0/plugin"
	"gopkg.in/h2non/gentleman.v0/utils"
	"io"
	"io/ioutil"
	"strings"
)

// String defines the HTTP request body based on the given string.
func String(data string) p.Plugin {
	return p.NewRequestPlugin(func(ctx *c.Context, h c.Handler) {
		ctx.Request.Body = utils.StringReader(data)
		ctx.Request.ContentLength = int64(bytes.NewBufferString(data).Len())
		h.Next(ctx)
	})
}

// JSON defines a JSON body in the outgoing request.
// Supports strings, array of bytes or buffer.
func JSON(data interface{}) p.Plugin {
	return p.NewRequestPlugin(func(ctx *c.Context, h c.Handler) {
		buf := &bytes.Buffer{}

		switch data.(type) {
		case string:
			buf.WriteString(data.(string))
		case []byte:
			buf.Write(data.([]byte))
		default:
			if err := json.NewEncoder(buf).Encode(data); err != nil {
				h.Error(ctx, err)
				return
			}
		}

		ctx.Request.Body = ioutil.NopCloser(buf)
		ctx.Request.ContentLength = int64(buf.Len())
		ctx.Request.Header.Set("Content-Type", "application/json")

		h.Next(ctx)
	})
}

// XML defines a XML body in the outgoing request.
// Supports strings, array of bytes or buffer.
func XML(data interface{}) p.Plugin {
	return p.NewRequestPlugin(func(ctx *c.Context, h c.Handler) {
		buf := &bytes.Buffer{}

		switch data.(type) {
		case string:
			buf.WriteString(data.(string))
		case []byte:
			buf.Write(data.([]byte))
		default:
			if err := xml.NewEncoder(buf).Encode(data); err != nil {
				h.Error(ctx, err)
				return
			}
		}

		ctx.Request.Body = ioutil.NopCloser(buf)
		ctx.Request.ContentLength = int64(buf.Len())
		ctx.Request.Header.Set("Content-Type", "application/xml")

		h.Next(ctx)
	})
}

// Reader defines a io.Reader stream as request body.
// Content-Type header won't be defined automatically, you have to declare it manually.
func Reader(body io.Reader) p.Plugin {
	return p.NewRequestPlugin(func(ctx *c.Context, h c.Handler) {
		rc, ok := body.(io.ReadCloser)
		if !ok && body != nil {
			rc = ioutil.NopCloser(body)
		}

		req := ctx.Request
		if body != nil {
			switch v := body.(type) {
			case *bytes.Buffer:
				req.ContentLength = int64(v.Len())
			case *bytes.Reader:
				req.ContentLength = int64(v.Len())
			case *strings.Reader:
				req.ContentLength = int64(v.Len())
			}
		}

		req.Body = rc
		h.Next(ctx)
	})
}