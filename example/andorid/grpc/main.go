package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	"github.com/cossim/hipush/api/pb/v1"
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

	ctx := context.Background()

	r := &v1.AndroidPushRequestData{
		Title:      "cossim",
		Content:    "hello",
		Topic:      "",
		TTL:        0,
		Priority:   "normal",
		CollapseID: "",
		Condition:  "",
		Sound:      "",
		Icon:       "",
		Data:       nil,
	}

	data, err := v1.ToStructPB(r)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON to structpb.Struct: %v", err)
	}

	req := &v1.PushRequest{
		AppID:    "xxx",
		AppName:  "cossim",
		Platform: "android",
		Token:    []string{"xxx"},
		Data:     data,
		Option:   nil,
	}

	fmt.Println("req => ", req.Data)

	resp, err := c.Push(ctx, req)
	if err != nil {
		log.Fatalf("could not push: %v", err)
	}
	fmt.Println("Push response:", resp)
}
