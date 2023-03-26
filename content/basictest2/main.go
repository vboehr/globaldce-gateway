package main

import (
	"fmt"
	"github.com/globaldce/globaldce-gateway/content"
	"context"
	//"time"
)

func main() {
	fmt.Println("Hello")

	// Create a context with a cancellation function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	maincontentclient:=content.Newcontentclient(ctx,"./")
	go maincontentclient.Initcontentclient()

	fmt.Println(maincontentclient.ScanDirectory("cooldapp","xxx"))
	
	cancel()

}