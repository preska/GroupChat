/*
  To build the app
  > go build mp1.go
  
  To start the chat server
  > ./mp1 [PORT NUMBER]

  For client join the chat
  > telnet localhost [POST NUMBER]
  OR
  > telnet [IP ADDRESS] [PORT NUMBER]
*/

package main

import (
  "bufio"
  "fmt"
  "log"
  "net"
  "os"
  "os/exec"
)

func main() {
  if len(os.Args) == 2 {
    startServer()
  } else if len(os.Args) == 4 {
    //name := os.Args[1]
    portNum := os.Args[2]
    //numPeople := os.Args[3]
    
    fmt.Println(portNum)
    
    cmd := exec.Command("telnet", "localhost", portNum)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
      log.Fatal(err)
    }    
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
    // 3 cases: 
    // 1) Accept new connections
    // 2) Clear old connections
    // 3) Broadcast the chat messages to each client
    select {

    // 1) Accept new connections
    case conn := <- newConnectionChannel:

      log.Printf("Accepted new client, #%d", numClients)

      // Add this connection to the `clientMap` map
      clientMap[conn] = numClients
      numClients += 1

      // Read incoming chat messages from clients and push them onto the messages channel to broadcast to the other clients
      go func(conn net.Conn, clientId int) {
        reader := bufio.NewReader(conn)
        for {
          incoming, err := reader.ReadString('\n')
          if err != nil {
            break
          }
          messages <- fmt.Sprintf("Client %d > %s", clientId, incoming)
        }

        // When we encouter `err` reading, send this 
        // connection to `oldConnectionChannel` for removal.
        //
        oldConnectionChannel <- conn

      }(conn, clientMap[conn])

    // Accept messages from connected clients
    case message := <- messages:

      // Loop over all connected clients
      for conn, _ := range clientMap {

        // Send them a message in a go-routine
        // so that the network operation doesn't block
        go func(conn net.Conn, message string) {
          _, err := conn.Write([]byte(message))

          // If there was an error communicating
          // with them, the connection is dead.
          if err != nil {
            oldConnectionChannel <- conn
          }
        }(conn, message)
      }
      log.Printf("New message: %s", message)
      log.Printf("Broadcast to %d clients", len(clientMap))

    // Remove dead clients
    case conn := <-oldConnectionChannel:
      log.Printf("Client %d disconnected", clientMap[conn])
      delete(clientMap, conn)
    }
  }
}