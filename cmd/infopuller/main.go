package main

import "infopuller/internal/app"

func main() {
	a, err := app.New()
	if err != nil {
		panic(err)
	}
	a.Run()
}
