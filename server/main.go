package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

var clients map[net.Conn]string

func handleConnectionError(connType string, user string, err error) {
	fmt.Printf("Connection error %s for user [%s]: %s", connType, user, err)
}

func portListener() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error listening: ", err)
		return
	}
	defer listener.Close() // With 'defer', ensure the listener is closed when the program exits
	fmt.Println("Server started on port 8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	defer conn.Close()

	err := setUpDisplayName(conn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("New connection from ", conn.RemoteAddr())
	for client := range clients {
		if client != conn {
			_, err := client.Write([]byte("[INFO] [" + time.Now().String() + "] " + clients[conn] + " has logged in.\n"))
			if err != nil {
				fmt.Printf("Error sending message to %s\n", client.RemoteAddr())
			}
		} else {
			_, err := client.Write([]byte("You have logged in.\n"))
			if err != nil {
				fmt.Printf("Error sending message to %s\n", client.RemoteAddr())
			}
		}
	}

	buffer := make([]byte, 1024)                 // will store user input
	commandRegex := regexp.MustCompile(`/(\w+)`) // regex for command inputs "/.+".
	for {                                        // handle user input
		n, err := conn.Read(buffer) // n represents the number of bytes received
		if err != nil {             // handle when user disconnects
			fmt.Printf("Connection from %s lost.", conn.RemoteAddr())
			for client := range clients { // Send message to other clients if someone disconnects
				if client != conn {
					_, err := client.Write([]byte("[INFO] [" + time.Now().String() + "] " + clients[conn] + " disconnected.\n"))
					if err != nil {
						fmt.Printf("Error sending message to %s\n", client.RemoteAddr())
					}
				}
			}
			break
		}
		/*
			The buffer will have the bytes of the message and the rest will be 0s
			e.g. ['h', 'e', 'l', 'l', 'o', '\n'] will basically be [104 101 108 108 111 10 0 0 0 0 ...] // length 1024
			buffer[:n] is basically getting the buffer values until they reach 0 (end of)
		*/
		message := string(buffer[:n])
		commandMatch := commandRegex.FindStringSubmatch(message)
		if commandMatch != nil {
			commands(commandMatch[1], conn) // command[1] is basically the match group of the regex. command[0] is the entire match
			continue
		}
		fmt.Printf("[%s]: %s\n", conn.RemoteAddr(), message)

		// Print message to other clients
		for client := range clients {
			if client != conn { // check if client is not the current connection
				_, err := client.Write([]byte("[MSG][" + clients[conn] + "]: " + message + "\n")) // Send message to other clients
				if err != nil {
					fmt.Printf("Error sending message to %s\n", client.RemoteAddr())
				}
			}
		}
	}

	// Remove disconnected client from list of clients
	delete(clients, conn)
}

func setUpDisplayName(conn net.Conn) error {
	_, err := conn.Write([]byte("Please enter your user name:\n"))
	if err != nil {
		fmt.Printf("Error sending message to %s\n", conn.RemoteAddr())
		return err
	}
	for { // loop until user enters valid username
		displayName := make([]byte, 1024)
		n, err := conn.Read(displayName)
		if err != nil {
			fmt.Printf("Connection from %s lost.", conn.RemoteAddr())
			return err
		}
		if !existingUser(string(displayName[:n])) {
			clients[conn] = strings.TrimSuffix(string(displayName[:n]), "\n")
			break
		} else {
			_, err := conn.Write([]byte("[ERR] This username already exists. Please enter another name name:\n"))
			if err != nil {
				fmt.Printf("Error sending message to %s\n", conn.RemoteAddr())
				return err
			}
		}
	}
	return nil
}

func commands(command string, conn net.Conn) {
	switch command {
	case "connected_users": // show currently connected users
		_, err := conn.Write([]byte("#####################\nConnected users:\n"))
		if err != nil {
			fmt.Printf("Error sending message to %s\n", conn.RemoteAddr())
		}
		for client := range clients {
			if client != conn {
				_, err := conn.Write([]byte("# - " + clients[client] + "\n")) // Send message to other clients
				if err != nil {
					handleConnectionError("listing "+clients[client]+" when displaying connected users", clients[conn], err)
				}
			}
		}
		_, err = conn.Write([]byte("#####################\n"))
		if err != nil {
			fmt.Printf("Error sending message to %s\n", conn.RemoteAddr())
		}

	case "rename": // change user name
		_, err := conn.Write([]byte("Enter new name:\n"))
		if err != nil {
			handleConnectionError("prompt to enter new name", clients[conn], err)
		}
		oldName := clients[conn]
		newDisplayName := make([]byte, 1024)
		n, err := conn.Read(newDisplayName)
		if err != nil {
			handleConnectionError("reading new display name", clients[conn], err)
		}
		clients[conn] = strings.TrimSuffix(string(newDisplayName[:n]), "\n")
		for client := range clients {
			if client == conn {
				_, err = conn.Write([]byte("Your display name has been successfully changed from [" +
					oldName + "] to [" + clients[conn] + "].\n"))
				if err != nil {
					handleConnectionError("Confirmation msg for name change", clients[conn], err)
				}
			} else {
				_, err = client.Write([]byte("[INFO][" + oldName + "] " + "changed name to [" + clients[conn] + "].\n"))
				if err != nil {
					handleConnectionError("Info msg to users about user name change", clients[client], err)
				}
			}
		}
	case "dm": // send direct message to a user
		_, err := conn.Write([]byte("Who would you like to send a direct message to?\n"))
		if err != nil {
			handleConnectionError("prompt to enter user to DM.", clients[conn], err)
		}
		dmUser := make([]byte, 1024)
		n, err := conn.Read(dmUser)
		if err != nil {
			handleConnectionError("reading new display name", clients[conn], err)
		}
		dmUser = dmUser[:n]
		if !existingUser(string(dmUser)) {
			_, err = conn.Write([]byte("User [" + string(dmUser) + "] does not exist.\n"))
			if err != nil {
				handleConnectionError("prompt to enter DM.", clients[conn], err)
			}
			return
		} else if clients[conn] == string(dmUser) {
			_, err = conn.Write([]byte("Cannot send DM to yourself.\n"))
			if err != nil {
				handleConnectionError("prompt to enter DM.", clients[conn], err)
			}
			return
		}
		_, err = conn.Write([]byte("Enter message to send to [" + string(dmUser) + "]:\n"))
		if err != nil {
			handleConnectionError("prompt to enter DM message.", clients[conn], err)
		}
		dm := make([]byte, 1024)
		n, err = conn.Read(dm)
		if err != nil {
			handleConnectionError("reading new display name", clients[conn], err)
		}
		for client, userName := range clients {
			if userName == string(dmUser) {
				_, err = client.Write([]byte("[MSG] [" + clients[conn] + "] (DM): " + string(dm[:n]) + "\n"))
				if err != nil {
					handleConnectionError("sending DM.", clients[conn], err)
				}
			}
		}

	default:
		_, err := conn.Write([]byte("Unknown command.\n"))
		if err != nil {
			fmt.Printf("Error sending message to %s\n", conn.RemoteAddr())
		}
	}
}

func existingUser(user string) bool {
	// Check if specified userName exists
	for _, userName := range clients {
		if user == userName {
			return true
		}
	}
	return false
}

func main() {
	clients = make(map[net.Conn]string)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go portListener()

	<-sigChan
	fmt.Printf("Shutting down server...")
	for client := range clients {
		client.Write([]byte("[INFO] Server is shutting down."))
		client.Close()
	}
	os.Exit(0)
}
