package main

import (
	proto "ChittyChat/grpc"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var mu sync.Mutex

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working")
	}
	var wg sync.WaitGroup

	msgCh := make(chan *proto.Message)
	inputCh := make(chan string)
	outputCh := make(chan string)

	currentUser, _ := user.Current()

	client := proto.NewChittyChatClient(conn)
	user := &proto.UserJoin{
		Name:    currentUser.Name,
		Lamport: 0,
	}

	stream, _ := client.JoinServer(context.Background(), user)
	go Listen(stream)

	time.Sleep(1000)
	msg1 := &proto.Message{
		Username: currentUser.Username,
		Msg:      "Hello Every One!",
		Lamport:  0,
	}
	_, _ = client.SendMessage(context.Background(), msg1)

	for {
	}
}

func Listen(stream grpc.ServerStreamingClient[proto.Message]) {
	for {
		msg, _ := stream.Recv()
		if msg == nil {
			continue
		}
		fmt.Printf("%s : %s (at time %d)", msg.Username, msg.Msg, msg.Lamport)
	}
}
