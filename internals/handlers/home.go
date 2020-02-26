package handlers

import (
	"fmt"
	"net/http"
)

// HomeHandler function for home router
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You've requested a home route")
}
