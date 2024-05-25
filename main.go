package main

import (
	"gon/gonweb"
	"log"
)

func main() {
	s := gonweb.MakeWebServer("web", ":8080", nil)
	s.Route("GET", "/", func(ctx *gonweb.GonContext) { ctx.JsonOk("hi") })
	s.Route("GET", "/home", func(ctx *gonweb.GonContext) { ctx.Ok("home") })
	s.Route("GET", "/home/user", func(ctx *gonweb.GonContext) { ctx.Ok("home/user") })
	s.Route("GET", "/login", func(ctx *gonweb.GonContext) { ctx.Ok("login") })
	s.Route("GET", "/*", func(ctx *gonweb.GonContext) { ctx.Ok("/* wildcard") })
	s.Route("GET", "/home/*", func(ctx *gonweb.GonContext) { ctx.Ok("/home/* wildcard") })
	s.Route("GET", "/home/user/:id", func(ctx *gonweb.GonContext) { ctx.Ok("/home/user/:id ") })
	s.Route("GET", "/home/~^ab$", func(ctx *gonweb.GonContext) { ctx.Ok("/home/~^ab$") })

	s.Route("POST", "/home/user/index", func(ctx *gonweb.GonContext) { ctx.JsonOk("post-hi") })

	log.Fatal(s.Start())
}
