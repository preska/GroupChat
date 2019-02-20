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
  //name := os.Args[1]

  port := os.Args[2]
  //numPeople := os.Args[3]

  host := "localhost"
  conn, err := net.Dial("tcp", host + ":" + port)

  if err != nil {
    if _, t := err.(*net.OpError); t {
      fmt.Println("Problem connecting to the server.")
    } else {
      fmt.Println(err.Error())
    }
    os.Exit(1)
  }

  fmt.Printf("READY\n")

  go read(conn)

  for {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("> ")
    text, _ := reader.ReadString('\n')

    conn.SetWriteDeadline(time.Now().Add(time.Second))

    _, err := conn.Write([]byte(text))
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
      text := scanner.Text()

      fmt.Printf("%s\n> ", text)

      if !ok {
        fmt.Println("Server is down.")
        break
      }
    }
  }
}