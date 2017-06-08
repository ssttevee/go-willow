package maple

import (
	"net"
)

const (
	MapleVersion uint16 = 185
	MaplePatch   string = "2"
)

type MapleServer struct {
	lis net.Listener
}

func (ni *MapleServer) Start() {
	for i := 0; ; i++ {
		conn, err := ni.lis.Accept()
		if err != nil {
			panic(err)
		}

		go newConnectionHandler(conn).Start(i)
	}
}

func NewServer(lis net.Listener) *MapleServer {
	return &MapleServer{
		lis: lis,
	}
}
