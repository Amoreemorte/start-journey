package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	request := []string{"First message"}
	var responce string

	err = client.Call("Server.Send", request, &responce)
	if err != nil {
		log.Print(err)
	}
	err = client.Call("Server.Send", []string{"Second message"}, &responce)
	if err != nil {
		log.Print(err)
	}

	err = client.Call("Server.Messages", struct{}{}, &request)
	if err != nil {
		log.Print(err)
	}

	fmt.Printf("%v", request)
}
