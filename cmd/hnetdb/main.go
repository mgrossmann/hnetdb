package main

import (
	"net/http"
)

func main() {

	server := http.NewServeMux()
	//server.HandleFunc("/users", users.Register())

	if err := http.ListenAndServe(":8080", server); err != nil {
		panic(err)
	}





}
