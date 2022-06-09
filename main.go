package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
)

func main() {
	fmt.Println("Hi")

	rest := Rest{}

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "81"), hystrixStreamHandler)

	rest.start()
}
