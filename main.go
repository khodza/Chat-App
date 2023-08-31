package main

import (
	"context"
	"fmt"
	"net/http"
)

func main() {
	setupAPI()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

func setupAPI() {

	ctx := context.Background()

	manager := NewManager(ctx)

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.serveWS)
	http.HandleFunc("/login", manager.loginHandler)
}
