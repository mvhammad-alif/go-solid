package main

import (
	"go-solid/internal/app"
)

func main() {
	e := app.InitHTTPServer()
	e.Logger.Fatal(e.Start(":1323"))
}
