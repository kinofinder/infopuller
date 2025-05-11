package main

import "infopuller/internal/utils/config"

func main() {
	err := config.Update()
	if err != nil {
		panic(err)
	}
}
