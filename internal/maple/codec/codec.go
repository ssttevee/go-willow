package codec

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"

	"golang.ssttevee.com/willow/internal/common"
)

type packetCodec struct {
	block  cipher.Block
	ivRecv IV
	ivSend IV

	serverVersion uint16
}

func (codec *packetCodec) ReadPacket(reader io.Reader) (common.ReadablePacket, error) {
	// read packet header (packet length)
	header := make([]byte, 4)
	_, err := reader.Read(header)
	if err != nil {
		return nil, err
	}

	if version := binary.LittleEndian.Uint16([]byte{
		header[0] ^ codec.ivRecv[2],
		header[1] ^ codec.ivRecv[3],
	}); version != codec.serverVersion {
		return nil, fmt.Errorf("Mismatched Versions: %d (client) != %d (server)", version, codec.serverVersion)
	}

	packetLength := binary.LittleEndian.Uint16([]byte{
		header[2] ^ header[0],
		header[3] ^ header[1],
	})

	// read packet
	packet := make([]byte, packetLength)
	if _, err := io.ReadFull(reader, packet); err != nil {
		return nil, err
	}

	decrypted := codec.Decode(packet)

	return newReceivedPacket(decrypted), nil
}

func (codec *packetCodec) WritePacket(writer io.Writer, packet common.Packet) error {
	data := packet.Bytes()

	header1 := make([]byte, 2)
	binary.LittleEndian.PutUint16(header1, 0xFFFF-codec.serverVersion)
	header1[0] ^= codec.ivSend[2]
	header1[1] ^= codec.ivSend[3]

	header2 := make([]byte, 2)
	binary.LittleEndian.PutUint16(header2, uint16(len(data)))
	header2[0] ^= header1[0]
	header2[1] ^= header1[1]

	var buf bytes.Buffer
	if _, err := buf.Write(header1); err != nil {
		return err
	}

	if _, err := buf.Write(header2); err != nil {
		return err
	}

	encrypted := codec.Encode(data)

	if _, err := buf.Write(encrypted); err != nil {
		return err
	}

	if _, err := buf.WriteTo(writer); err != nil {
		return err
	}

	return nil
}

func (codec *packetCodec) Encode(data []byte) []byte {
	return codec.crypt(codec.ivSend, data)
}

func (codec *packetCodec) Decode(data []byte) []byte {
	return codec.crypt(codec.ivRecv, data)
}

func (codec *packetCodec) ReceiverIV() []byte {
	return []byte(codec.ivRecv)[:]
}

func (codec *packetCodec) SenderIV() []byte {
	return []byte(codec.ivSend)[:]
}

func New(key []byte, serverVersion uint16) (*packetCodec, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//generate IVs
	ivs := make([]byte, 8)
	rand.Read(ivs)

	return &packetCodec{
		block:  block,
		ivRecv: IV(ivs[:4]),
		ivSend: IV(ivs[4:]),

		serverVersion: serverVersion,
	}, nil
}

func (codec *packetCodec) crypt(iv IV, data []byte) []byte {
	out := make([]byte, len(data))

	remaining := len(data)
	llength := 0x5B0
	start := 0
	for remaining > 0 {
		morph_key := []byte{
			iv[0], iv[1], iv[2], iv[3],
			iv[0], iv[1], iv[2], iv[3],
			iv[0], iv[1], iv[2], iv[3],
			iv[0], iv[1], iv[2], iv[3],
		}

		if remaining < llength {
			llength = remaining
		}

		for x := start; x < start+llength; x++ {
			if (x-start)%16 == 0 {
				codec.block.Encrypt(morph_key, morph_key)
			}

			out[x] = data[x] ^ morph_key[(x-start)%16]
		}

		start += llength
		remaining -= llength
		llength = 0x5B4
	}

	cipher.NewOFB(codec.block, []byte{
		iv[0], iv[1], iv[2], iv[3],
		iv[0], iv[1], iv[2], iv[3],
		iv[0], iv[1], iv[2], iv[3],
		iv[0], iv[1], iv[2], iv[3],
	}).XORKeyStream(out, data)

	iv.Update()

	return out
}
