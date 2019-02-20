/*
  To build the app
  > go build mp1client.go

  To create a client
  > ./mp1client [NAME] [PORT NUMBER] [NUMBER OF PEOPLE]
*/

package main

import (
  "bufio"
  "fmt"
  "net"
  "os"
  "time"
)

func main() {
  name := os.Args[1] + "\n"
  port := os.Args[2]
  //num := os.Args[3]

  host := "localhost"
  conn, err := net.Dial("tcp", host + ":" + port) // connect to a tcp server at the given port

  if err != nil {
    if _, c := err.(*net.OpError); c {
      fmt.Println("Problem connecting to the server.")
    } else {
      fmt.Println(err.Error())
    }
    os.Exit(1)
  }

  fmt.Printf("READY\n")

  go read(conn) // get message from user input

  // send name of client to the server
  _, err = conn.Write([]byte(name)) 
  if err != nil {
    fmt.Println(err.Error())
  }

  for {
    reader := bufio.NewReader(os.Stdin)
    message, _ := reader.ReadString('\n')

    conn.SetWriteDeadline(time.Now().Add(time.Second)) // detect a client that is not reading data (or reading data at an unacceptable rate)

    _, err := conn.Write([]byte(message)) // send message to the server
    if err != nil {
      fmt.Println(err.Error())
      break
    }
  }
}

func read(conn net.Conn) {
  for {
    scanner := bufio.NewScanner(conn)
    for {
      ok := scanner.Scan()
      message := scanner.Text()

      fmt.Printf("%s\n> ", message)

      if !ok {
        fmt.Println("Server is down.")
        break
      }
    }
  }
}