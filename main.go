package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ManojaD2004/route"
	t "github.com/ManojaD2004/types"
	"github.com/redis/go-redis/v9"
	"gofr.dev/examples/using-add-rest-handlers/migrations"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"gofr.dev/pkg/gofr/websocket"
)

type Customer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type user struct {
	ID         int    `json:"id"`
	Name       string `json:"name"  sql:"not_null"`
	Age        int    `json:"age"`
	IsEmployed bool   `json:"isEmployed"`
}

// GetAll : User can overwrite the specific handlers by implementing them like this
// func (u *user) GetAll(c *gofr.Context) (interface{}, error) {
// 	return "user GetAll called", nil
// }

func main() {
	// initialise gofr object
	app := gofr.New()
	// register route greet
	app.Metrics().NewCounter("transaction_success", "used to track the count of successful transactions")
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
		entries, _ := c.File.ReadDir("./")

		for _, entry := range entries {
			entryType := "File"

			if entry.IsDir() {
				entryType = "Dir"
			}

			fmt.Printf("%v: %v Size: %v Last Modified Time : %v\n", entryType, entry.Name(), entry.Size(), entry.ModTime())
		}
		c.Metrics().IncrementCounter(c, "transaction_success")
		span := c.Trace("my-custom-span")
		defer span.End()
		return t.New1{Name: "Hello World"}, nil
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
	app.POST("/user", route.UserGetRoute)

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

	app.POST("/test-api", TestApiService)
	// Add migrations to run
	app.Migrate(migrations.All())

	// AddRESTHandlers creates CRUD handles for the given entity
	err := app.AddRESTHandlers(&user{})
	if err != nil {
		return
	}

	// Cron Jobs
	// app.AddCronJob("*/30 * * * * *", "30 second job", func(ctx *gofr.Context) {
	// 	fmt.Println("current time is", time.Now())
	// })
	// app.EnableBasicAuth("username" ,"password")
	app.UseMiddlewareWithContainer(customMiddleware)
	app.AddHTTPService("test-api-service", "https://jsonplaceholder.typicode.com")
	wsUpgrader := websocket.NewWSUpgrader(
		websocket.WithHandshakeTimeout(5*time.Second), // Set handshake timeout
		websocket.WithReadBufferSize(2048),            // Set read buffer size
		websocket.WithWriteBufferSize(2048),           // Set write buffer size
		websocket.WithSubprotocols("chat", "binary"),  // Specify subprotocols
		websocket.WithCompression(),                   // Enable compression
	)

	app.OverrideWebsocketUpgrader(wsUpgrader)
	app.WebSocket("/ws", WSHandler)
	// it can be over-ridden through configs
	app.Run()
}

func WSHandler(ctx *gofr.Context) (interface{}, error) {
	var message string

	err := ctx.Bind(&message)
	if err != nil {
		ctx.Logger.Errorf("Error binding message: %v", err)
		return nil, err
	}

	ctx.Logger.Infof("Received message: %s", message)

	err = ctx.WriteMessageToSocket("Hello! GoFr")
	if err != nil {
		return nil, err
	}

	return message, nil
}

func customMiddleware(c *container.Container, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Logger.Log("Hey! Welcome to GoFr")
		handler.ServeHTTP(w, r)
	})
}

type TestApiType struct {
	UserId    int    `json:"userId"`
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func TestApiService(ctx *gofr.Context) (interface{}, error) {
	fmt.Println("I am In test-api")
	testApi := ctx.GetHTTPService("test-api-service")
	// Use the Get method to call the GET /user endpoint of payments service
	resp, err := testApi.Get(ctx, "todos/1", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ret := TestApiType{}
	json.Unmarshal(body, &ret)
	return ret, nil
}
