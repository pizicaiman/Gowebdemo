package main

import (
	"app"
	"fmt"
	"log"
)

func main() {
	//调用 app包里面的启动
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("go web start ..!")
}
