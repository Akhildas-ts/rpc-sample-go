package main

import (
	"fmt"
	"net"
	"net/rpc"
	"time"
)

// Args holds the arguments for the RPC method
type Args struct {
	A, B int
}

// Calculator is the service with methods that can be called remotely
type Calculator struct{}

// Sum is an RPC method for adding two numbers
func (c *Calculator) Sum(args *Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}

func startServer() {
	calculator := new(Calculator)
	rpc.Register(calculator) // Register the Calculator service

	// Listen on port 1234
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error setting up listener:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 1234...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func startClient() {
	// Give the server a moment to start up
	time.Sleep(1 * time.Second)

	// Connect to the RPC server
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error connecting to RPC server:", err)
		return
	}
	defer client.Close()

	// Prepare arguments and a variable to store the reply
	args := &Args{A: 5, B: 3}
	var reply int

	// Call the Sum method on the server
	err = client.Call("Calculator.Sum", args, &reply)
	if err != nil {
		fmt.Println("Error calling Calculator.Sum:", err)
		return
	}

	fmt.Println("Result of Calculator.Sum:", reply) // Expected output: Result of Calculator.Sum: 8
}

func main() {
	// Run the server in a separate goroutine
	go startServer()

	// Run the client
	startClient()
}
