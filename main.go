package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	listenPort := os.Getenv("SIMPLE_LISTEN_PORT")
	targetHost := os.Getenv("SIMPLE_TARGET_HOST")
	targetPort := os.Getenv("SIMPLE_TARGET_PORT")

	fmt.Printf("config: listening on: %v calling: %v:%v\n", listenPort, targetPort, targetHost)

	handler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("received request %v %v %v %v\n", req.Method, req.Host, req.URL, req.URL.User)
		name := req.URL.Query().Get("name")

		if name == "forward" {
			fmt.Printf("==> calling target")
			response, err := http.Get("http://" + targetHost + ":" + targetPort + "/?name=donotforward")
			if err != nil {
				fmt.Printf("error occurred: %v\n", err)
			}
			fmt.Printf("response: %v\n", response.StatusCode)
		}

		fmt.Fprintf(w, "Hello %v\n", name)
	}

	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(":"+listenPort, nil); err != nil {
		fmt.Println("server encountered an error:", err)
	}
}
