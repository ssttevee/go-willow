package main

import (
	"fmt"
	"net"
	"os"

	"golang.ssttevee.com/willow/internal/maple"
	_ "golang.ssttevee.com/willow/internal/maple/handlers"
)

// netsh int ip add addr 1 8.31.99.141/0 st=ac sk=tr
// "C:\AriesMS\MapleStory.exe" WebStart some:auth:token 8.31.99.141 8484
// netsh int ip delete addr 1 8.31.99.141

func main() {
	lis, err := net.Listen("tcp", ":8484")
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stdout, "Listening on %s\n", lis.Addr().String())

	maple.NewServer(lis).Start()
}
