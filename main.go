package main

import (
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
	manager := NewManager()

	http.Handle("/", http.FileServer(http.Dir("./fronted")))
	http.HandleFunc("/ws", manager.serveWS)
}
