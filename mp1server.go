/*
  To build the app
  > go build mp1.go
  
  To start the chat server
  > ./mp1 [PORT NUMBER]
*/

package main

import (
  "bufio"
  "fmt"
  "net"
  "os"
)

func main() {
  if len(os.Args) == 2 {
    startServer()
  } else {
    fmt.Println("Enter a port number.")
    os.Exit(1)
  }
}

func startServer() {
  numClients := 0  // Number of people currently connected to the chat server
  clientMap := make(map[net.Conn]int) // Map of people connected to the chat server (key = net.Conn object, value = client ids)
  newConnectionChannel := make(chan net.Conn) // Channel of new connections
  oldConnectionChannel := make(chan net.Conn) // Channel of old connections to remove from the clientMap
  messages := make(chan string) // Channel of chat messages to broadcast

  portNum := os.Args[1]
  
  // Start the TCP server
  server, err := net.Listen("tcp", ":" + portNum)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  // Accept and start a thread for each connection
  go func() {
    for {
      conn, err := server.Accept()
      if err != nil {
        fmt.Println(err)
        os.Exit(1)
      }
      // Push the new connection into the newConnectionChannel channel
      newConnectionChannel <- conn
    }
  }()

  // Keep the server running until interrupted 
  for {
    // 3 cases: (1) Accept new connections   (2) Clear old connections   (3) Broadcast the chat messages to each client
    select {
    // Accept new connections
    case conn := <- newConnectionChannel:
      clientMap[conn] = numClients // Add this connection to the `clientMap` map
      numClients++

      // Read incoming chat messages from clients and push them onto the messages channel to broadcast to the other clients
      go func(conn net.Conn, clientId int) {
        messages <- fmt.Sprintf("Client %d has joined\n", numClients - 1)
        reader := bufio.NewReader(conn)
        for {
          message, err := reader.ReadString('\n')
          if err != nil {
            break
          }
          messages <- fmt.Sprintf("Client %d: %s", clientId, message)
        }
        oldConnectionChannel <- conn
        messages <- fmt.Sprintf("Client %d has left\n", clientMap[conn])

      }(conn, clientMap[conn])

    // Accept messages from connected clients
    case message := <- messages:
      // Loop over all connected clients
      for conn, _ := range clientMap {
        //if currClientId != clientMap[conn] {
          // Send the connected clients a message except to the client that sent the message
          go func(conn net.Conn, message string) {
            _, err := conn.Write([]byte(message))

            if err != nil { // connection is dead
              oldConnectionChannel <- conn
            }
          }(conn, message)
        //}
      }
      fmt.Printf("%s", message)

    // Remove old clients
    case conn := <- oldConnectionChannel:
      delete(clientMap, conn)
    }
  }
}