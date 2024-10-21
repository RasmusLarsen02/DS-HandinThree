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
	wg.Add(2)

	go Listen(stream, msgCh, &wg)

	go ReadInput(inputCh, &wg, client)

	go PrintMessages(outputCh)

	for {
		select {
		case msg := <-msgCh:
			mu.Lock()
			str := fmt.Sprintf("Message: %s: %s (Lamport: %d)", msg.Username, msg.Msg, msg.Lamport)
			outputCh <- str
			mu.Unlock()

		case input, ok := <-inputCh:
			if !ok {
				close(outputCh)
				fmt.Println("Exiting...")
				return
			}
			msg1 := &proto.Message{
				Username: currentUser.Name,
				Msg:      input,
				Lamport:  0,
			}
			_, _ = client.SendMessage(context.Background(), msg1)

		}
	}
	wg.Wait()
}

func PrintMessages(outputCh chan string) {
	for {
		msg := <-outputCh
		mu.Lock()
		fmt.Println(msg)
		mu.Unlock()
	}

}

func Listen(stream proto.ChittyChat_JoinServerClient, msgCh chan<- *proto.Message, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return
		}
		if msg != nil {
			msgCh <- msg
		}
	}
}

func ReadInput(inputCh chan<- string, wg *sync.WaitGroup, client proto.ChittyChatClient) {
	time.Sleep(2 * time.Second)
	defer wg.Done()
	reader := bufio.NewReader(os.Stdin)

	for {
		time.Sleep(1 * time.Second)
		mu.Lock()
		fmt.Println("Write your message or type exit to leave: ")
		mu.Unlock()
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "exit" {
			close(inputCh)
			LeaveServer(client)
			return
		}
		inputCh <- input
	}
}

func LeaveServer(client proto.ChittyChatClient) {
	currentUser, _ := user.Current()
	user := &proto.UserJoin{
		Name:    currentUser.Name,
		Lamport: 0,
	}

	_, err := client.JoinServer(context.Background(), user)
	if err != nil {
		log.Fatalf("did not work")
	}
}
