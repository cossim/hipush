package main

import (
	"context"
	"fmt"
	xp "github.com/yilee/xiaomi-push"
)

var client = xp.NewClient("WX9E+eHBbql3gl6RYHEGXQ==", []string{"com.dootask.task"})
var regID1 = "LCDIsLAcA6ZLexnPIfZp4bGfDTShKdoRUcxP8pii8OxMOuHNcNyIKlJL0QQYiMuN"

func main() {
	msg1 := xp.NewAndroidMessage("title", "body").
		SetPayload("this is payload1")
	resp, err := client.Send(context.Background(), msg1, regID1)
	if err != nil {
		panic(err)
	}

	fmt.Println("resp => ", resp)

}
