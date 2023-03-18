package rpc

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/globaldce/globaldce-gateway/content"
)

type Request struct {
    Jsonrpc string          `json:"jsonrpc"`
    Method  string          `json:"method"`
    Params  []string `json:"params"`
    ID      int             `json:"id"`
}

type Response struct {
    Jsonrpc string      `json:"jsonrpc"`
    Result  *Result      `json:"result,omitempty"`//interface{} `json:"result,omitempty"`
    Error   *Error      `json:"error,omitempty"`
    ID      int         `json:"id"`
}

type Error struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type Result struct {
    Type    string      `json:"type"`
    Data    interface{} `json:"data,omitempty"`
}


var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // allow all origins
    },
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

var Mncc *content.ContentClient

func RPCInit(tmpMncc *content.ContentClient) {
    Mncc=tmpMncc

    http.HandleFunc("/rpc", RPCHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func RPCHandler(w http.ResponseWriter, r *http.Request) {
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
                Result:  &Result{Type:"pong"},
                ID:      request.ID,
            }
        case "CacheTorrent":
            //log.Printf("cacheTorrent params %s",request.Params[0])
            
            tmpResult:=runCacheTorrent(request.Params)
            response = Response{
                Jsonrpc: "2.0",
                Result:  &tmpResult,
                ID:      request.ID,
            }
        //maincontentclient.AddCacheTorrentRequest("cooldapp","",tmpmagnet)
        //maincontentclient.ProtorizeTorrentPiecesInterval(tmpmagnet,".mp4",0,20)
        //maincontentclient.ProtorizeTorrentAllPieces(tmpmagnet,".mp4")
        case "ProtorizeTorrentPiecesInterval":
            tmpResult:=runProtorizeTorrentPiecesInterval(request.Params)
            response = Response{
                Jsonrpc: "2.0",
                Result:  &tmpResult,//"ProtorizeTorrentPiecesIntervalSuccess",
                ID:      request.ID,
            }
        case "ProtorizeTorrentDurationPercentageInterval":
            tmpResult:=runProtorizeTorrentDurationPercentageInterval(request.Params)
            response = Response{
                Jsonrpc: "2.0",
                Result:  &tmpResult,//"ProtorizeTorrentPiecesIntervalSuccess",
                ID:      request.ID,
            }
        case "GetTorrentDetails":
            tmpResult:=runGetTorrentDetails(request.Params)
            response = Response{
                Jsonrpc: "2.0",
                Result:  &tmpResult,//tmpTorrentDetailsString,
                ID:      request.ID,
            }
        //************************
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
