package main

import (
	"mynet/blitz"
	"mynet/network/tcp"
)

func main() {
	var app blitz.App
	var proto blitz.Protocol
	proto.SetNotify(&app)
	tcp.ListenAndAccept(":8080", &proto)
}
