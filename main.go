package main

import (
    "gofr.dev/pkg/gofr" 
    "fmt"
)

type New1 struct {
	Name string `json:"name"`
}

func main() {
	// initialise gofr object
	app := gofr.New()
	// register route greet
	app.GET("/greet", func(ctx *gofr.Context) (interface{}, error) {
		return "Hello World!", nil
	})
    app.POST("/greet1", func(c *gofr.Context) (interface{}, error) {
        p := New1{}
        c.Bind(&p)
        fmt.Println("%v The changes body", p)
        return "Hello World!", nil
    })
	app.GET("/greet1", func(ctx *gofr.Context) (interface{}, error) {
		return New1{"Hello World"}, nil
	})

	// Runs the server, it will listen on the default port 8000.
	// it can be over-ridden through configs
	app.Run()
}
