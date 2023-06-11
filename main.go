package main

import (
	"github.com/DwGoing/smart_market/internal/app"
	"github.com/alibaba/ioc-golang"
)

func main() {
	if err := ioc.Load(); err != nil {
		panic(err)
	}
	app, err := app.GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.Run()
}
