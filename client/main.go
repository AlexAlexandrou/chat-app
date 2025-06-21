package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/AlexAlexandrou/chat-app/client/styles"
	"github.com/AlexAlexandrou/chat-app/config"
)

func main() {
	config.LoadEnvVars()

	serverAddress := config.GetEnv("SERVER_ADDRESS", "localhost")
	port := config.GetEnv("SERVER_PORT", "8080")
	conn, err := net.Dial("tcp", serverAddress+":"+port) // The address will change to an actual address
	if err != nil {
		fmt.Printf("%s", styles.StyleMsg("[ERR] Error connecting to server."))
		return
	}
	defer conn.Close()

	go func() {
		serverScanner := bufio.NewScanner(conn) // reads from server
		for serverScanner.Scan() {
			serverMsg := serverScanner.Text()
			colouredMsg := styles.StyleMsg(serverMsg)
			fmt.Printf("\r%s\n> ", colouredMsg)
		}

		// Stop client when connection to server is lost
		fmt.Printf("\n")
		fmt.Printf("%s", styles.StyleMsg("[ERR] Connection to Server lost. Shutting down client."))
		os.Exit(0)
	}()

	inputScanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">")
	for inputScanner.Scan() {
		fmt.Print(">")
		input := inputScanner.Text()
		if input == "/exit" {
			return
		}
		fmt.Fprint(conn, input)
	}
}
