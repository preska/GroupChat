package main

import (
    "fmt"
    "net"
    "os"
)

const (
    CONN_HOST = "localhost"
    CONN_TYPE = "tcp"
)

func main() {
    startServer()
}

func startServer() {
  //Get command line arguments (Name, port number, and number of people in chat)
  //var user_name = os.Args[1]
  var port_num = os.Args[2]
  //var num_people = os.Args[3]

  // Listen for incoming connections.
  l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+ port_num)
  if err != nil {
    fmt.Println("Error listening:", err.Error())
    os.Exit(1)
  }
  
  // Close the listener when the application closes.
  defer l.Close()
  fmt.Println("READY")
  fmt.Println("#Listening on " + CONN_HOST + ":" + port_num)
  for {
    // Listen for an incoming connection.
    conn, err := l.Accept()
    if err != nil {
      fmt.Println("Error accepting: ", err.Error())
      os.Exit(1)
    }
    // Handle connections in a new goroutine.
    go handleRequest(conn)
  }
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
  // Make a buffer to hold incoming data.
  buf := make([]byte, 1024)
  // Read the incoming connection into the buffer.
  _, err := conn.Read(buf)
  if err != nil {
    fmt.Println("Error reading:", err.Error())
  }
  // Send a response back to person contacting us.
  conn.Write([]byte("Message received."))
  // Close the connection when you're done with it.
  conn.Close()
}