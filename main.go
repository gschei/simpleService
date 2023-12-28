package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
)

func initSpiffeServer(ctx context.Context, port string) {

	// export SPIFFE_ENDPOINT_SOCKET=unix:///run/spire/sockets/agent.sock

	clientID := spiffeid.RequireFromString("spiffe://example.org/gilbert/testclient")

	listener, err := spiffetls.Listen(ctx, "tcp", "0.0.0.0:"+port, tlsconfig.AuthorizeID(clientID))
	if err != nil {
		fmt.Printf("error setting up spiffe: %v\n", err)
	}
	defer listener.Close()

	fmt.Printf("SPIFFE: listening on port %v\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error accepting connections: %v\n", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Printf("error reading incoming data: %v", err)
	}

	fmt.Printf("received in Spiffe Listener: %v", req)

	if _, err = conn.Write([]byte("Seawas")); err != nil {
		fmt.Printf("error sending response: %v", err)
		return
	}
}

func main() {
	listenPort := os.Getenv("SIMPLE_LISTEN_PORT")
	targetHost := os.Getenv("SIMPLE_TARGET_HOST")
	targetPort := os.Getenv("SIMPLE_TARGET_PORT")

	fmt.Printf("config: listening on: %v calling: %v:%v\n", listenPort, targetPort, targetHost)

	go initSpiffeServer(context.Background(), "1"+listenPort)

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
