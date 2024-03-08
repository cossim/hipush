package main

import (
	"fmt"
	vv "github.com/cossim/vivo-push"
)

func main() {
	client, err := vv.NewClient("105562603", "ea6a826256eb2896a2f36743d859dbd0", "b14a60d5-bcb1-4836-8d88-15f2ac040eee")
	if err != nil {
		return
	}

	// 单推
	msg1 := vv.NewVivoMessage("hi baby", "hi")
	_, err = client.Send(msg1, "v2-CQanxnrM-uZu6i_y_E3PpymvRvSJhhFjaOwQxbGE-jJ0BNLa0IUm")
	if err != nil {
		fmt.Println(err)
		return
	}

	//// 群推
	//msg2 := vv.NewListPayloadMessage("hello baby", "hello")
	//_, err = client.SendList(msg2, []string{"regID1", "regID2"})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	////全量推送
	//msg3 := vv.NewListPayloadMessage("hi all baby", "hi all")
	//_, err = client.SendAll(msg3)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	return
}
