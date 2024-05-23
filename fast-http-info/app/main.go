// Package helloworld provides a set of Cloud Functions samples.
package hello

import (
	"cloud.google.com/go/functions/metadata"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
)

func init() {
	functions.HTTP("entryPoint", entryPoint)
}

func entryPoint(w http.ResponseWriter, r *http.Request) {
	m, _ := metadata.FromContext(r.Context())
	ctx := &fasthttp.RequestCtx{}
	ctx.SetUserValue("metadata", m)
	ctx.Request.Header.SetMethod(r.Method)
	ctx.Request.SetRequestURI(r.URL.String())
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// Handle the error appropriately
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	ctx.Request.SetBody(body)
	for key, values := range r.Header {
		for _, value := range values {
			ctx.Request.Header.Add(key, value)
		}
	}
	handleFastHTTPRequest(ctx)
	w.WriteHeader(ctx.Response.StatusCode())
	w.Write(ctx.Response.Body())
}

func handleFastHTTPRequest(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world!\n\n")

	fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
	fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
	fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
	fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
	fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
	fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
	fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
	fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
	fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
	fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())

	fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)

	ctx.SetContentType("text/plain; charset=utf8")

	// Set arbitrary headers
	ctx.Response.Header.Set("X-My-Header", "my-header-value")

	// Set cookies
	var c fasthttp.Cookie
	c.SetKey("cookie-name")
	c.SetValue("cookie-value")
	ctx.Response.Header.SetCookie(&c)
}
