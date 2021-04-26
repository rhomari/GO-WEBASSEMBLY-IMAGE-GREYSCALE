package main

import "net/http"

func main() {
	http.ListenAndServe(":2304", http.FileServer(http.Dir(".")))
}
