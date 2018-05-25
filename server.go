package main

// https://gist.github.com/reterVision/33a72d70194d4a3c272e

import (
    "os"
    "./server"
)

func main() {
    var serverHost string
    if len(os.Args) == 2 {
        serverHost = os.Args[1]
    } else {
        serverHost = ":9999"
    }

    chatServer := server.NewUDPChatServer(serverHost)
    chatServer.Listen()
}

