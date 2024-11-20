package main

import (
	"errors"
	"fmt"
	t "github.com/ManojaD2004/types"
	"github.com/redis/go-redis/v9"
	"gofr.dev/pkg/gofr"
	"os"
	"time"
)

type Customer struct {
	ID   int    `json:"id"`
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
		p := t.New1{}
		c.Bind(&p)
		// Wow many functionality
		fmt.Println(p, "The changes body", os.Getenv("SOME_ENV"))
		// Another way of sending response
		// return map[string]string{"hello": "Hello World"}, nil
		err := errors.New("some error")
		return t.New1{Name: "Hello World"}, err
	})
	app.POST("/greet2", func(c *gofr.Context) (interface{}, error) {
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
	app.GET("/redis", func(ctx *gofr.Context) (interface{}, error) {
		// Get the value using the Redis instance
		redKey := ctx.Param("key")
		val, err := ctx.Redis.Get(ctx.Context, redKey).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			// If the key is not found, we are not considering this an error and returning ""
			return nil, err
		}
		fmt.Println(val)
		return val, nil
	})
	app.POST("/redis", func(ctx *gofr.Context) (interface{}, error) {
		r := t.Redis1{}
		ctx.Bind(&r)
		val, err := ctx.Redis.Set(ctx.Context, r.Key, r.Val, time.Duration(0)).Result()
		fmt.Println(val, err)
		if err != nil && !errors.Is(err, redis.Nil) {
			return nil, err
		}
		return val, nil
	})

	app.POST("/customer/{name}", func(ctx *gofr.Context) (interface{}, error) {
		name := ctx.PathParam("name")

		// Inserting a customer row in database using SQL
		val, err := ctx.SQL.ExecContext(ctx, "INSERT INTO customers (name) VALUES (?)", name)

		return val, err
	})

	app.GET("/customer", func(ctx *gofr.Context) (interface{}, error) {
		var customers []Customer

		// Getting the customer from the database using SQL
		rows, err := ctx.SQL.QueryContext(ctx, "SELECT * FROM customers")
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var customer Customer
			if err := rows.Scan(&customer.ID, &customer.Name); err != nil {
				return nil, err
			}

			customers = append(customers, customer)
		}

		// return the customer
		return customers, nil
	})

	// it can be over-ridden through configs
	app.Run()
}
