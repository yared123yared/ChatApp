package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"github.com/gorilla/websocket"
)

type M map[string]interface{

}

const MESSAGE_NEW_USER="New User"
const MESSAGE_CHAT ="Chat"
const MESSAGE_LEAVE="Leave"

var connections=make([]*WebSocketConnection, 0)

type SocketResponse struct {
	From string
	Type string
	Message string
}
type  WebSocketConnection struct {
	*websocket.Conn
	Username string
}
type  SocketPayload struct{
	Message string
}

func main(){
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		content,err:=ioutil.ReadFile("index.html")
		if err!=nil{
			http.Error(writer,"Could not open requested file", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(writer, "%s",content)

	})

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		currentGorillaConn, err:=websocket.Upgrade(writer,request,writer.Header(),1024,1024)
		if err !=nil{
			http.Error(writer,"Could not open websocket connection", http.StatusBadRequest)

		}
		userName:=request.URL.Query().Get("username")
		currentConn:=WebSocketConnection{Conn:currentGorillaConn,Username: userName}
		connections=append(connections,&currentConn)

		go handleIO(&currentConn,connections)

	})

	fmt.Println("Server startgin at :8000")
	http.ListenAndServe(":8000", nil)

}
func handleIO(currentConn *WebSocketConnection, connections []*WebSocketConnection){
	defer func() {
		if r:=recover();r !=nil{
			log.Println("ERROR", fmt.Sprintf("%v",r))
		}
	}()
	broadcastMessage(currentConn,MESSAGE_NEW_USER,"")
	for{
		payload:=SocketPayload{}
		err:=currentConn.ReadJSON(&payload)
		if err!=nil{
			if strings.Contains(err.Error(),"websocket:close"){
				broadcastMessage(currentConn,MESSAGE_LEAVE,"")
				ejectConnection(currentConn)
				return
			}
			log.Println("ERROR" , err.Error())
			continue
		}
		broadcastMessage(currentConn,MESSAGE_CHAT,payload.Message)

	}

}
func ejectConnection(currentConn *WebSocketConnection){
	filtered:=
}
func broadcastMessage(currentConn *WebSocketConnection,kind,message string){

}

