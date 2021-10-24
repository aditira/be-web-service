package helper

import (
	"fmt"
	"net/http"
)

// Function for handling messages
func PrintMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

// Function for handling errors (upper-case letters is Exported)
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckMethod(w http.ResponseWriter, r *http.Request, method string) {
	if r.Method != method {
		http.Error(w, "Unsupported http method", http.StatusBadRequest)
		return
	}
}
