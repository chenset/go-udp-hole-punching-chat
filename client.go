package main

import (
    "os"
    "log"
    "./client"
)

func main() {
    if len(os.Args) < 5 {
        log.Fatal("Usage: ", os.Args[0], " port serverAddr username peername")
    }

    port := os.Args[1]
    serverAddr := os.Args[2]
    username := os.Args[3]
    peer := os.Args[4]

    log.Print("port: ", port, ", serverAddr: ", serverAddr, ", username: ", username, ", peer: ", peer)

    chatClient := client.NewUDPChatClient(os.Args[1], os.Args[2], os.Args[3], os.Args[4])
    chatClient.Start()
}

