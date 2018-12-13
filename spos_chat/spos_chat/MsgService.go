// MsgService
package main

import (
//	"fmt"
)

const (
	SEND_MSG_OK                = 0
	SEND_MSG_USER_LOGOUT       = 1
	SEND_MSG_NO_DEST_USER      = 2
	SEND_MSG_DEST_USER_LOGOUT  = 3
	SEND_MSG_SERVER_RECV_ERR   = 4
	SEMD_MSG_DEST_USER_OFFLINE = 5
)

const (
	MSG_CHANNEL_NUM = 100
)

// Streamer
type UniMsg struct {
	Mtype      int32 //0--text; 1--voice; 2--image; 3--video
	Sid        string
	Smac       string
	Did        string
	Dmac       string
	SampleRate int32
	MsgId      string
	MsgNum     int32
	MsgLen     int32
	MsgStatus  int32
	Msg        []byte
}
