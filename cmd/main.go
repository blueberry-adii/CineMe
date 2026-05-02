package main

import "net/http"

func main() {
	dir := http.Dir("./static")
	if err := http.ListenAndServe(":80", http.FileServer(dir)); err != nil {
		panic(err)
	}
}
