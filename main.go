package main

import (
	"gofr.dev/pkg/gofr"
	"errors"
	"fmt"
	"os"
	t "github.com/ManojaD2004/types"
)

func main() {
	// initialise gofr object
	app := gofr.New()
	// register route greet
	app.GET("/greet", func(ctx *gofr.Context) (interface{}, error) {
		return "Hello World!", nil
	})
    app.POST("/greet1", func(c *gofr.Context) (interface{}, error) {
        p := t.New1{}
        c.Bind(&p)
		// Wow many functionality
        fmt.Println(p, "The changes body", os.Getenv("SOME_ENV"))
		// Another way of sending response
        // return map[string]string{"hello": "Hello World"}, nil
		err := errors.New("some error")
		return t.New1{Name: "Hello World"}, err
    })
	app.GET("/greet1", func(ctx *gofr.Context) (interface{}, error) {
		return t.New1{Name: "Hello World"}, nil
	})

	// Runs the server, it will listen on the default port 8000.
	// it can be over-ridden through configs
	app.Run()
}
