package maple

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"golang.ssttevee.com/willow/internal/common"
	"golang.ssttevee.com/willow/internal/maple/codec"
	"golang.ssttevee.com/willow/internal/maple/packet"
)

// aesKey is the encryption key extracted from MapleStory v185.2 game files
var aesKey = []byte{
	0x29, 0x00, 0x00, 0x00,
	0xF6, 0x00, 0x00, 0x00,
	0x18, 0x00, 0x00, 0x00,
	0x5E, 0x00, 0x00, 0x00,
	0xCA, 0x00, 0x00, 0x00,
	0x5A, 0x00, 0x00, 0x00,
	0x40, 0x00, 0x00, 0x00,
	0x61, 0x00, 0x00, 0x00,
}

type ConnectionHandler struct {
	conn   net.Conn
	number int
}

func newConnectionHandler(conn net.Conn) *ConnectionHandler {
	return &ConnectionHandler{
		conn: conn,
	}
}

func (client *ConnectionHandler) Start(n int) {
	quit := make(chan bool)
	pingTicker := time.NewTicker(30 * time.Second)

	defer func() {
		client.conn.Close()
		pingTicker.Stop()
		quit <- true
	}()

	client.number = n

	fmt.Fprintf(os.Stdout, "[#%d] Received connection from: %s\n", client.number, client.conn.RemoteAddr().String())

	c, err := codec.New(aesKey, MapleVersion)
	if err != nil {
		fmt.Fprintf(os.Stdout, "[#%d][ERRO] %v\n", client.number, err)
		return
	}

	helloPacket := packet.NewHelloPacket(MapleVersion, MaplePatch, c.ReceiverIV(), c.SenderIV())

	fmt.Fprintf(os.Stdout, "[#%d][SEND] Sent Handshake\n", client.number)
	client.conn.Write(helloPacket.Bytes())

	packetWriter := common.PacketWriterFunc(func(packet common.Packet) error {
		printPacket(client.number, "SEND", packet)
		return c.WritePacket(client.conn, packet)
	})

	// start goroutine for pings
	go func() {
		for {
			select {
			case <-pingTicker.C:
				if err := packetWriter.Write(packet.NewPacketWriter(common.SendOpCode_Ping)); err != nil {
					fmt.Fprintf(os.Stdout, "[#%d][ERRO] %v\n", client.number, err)
					break
				}
			case <-quit:
				break
			}
		}
	}()

	for {
		p, err := c.ReadPacket(client.conn)
		if err != nil {
			if err == io.EOF {
				fmt.Fprintf(os.Stdout, "[#%d] Client closed connection\n", client.number)
			} else {
				fmt.Fprintf(os.Stdout, "[#%d][ERRO] %v\n", client.number, err)
				fmt.Fprintf(os.Stdout, "[#%d] Closing connection\n", client.number)
			}

			break
		}

		printPacket(client.number, "RECV", p)

		if err := handlePacket(p, packetWriter); err != nil {
			if err == ErrUnhandledPacket {
				fmt.Fprintf(os.Stdout, "[#%d][WARN] %v\n", client.number, err)
			}
		}
	}
}

func printPacket(clientId int, mode string, packet common.Packet) {
	fmt.Fprintf(
		os.Stdout,
		"[#%d][%s] [0x%X] %#v\n",
		clientId,
		mode,
		packet.OpCode(),
		packet.Bytes()[2:],
	)
}
