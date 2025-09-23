package goweb

import (
	"fmt"
	"net/http"
)
s
func main() {
	port := ":3333"
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Get Users")
	io.WriteString(w, "Hello, World!")
}

func getCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Get Courses")
	io.WriteString(w, "Hello, World!")
}