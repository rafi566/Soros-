package main

import (
    "log"
    "net/http"
    "os"

    "soros/internal/server"
)

func main() {
    addr := getAddr()
    api := server.NewAPIService()

    log.Printf("starting API server on %s", addr)
    if err := http.ListenAndServe(addr, api.Router()); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}

func getAddr() string {
    if addr := os.Getenv("PORT"); addr != "" {
        return ":" + addr
    }
    return ":8080"
}
