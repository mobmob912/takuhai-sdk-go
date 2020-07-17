package main

import (
	"io/ioutil"
	"log"

	takuhai "github.com/tockn/takuhai-sdk-go"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	e := takuhai.NewClient()
	return e.Run(echo)
}

func echo(c *takuhai.Context) {
	respBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Fail([]byte("invalid output"))
		return
	}
	log.Println("echo: ", string(respBody))
	c.Next(respBody)
}
