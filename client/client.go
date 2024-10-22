package main

import (
	proto "ChittyChat/grpc"
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/user"
	"strconv"
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
		Name:    currentUser.Username + strconv.Itoa(rand.IntN(1000)),
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
			str := fmt.Sprintf("%s: %s (Lamport: %d)", msg.Username, msg.Msg, msg.Lamport)
			outputCh <- str
			mu.Unlock()

		case input, _ := <-inputCh:
			if input == "exit" {
				close(outputCh)
				LeaveServer(client, user)
				fmt.Println("Exiting...")
				return
			}
			msg1 := &proto.Message{
				Username: user.Name,
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
			inputCh <- "exit"
			close(inputCh)
			return
		}
		inputCh <- input
	}
}

func LeaveServer(client proto.ChittyChatClient, user *proto.UserJoin) {

	userl := &proto.UserLeave{
		Name:    user.Name,
		Lamport: user.Lamport,
	}
	_, err := client.LeaveServer(context.Background(), userl)
	if err != nil {
		log.Fatalf("did not work")
	}
}
