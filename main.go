package main

import "time"

func main() {
	go server()
	time.Sleep(100 * time.Millisecond)
	client()
	return
}
