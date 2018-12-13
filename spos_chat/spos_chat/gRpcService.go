// gRpcService
package main

import (
	//"container/list"
	//"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"io"
	"log"
	"spos_auth/datas"
	"spos_auth/umisc"
	proto "spos_proto/uni/voice" //
	//	"strconv"
	"github.com/dgrijalva/jwt-go"
	//	"github.com/dgrijalva/jwt-go/request"
	"strings"
	"time"
)

const (
	SecretKey = "unitone_spos"
)

type SposChat struct{}

// SayHello implements helloworld.GreeterServer
func (s *SposChat) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	fmt.Printf("%s:SposChat SayHello111 Text=%s\n", umisc.CreateDateString(), in.Name)
	return &proto.HelloReply{Message: "Hello111 " + in.Name + "; Unichat welcome!"}, nil
}

func (s *SposChat) Login(ctx context.Context, in *proto.ClientLoginReq) (*proto.ClientLoginRes, error) {

	fmt.Printf("%s:Login id=%s; password=%s;name=%s; mac=%s\n", umisc.CreateDateString(), in.Sid, in.Password, in.Name, in.Mac)
	var d_user datas.UserInfo
	d_user.GetUserById(in.Sid)
	fmt.Printf("LoginHandler: password=%s\n", d_user.EncryptedPassword)

	if len(d_user.EncryptedPassword) > 0 {
		Encryptor := umisc.BcryptEncryptorNew(&umisc.Config{})
		if err := Encryptor.Compare(d_user.EncryptedPassword, in.Password); err == nil {
			fmt.Println("LoginHandler:Password ok, update user work list")
			msg_ch := make(chan UniMsg, MSG_CHANNEL_NUM)
			user := UniWorkUser{UniUserInfo: UniUserInfo{Id: in.Sid, Gid: "000", Name: in.Name, Mac: in.Mac}, IsLogin: true, LoginTime: umisc.CreateDateString(), IsOnline: true, IsRecvOk: true, LastHeartTime: time.Now(), MsgCh: &msg_ch}
			UpdateUserToWorkList(user)

			//create token begin
			token := jwt.New(jwt.SigningMethodHS256)

			fmt.Printf("LoginHandler: user.UID=%s\n", in.Sid)

			basic_info := datas.Basic{UID: in.Sid}

			claims := basic_info.ToClaims()
			claims.ExpiresAt = time.Now().Add(time.Hour * time.Duration(10)).Unix()
			claims.IssuedAt = time.Now().Unix()

			token.Claims = claims

			//	tokenString, err := token.SignedString(signKey)
			tokenString, err := token.SignedString([]byte(SecretKey))

			fmt.Printf("LoginHandler:tokenString=%s; err=%v\n", tokenString, err)

			return &proto.ClientLoginRes{Ret: true, ResStr: GetLoginResInfo(), ResToken: tokenString}, nil
		} else {
			fmt.Printf("LoginHandler:Password error\n")
			return &proto.ClientLoginRes{Ret: false}, nil
		}
	} else {
		fmt.Printf("spos login account error\n")
		return &proto.ClientLoginRes{Ret: false}, nil
	}

}

func (s *SposChat) Logout(ctx context.Context, in *proto.ClientLogoutReq) (*proto.ClientLogoutRes, error) {

	fmt.Printf("%s:Logout id=%s; name=%s; mac=%s\n", umisc.CreateDateString(), in.Sid, in.Name, in.Mac)
	//	msg_ch := make(chan UniMsg, 10)
	//	user := UniWorkUser{UniUserInfo: UniUserInfo{Id: in.Sid, Gid: "000", Name: in.Name, Mac: in.Mac}, IsLogin: false, IsOnline: false, LogoutTime: umisc.CreateDateString(), MsgCh: &msg_ch}
	//	RemoveUserFromWorkList(user)
	user := GetWorkUserByIdAndMac(in.Sid, in.Mac)
	if user != nil {
		user.LogoutTime = umisc.CreateDateString()
		user.IsLogin = false
		user.IsOnline = false
		UpdateUserToWorkList(*user)

		return &proto.ClientLogoutRes{Ret: true}, nil
	} else {
		return &proto.ClientLogoutRes{Ret: false}, nil
	}

}

func (s *SposChat) HeartMsg(ctx context.Context, in *proto.ClientHeartReq) (*proto.ClientHeartRes, error) {

	fmt.Printf("%s:HeartReq id=%s; mac=%s\n", umisc.CreateDateString(), in.Sid, in.Mac)
	user := GetWorkUserByIdAndMac(in.Sid, in.Mac)

	if user != nil {
		user.LastHeartTime = time.Now()
		user.IsLogin = true
		user.IsOnline = true
		UpdateUserToWorkList(*user)
		return &proto.ClientHeartRes{Ret: true, Msg: GetHeartResInfo()}, nil
	}

	return &proto.ClientHeartRes{Ret: false, Msg: "user not work"}, nil

}

func (s *SposChat) RecvMsgStream(in *proto.RecvMsgReq, gs proto.Chat_RecvMsgStreamServer) error {

	fmt.Printf("%s:RecvMsgStream id=%s; name=%s; mac=%s\n", umisc.CreateDateString(), in.Sid, in.Name, in.Smac)

	ctx := gs.Context()

	user := GetWorkUserByIdAndMac(in.Sid, in.Smac)
	if user != nil {
		user.IsOnline = true
		user.IsRecvOk = true
		UpdateUserToWorkList(*user)
	}

	for user != nil {

		select {
		case <-ctx.Done():
			log.Println("RecvMsgStream::recv client end termination signal by context ")
			return ctx.Err()
		default:

			if !user.IsLogin || !user.IsOnline {
				return nil
			}

			log.Printf("RecvMsgStream::id=%sser::pull msg from msgCh channel,waiting...\n", user.Id)
			mch := <-*user.MsgCh

			log.Printf("RecvMsgStream::id=%sser::pull msg sid:%s; smac:%s; did:%s; dmac:%s\n",
				user.Id, mch.Sid, mch.Smac, mch.Did, mch.Dmac)
			log.Printf("RecvMsgStream::id=%sser::pull msg type:%d; sample_rate:%d; id:%s; no:%d; len:%d; status:%d\n",
				user.Id, mch.Mtype, mch.SampleRate, mch.MsgId, mch.MsgNum, mch.MsgLen, mch.MsgStatus)

			log.Printf("RecvMsgStream::id=%sser::pulled msg ok,send msg to client\n", user.Id)
			if err := gs.Send(&proto.RecvMsgRes{Mtype: proto.RecvMsgRes_MsgType(mch.Mtype), Sid: mch.Sid, Smac: mch.Smac, Did: mch.Did, Dmac: mch.Dmac, SampleRate: mch.SampleRate, MsgId: mch.MsgId, MsgNum: mch.MsgNum, MsgLen: mch.MsgLen, MsgStatus: mch.MsgStatus, Msg: mch.Msg}); err != nil {
				log.Printf("RecvMsgStream::id=%sser::send msg to client fail, err=%v\n", user.Id, err)
				return err
			}
		}
	}

	return nil
}

// SendMsgStream
func (s *SposChat) SendMsgStream(stream proto.Chat_SendMsgStreamServer) error {

	ctx := stream.Context()

	log.Println("SendMsgStream begin...")

	for {
		select {
		case <-ctx.Done():
			log.Println("SendMsgStream::recv client end termination signal by context ")
			return ctx.Err()
		default:
			// recv client msg
			in, err := stream.Recv()

			if err == io.EOF {
				log.Println("SendMsgStream::client send stream end")
				stream.SendAndClose(&proto.SendMsgRes{Ret: SEND_MSG_OK, Msg: "client send stream end"})
				return nil
			}

			if err != nil {
				log.Println("SendMsgStream::recv data error:", err)
				stream.SendAndClose(&proto.SendMsgRes{Ret: SEND_MSG_SERVER_RECV_ERR, Msg: "server recv data error"})
				return err
			}

			src_user := GetWorkUserByIdAndMac(in.Sid, in.Smac)
			if src_user == nil {
				log.Printf("SendMsgStream::id=%sser::not login or invaild user!\n", in.Sid)
				stream.SendAndClose(&proto.SendMsgRes{Ret: SEND_MSG_USER_LOGOUT, Msg: "not login or invaild user!"})
				return errors.New("not login or invaild user!")

			}

			log.Printf("SendMsgStream::id=%sser::recv msg sid:%s; mac=%s; did=%s; dmac=%s\n",
				src_user.Id, in.Sid, in.Smac, in.Did, in.Dmac)

			log.Printf("SendMsgStream::id=%sser::recv msg type:=%d; sample_rate:%d;id:%s; no:%d;len:%d; status=%d;\n",
				src_user.Id, in.Mtype, in.SampleRate, in.MsgId, in.MsgNum, in.MsgLen, in.MsgStatus)

			log.Printf("SendMsgStream::id=%sser::msg codeing...\n", src_user.Id)
			if !strings.Contains(in.Did, "0000") {
				dest_user := GetWorkUserById(in.Did)
				send_fail_num := 0
				if len(dest_user) > 0 {
					for index, us := range dest_user {

						if us.IsLogin && us.IsOnline {
							log.Printf("SendMsgStream::id=%sser::msg pushing to index=%d dest user msg channel\n", src_user.Id, index)

							m := UniMsg{Sid: in.Sid, Did: in.Did, Smac: in.Smac, Dmac: in.Dmac,
								Mtype: int32(in.Mtype), SampleRate: in.SampleRate, MsgId: in.MsgId, MsgNum: in.MsgNum, MsgLen: in.MsgLen, MsgStatus: in.MsgStatus, Msg: in.Msg}
							log.Printf("SendMsgStream::id=%sser::dest user id=%s msg channel len =%d\n", src_user.Id, in.Did, len(*us.MsgCh))
							if len(*us.MsgCh) >= MSG_CHANNEL_NUM {
								log.Printf("SendMsgStream::id=%sser::get dest user id=%s not online,not push, send msg fail\n", src_user.Id, in.Did)
								us.IsOnline = false
								us.IsRecvOk = false
								UpdateUserToWorkList(us)
								//stream.SendAndClose(&proto.SendMsgRes{Ret: SEMD_MSG_DEST_USER_OFFLINE, Msg: "dest user offline"})
								send_fail_num++
								continue
							} else {
								*us.MsgCh <- m
								log.Printf("SendMsgStream::id=%sser::msg push to dest user id=%s msg channel success\n", src_user.Id, us.Id)

							}

						} else {
							send_fail_num++
							log.Printf("SendMsgStream::id=%sser::get dest user id=%s not login or not online,not push\n", src_user.Id, in.Did)
							//stream.SendAndClose(&proto.SendMsgRes{Ret: SEND_MSG_DEST_USER_LOGOUT, Msg: "dest user not login"})
							continue
						}

					}

					if send_fail_num == len(dest_user) {
						stream.SendAndClose(&proto.SendMsgRes{Ret: SEMD_MSG_DEST_USER_OFFLINE, Msg: "dest user not login or offline"})
					}
				} else {
					log.Printf("SendMsgStream::id=%sser::get dest user id=%s fail,not push\n", src_user.Id, in.Did)
					stream.SendAndClose(&proto.SendMsgRes{Ret: SEND_MSG_NO_DEST_USER, Msg: "no dest user!"})
					//					return nil
				}

			}

		}
	}

	stream.SendAndClose(&proto.SendMsgRes{Ret: SEND_MSG_OK, Msg: "SendMsgStream:send msg stream Ok"})

	return nil

}

// BidStream
func (s *SposChat) BidStream(stream proto.Chat_BidStreamServer) error {

	//	ctx := stream.Context()

	//	var src_user *UniWorkUser

	//	log.Println("BidStream begin...")

	//	//go func() {

	//	for {
	//		select {
	//		case <-ctx.Done():
	//			log.Println("BidStream::recv client end termination signal by context ")
	//			return ctx.Err()
	//		default:
	//			// recv client msg
	//			in, err := stream.Recv()

	//			if err == io.EOF {
	//				log.Println("BidStream::client send stream end")
	//				return nil
	//			}

	//			if err != nil {
	//				log.Println("BidStream::recv data error:", err)
	//				return err
	//			}

	//			//
	//			//				id = in.Sid
	//			src_user = GetWorkUser(in.Sid)
	//			if src_user == nil {
	//				log.Printf("BidStream::id=%sser::not login or invaild user!\n", in.Sid)
	//				return errors.New("BidStream::not login or invaild user")

	//			}

	//			log.Printf("BidStream::id=%sser::recv msg type:%d\n", src_user.Id, in.Mtype)
	//			log.Printf("BidStream::id=%sser::recv msg sid:%s\n", src_user.Id, in.Sid)
	//			log.Printf("BidStream::id=%sser::recv msg did:%s\n", src_user.Id, in.Did)
	//			log.Printf("BidStream::id=%sser::recv msg sample_rate:%d\n", src_user.Id, in.SampleRate)
	//			log.Printf("BidStream::id=%sser::recv msg id:%s\n", src_user.Id, in.MsgId)
	//			log.Printf("BidStream::id=%sser::recv msg num:%d\n", src_user.Id, in.MsgNum)
	//			log.Printf("BidStream::id=%sser::recv msg len:%d\n", src_user.Id, in.MsgLen)
	//			log.Printf("BidStream::id=%sser::recv msg status:%d\n", src_user.Id, in.MsgStatus)
	//			//				log.Printf("id=%sser::recv msg bytes:%s\n", src_user.Id, string(in.Msg))

	//			switch string(in.Msg) {
	//			case "end\n":
	//				log.Println("BidStream::send end command")
	//				if err := stream.Send(&proto.ServerSendMsgToClient{Msg: []byte("recv end command")}); err != nil {
	//					return err
	//				}
	//				// recv end command, return nil
	//				return nil
	//			case "return datas\n":
	//				log.Println("BidStream::recv: return datas command")
	//				//
	//				for i := 0; i < 10; i++ {
	//					if err := stream.Send(&proto.ServerSendMsgToClient{Msg: []byte("datas #" + strconv.Itoa(i))}); err != nil {
	//						return err
	//					}
	//				}
	//			default:

	//				log.Printf("BidStream::id=%sser::msg codeing...\n", src_user.Id)
	//				if !strings.Contains(in.Did, "0000") {
	//					dest_user := GetWorkUser(in.Did)

	//					if dest_user != nil {
	//						log.Printf("BidStream::id=%sser::msg pushing to dest user msg channel\n", src_user.Id)

	//						m := UniMsg{Mtype: int32(in.Mtype), Sid: in.Sid, Did: in.Did, SampleRate: in.SampleRate, MsgId: in.MsgId, MsgNum: in.MsgNum, MsgLen: in.MsgLen, MsgStatus: in.MsgStatus, Msg: in.Msg}
	//						log.Printf("BidStream::id=%sser::dest user id=%s msg channel len =%d\n", src_user.Id, in.Did, len(*dest_user.MsgCh))

	//						*dest_user.MsgCh <- m
	//						log.Printf("BidStream::id=%sser::msg push to dest user id=%s msg channel end\n", src_user.Id, dest_user.Id)

	//					} else {
	//						log.Printf("BidStream::id=%sser::get dest user id=%s fail,not push\n", src_user.Id, in.Did)

	//					}

	//				}

	//				/*********************Test*******************/
	//				//					if err := stream.Send(&proto.ServerSendMsgToClient{Mtype: proto.ServerSendMsgToClient_MsgType(in.Mtype), Sid: in.Sid, Did: in.Did, SampleRate: in.SampleRate, Msg: in.Msg}); err != nil {
	//				//						return //err
	//				//					}
	//				//					if err := stream.Send(&proto.ServerSendMsgToClient{Msg: []byte("server return: " + string(in.Msg))}); err != nil {
	//				//						return //err
	//				//					}
	//			}
	//		}
	//	}

	return nil
}
