package main

import (
	"github.com/valyala/fasthttp"
	"log"
)

const serverAddr = ":8081"

func main() {
	server := &fasthttp.Server{}

	server.Handler = func(ctx *fasthttp.RequestCtx) {
		log.Printf("%s %s\n", ctx.Method(), ctx.Request.URI().RequestURI())

		ctx.Response.SetBody([]byte("Hello!\n"))
		ctx.Response.SetStatusCode(200)
	}

	log.Println("Running on", serverAddr)
	log.Fatalln(fasthttp.ListenAndServe(serverAddr, server.Handler))
}
