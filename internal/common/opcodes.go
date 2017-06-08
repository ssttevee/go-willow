package common

type OpCode uint16

const (
	RecvOpCode_Unknown_0067 OpCode = 0x0067 // 1 byte: locale + uint16: version + uint16: patch + 1 byte: probably null // maybe login acknowledgement
	RecvOpCode_ListServers  OpCode = 0x0068 // 1 byte: probably boolean
	RecvOpCode_Unknown_006A OpCode = 0x006A // 1 byte: probably boolean + maplestring: auth token + 27 bytes: unknown
	RecvOpCode_Unknown_0072 OpCode = 0x0072 // length: 0
	RecvOpCode_Unknown_0086 OpCode = 0x0086 // length: 4  "1A 36 00 00" // maybe pong
	RecvOpCode_Unknown_0093 OpCode = 0x0093 // length: 4  "EC 6C FE 83" // maybe pong
	RecvOpCode_Unknown_0098 OpCode = 0x0098 // length: 4  "A1 AA 8E 1F" // maybe part of handshake?
	RecvOpCode_Unknown_009A OpCode = 0x009A // uint32: startup time?
	RecvOpCode_Unknown_00A2 OpCode = 0x00A2 // length: 2  "2D 00"
	RecvOpCode_Unknown_00B1 OpCode = 0x00B1 // length: 0
	RecvOpCode_Unknown_02C7 OpCode = 0x02C7 // length: 0
)

const (
	SendOpCode_WorldChannelList            OpCode = 0x0001 // 1 byte: null + maplestring: world name + 1 byte: (null or message type) + maplestring: world message + 3 bytes: "64 00 64" + 1 byte: "14" | "0F"
	SendOpCode_Unknown_0002                OpCode = 0x0002 // length: 4  "FD 00 00 00"
	SendOpCode_RecommendedWorldDescription OpCode = 0x0003 // 1 byte: world + bytes: unknown + maplestring: description
	SendOpCode_Ping                        OpCode = 0x0012 // 0 bytes
	SendOpCode_Unknown_0017                OpCode = 0x0017 // length: 5  "00 00 00 00 00"
	SendOpCode_Unknown_0018                OpCode = 0x0018 // length: 0
	SendOpCode_Unknown_0026                OpCode = 0x0026 // length: 1  "00"
	SendOpCode_Unknown_0032                OpCode = 0x0032 // length: 1  "00"
	SendOpCode_Unknown_0036                OpCode = 0x0036 // 5 bytes: "70 96 0B CE BB" + 3 bytes: "DF D2 01" + maplestring: "MapLogin1" | "MapLogin2"
	SendOpCode_Unknown_0039                OpCode = 0x0039 // length: 1  "00"
	SendOpCode_Unknown_7777                OpCode = 0x7777
	SendOpCode_Unknown_7778                OpCode = 0x7778
)
