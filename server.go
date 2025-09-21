package main

import (
	"go-solid/internal/app"
)

func main() {
	e, err := app.InitHTTPServer()
	if (err != nil) {
		panic(err)
	}

	e.Logger.Fatal(e.Start(":1323"))
}
