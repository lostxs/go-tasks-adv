package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {
		number := rand.IntN(6) + 1
		w.Write([]byte(strconv.Itoa(number)))
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server listening op port 8081")
	server.ListenAndServe()
}
