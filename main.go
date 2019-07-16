package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	InitRouter()

	fmt.Println("---------Server Start!---------")
	fmt.Println("Port: ", 10001)
	log.Fatal(http.ListenAndServe(":10001", nil))
}
