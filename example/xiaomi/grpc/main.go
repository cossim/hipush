package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	"github.com/cossim/hipush/api/grpc/v1"
	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:7071", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := v1.NewPushServiceClient(conn)

	// Contact the server and print out its response.
	ctx := context.Background()
	req := &v1.PushRequest{
		Platform: "xiaomi",
		Tokens: []string{
			"xxx",
		},
		Title:            "cossim",
		Message:          "",
		Topic:            "",
		Key:              "",
		Category:         "",
		Sound:            "",
		Alert:            nil,
		Badge:            0,
		ThreadID:         "",
		Data:             nil,
		Image:            "",
		ID:               "",
		PushType:         "",
		AppID:            "",
		Priority:         0,
		ContentAvailable: false,
		MutableContent:   false,
		Development:      false,
		Option:           nil,
	}
	resp, err := c.Push(ctx, req)
	if err != nil {
		log.Fatalf("could not push: %v", err)
	}
	fmt.Println("Push response:", resp)
}
