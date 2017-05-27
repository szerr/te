package main

import (
    "github.com/gorilla/websocket"
    "log"
    "net/http"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:     1024,
    WriteBufferSize:    1024,
}

var MES_CH = make(chan []byte, 1024)
var Roon = make(map[*websocket.Conn]string)

func webSocketHandler(w http.ResponseWriter, r *http.Request){
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    messageType, message, err := ws.ReadJSON()
    if err != nil {
        log.Println(err)
        return
    }
    if messageType != websocket.TextMessage {
        return
    }
    for {
        messageType, message, err := ws.ReadJSON()
        if err != nil {
            log.Println(err)
            return
        }
        MES_CH <- message
        if messageType != websocket.TextMessage {
            return
        }
    }
}

func broHandler(ch chan []byte){
    for {
        log.Println(<-ch, 111)
    }
    //err = ws.WriteMessage(messageType, message)
    //if err != nil {
    //    log.Println(err)
    //    return
    //}
}

func main() {
    fs := http.FileServer(http.Dir("./status"))
    http.Handle("/", fs)
    http.HandleFunc("/ws", webSocketHandler)
    go broHandler(MES_CH)
    log.Println("http server started on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ",err)
    }
}

