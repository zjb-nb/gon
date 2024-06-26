package main

import (
	"gon/gonweb"
	"gon/middleware"
	"log"
	"net/http"
	"time"
)

func main() {
	respm := middleware.NewRespMiddleware()
	g := gonweb.NewGracefulshutdown()
	respm.Addpage(http.StatusNotFound, middleware.NotFoundPage())
	s := gonweb.MakeWebServer("web", ":8080", nil,
		g.RejectFilterBuilder,
		middleware.NewTimeMiddleware().ComputerTimeSpend,
		respm.SendRespBuilder,
		respm.HookBuilder,
		middleware.SayHelloFilterBuilder)
	{
		s.Route("GET", "/", func(ctx *gonweb.GonContext) { ctx.JsonOk("hi") })
		s.Route("GET", "/home", func(ctx *gonweb.GonContext) {
			time.Sleep(time.Second * 5)
			ctx.Ok("home")
		})
		s.Route("GET", "/home/user", func(ctx *gonweb.GonContext) { ctx.Ok("home/user") })
		s.Route("GET", "/login", func(ctx *gonweb.GonContext) { ctx.Ok("login") })
		s.Route("GET", "/*", func(ctx *gonweb.GonContext) { ctx.Ok("/* wildcard") })
		s.Route("GET", "/home/*", func(ctx *gonweb.GonContext) { ctx.Ok("/home/* wildcard") })
		s.Route("GET", "/home/user/:id", func(ctx *gonweb.GonContext) { ctx.Ok("/home/user/:id ") })
		s.Route("GET", "/home/~^ab$", func(ctx *gonweb.GonContext) { ctx.Ok("/home/~^ab$") })

		s.Route("POST", "/home/user/index", func(ctx *gonweb.GonContext) { ctx.JsonOk("post-hi") })
	}

	go func() {
		log.Fatal(s.Start())
	}()

	gonweb.WaitShutdown(time.Minute, time.Second*10, g.RejectHookBuilder(), gonweb.TestHookBuilder(s))
}
