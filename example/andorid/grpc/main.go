package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
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

	ap := &v1.AndroidPushRequestData{
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

	marshaler := &jsonpb.Marshaler{}
	jsonString, err := marshaler.MarshalToString(ap)
	if err != nil {
		log.Fatalf("Failed to marshal struct to JSON: %v", err)
	}

	structValue := &structpb.Struct{}
	err = jsonpb.UnmarshalString(jsonString, structValue)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON to structpb.Struct: %v", err)
	}

	req := &v1.PushRequest{
		AppID:    "xxx",
		AppName:  "cossim",
		Platform: "android",
		Token:    []string{"xxx"},
		Data:     structValue,
		Option:   nil,
	}

	fmt.Println("req => ", req.Data)

	resp, err := c.Push(ctx, req)
	if err != nil {
		log.Fatalf("could not push: %v", err)
	}
	fmt.Println("Push response:", resp)
}
