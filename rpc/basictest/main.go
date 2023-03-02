package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

type Request struct {
    Jsonrpc string          `json:"jsonrpc"`
    Method  string          `json:"method"`
    Params  json.RawMessage `json:"params"`
    ID      int             `json:"id"`
}

type Response struct {
    Jsonrpc string      `json:"jsonrpc"`
    Result  interface{} `json:"result,omitempty"`
    Error   *Error      `json:"error,omitempty"`
    ID      int         `json:"id"`
}

type Error struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func main() {
    http.HandleFunc("/rpc", rpcHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func rpcHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()

    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }

        // Parse incoming JSON-RPC request
        var request Request
        err = json.Unmarshal(message, &request)
        if err != nil {
            log.Println(err)
            return
        }

        // Handle incoming JSON-RPC request
        var response Response
        switch request.Method {
        case "ping":
            response = Response{
                Jsonrpc: "2.0",
                Result:  "pong",
                ID:      request.ID,
            }
        case "hello":
            response = Response{
                Jsonrpc: "2.0",
                Result:  "world",
                ID:      request.ID,
            }
        default:
            response = Response{
                Jsonrpc: "2.0",
                Error: &Error{
                    Code:    -32601,
                    Message: "Method not found",
                },
                ID: request.ID,
            }
        }

        // Send response back to client
        responseBytes, err := json.Marshal(response)
        if err != nil {
            log.Println(err)
            return
        }

        err = conn.WriteMessage(messageType, responseBytes)
        if err != nil {
            log.Println(err)
            return
        }
    }
}
