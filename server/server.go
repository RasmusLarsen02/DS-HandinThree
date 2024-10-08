package main

import (
	proto "ChittyChat/grpc"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type ChittyChatServer struct {
	proto.UnimplementedChittyChatServer
}

func (cs *ChittyChatServer) SendMessage(ctx context.Context, in *proto.Message) (*proto.Empty, error) {
	return &proto.Empty{}, nil
}

func main() {
	cs := &ChittyChatServer{}
	cs.start_server()
}

func (cs *ChittyChatServer) start_server() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Did not work")
	}

	proto.RegisterChittyChatServer(grpcServer, cs)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("Did not work")
	}
}
