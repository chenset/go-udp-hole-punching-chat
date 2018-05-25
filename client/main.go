package main

import (
    "os"
    "log"
    "fmt"
    "net"
    "encoding/json"
    "time"
)

type ChatRequest struct {
    Action string
    Username string
    Message string
}

func main() {
    if len(os.Args) < 5 {
        log.Fatal("Usage: ", os.Args[0], " port serverAddr username peername")
    }

    port := fmt.Sprintf(":%s", os.Args[1])
    serverAddr := os.Args[2]
    username := os.Args[3]
    peer := os.Args[4]

    log.Print("port: ", port, ", serverAddr: ", serverAddr, ", username: ", username, ", peer: ", peer)

    buf := make([]byte, 2048)

    saddr, err := net.ResolveUDPAddr("udp4", serverAddr)
    if err != nil {
        log.Print("Resolve server address failed.")
        log.Fatal(err)
    }
    log.Print("saddr: ", saddr)

    addr, err := net.ResolveUDPAddr("udp4", port)
    if err != nil {
        log.Print("Resolve locl addres failed.")
        log.Fatal(err)
    }
    log.Print("addr: ", addr)

    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        log.Print("Listen UDP failed.")
        log.Fatal(err)
    }

    initChatRequest := ChatRequest{
        "New",
        username,
        "",
    }
    jsonRequest, err := json.Marshal(initChatRequest)
    if err != nil {
        log.Print("Marshal Register information failed.")
        log.Fatal(err)
    }
    _, err = conn.WriteToUDP(jsonRequest, saddr)
    if err != nil {
        log.Fatal(err)
    }

    log.Print("Waiting for server response...")
    _, _, err = conn.ReadFromUDP(buf)
    if err != nil {
        log.Print("Register to server failed.")
        log.Fatal(err)
    }

    connectChatRequest := ChatRequest{
        "Get",
        username,
        peer,
    }
    jsonRequest, err = json.Marshal(connectChatRequest)
    if err != nil {
        log.Print("Marshal connection information failed.")
        log.Fatal(err)
    }

    var serverResponse ChatRequest
    for i := 0; i < 3; i++ {
        conn.WriteToUDP(jsonRequest, saddr)
        n, _, err := conn.ReadFromUDP(buf)
        if err != nil {
            log.Print("Get peer address from server failed.")
            log.Fatal(err)
        }
        err = json.Unmarshal(buf[:n], &serverResponse)
        if err != nil {
            log.Print("Unmarshal server response failed.")
            log.Fatal(err)
        }
        if serverResponse.Message != "" {
            break
        }
        time.Sleep(10 * time.Second)
    }

    if serverResponse.Message == "" {
        log.Fatal("Cannot get peer's address")
    }
    log.Print("Peer address: ", serverResponse.Message)
    peerAddr, err := net.ResolveUDPAddr("udp4", serverResponse.Message)
    if err != nil {
        log.Print("Resolve peer addres failed.")
        log.Fatal(err)
    }

    go listen(conn)
    for {
        fmt.Print("Input message: ")
        message := make([]byte, 2048)
        fmt.Scanln(&message)
        messageRequest := ChatRequest{
            "Chat",
            username,
            string(message),
        }
        jsonRequest, err = json.Marshal(messageRequest)
        if err != nil {
            log.Print("Error: ", err)
            continue
        }
        conn.WriteToUDP(jsonRequest, peerAddr)
    }
}

func listen(conn *net.UDPConn) {
    for {
        buf := make([]byte, 2048)
        n, _, err := conn.ReadFromUDP(buf)
        if err != nil {
            log.Print(err)
            continue
        }

        var message ChatRequest
        err = json.Unmarshal(buf[:n], &message)
        if err != nil {
            log.Print(err)
            continue
        }
        fmt.Println(message.Username, ":", message.Message)
    }
}
