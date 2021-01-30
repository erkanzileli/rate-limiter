package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
)

const serverAddr = ":8081"

func main() {
	server := &fasthttp.Server{}

	server.Handler = func(ctx *fasthttp.RequestCtx) {
		fmt.Printf("%s %s\n", ctx.Method(), ctx.Request.URI().RequestURI())

		ctx.Response.SetBody([]byte("Hello!\n"))
		ctx.Response.SetStatusCode(200)
	}

	fmt.Println("Running on", serverAddr)
	log.Fatalln(fasthttp.ListenAndServe(serverAddr, server.Handler))
}
