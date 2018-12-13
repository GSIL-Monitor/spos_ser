// LangGrpcServer
package main

import (
	speech "cloud.google.com/go/speech/apiv1"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"io"
	"log"
	//	"net"
	//	"os"
	"strings"

	"spos_lang/trans"
	"spos_lang/tts"
	"spos_lang/ulog"
	pb "spos_proto/uni"

	"github.com/dgrijalva/jwt-go"
	"spos_auth/datas"
	"spos_auth/umisc"
	"time"
)

const (
	SecretKey = "unitone_spos"
)

type lang_server struct{}

// SayHello implements helloworld.GreeterServer
func (s *lang_server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Printf("%s:Spos Lang SayHello222 Text=%s\n", ulog.CreateDateString(), in.Name)
	return &pb.HelloReply{Message: "Hello222 " + in.Name}, nil
}

func (s *lang_server) Login(ctx context.Context, in *pb.ClientLoginReq) (*pb.ClientLoginRes, error) {

	fmt.Printf("%s:Login id=%s; password=%s;name=%s; mac=%s\n", umisc.CreateDateString(), in.Sid, in.Password, in.Name, in.Mac)
	var d_user datas.UserInfo
	d_user.GetUserById(in.Sid)
	fmt.Printf("LoginHandler: password=%s\n", d_user.EncryptedPassword)

	if len(d_user.EncryptedPassword) > 0 {
		Encryptor := umisc.BcryptEncryptorNew(&umisc.Config{})
		if err := Encryptor.Compare(d_user.EncryptedPassword, in.Password); err == nil {
			fmt.Println("LoginHandler:Password ok, update user work list")
			//			msg_ch := make(chan UniMsg, MSG_CHANNEL_NUM)
			user := UniWorkUser{UniUserInfo: UniUserInfo{Id: in.Sid, Gid: "000", Name: in.Name, Mac: in.Mac}, IsLogin: true, LoginTime: umisc.CreateDateString(), IsOnline: true, IsRecvOk: true, LastHeartTime: time.Now()}
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

			return &pb.ClientLoginRes{Ret: true, ResStr: GetLoginResInfo(), ResToken: tokenString}, nil
		} else {
			fmt.Printf("LoginHandler:Password error\n")
			return &pb.ClientLoginRes{Ret: false}, nil
		}
	} else {
		fmt.Printf("spos login account error\n")
		return &pb.ClientLoginRes{Ret: false}, nil
	}

}

func (s *lang_server) Logout(ctx context.Context, in *pb.ClientLogoutReq) (*pb.ClientLogoutRes, error) {

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

		return &pb.ClientLogoutRes{Ret: true}, nil
	} else {
		return &pb.ClientLogoutRes{Ret: false}, nil
	}

}

func (s *lang_server) HeartMsg(ctx context.Context, in *pb.ClientHeartReq) (*pb.ClientHeartRes, error) {

	fmt.Printf("%s:HeartReq id=%s; mac=%s\n", umisc.CreateDateString(), in.Sid, in.Mac)
	user := GetWorkUserByIdAndMac(in.Sid, in.Mac)

	if user != nil {
		user.LastHeartTime = time.Now()
		user.IsLogin = true
		user.IsOnline = true
		UpdateUserToWorkList(*user)
		return &pb.ClientHeartRes{Ret: true, Msg: GetHeartResInfo()}, nil
	}

	return &pb.ClientHeartRes{Ret: false, Msg: "user not work"}, nil

}

func (s *lang_server) TtsRequest(in *pb.TtsReq, gs pb.Lang_TtsRequestServer) error {

	fmt.Printf("%s:TtsRequest Text=%s\n", ulog.CreateDateString(), in.Text)

	audio_data, _ := tts.SynthesizeTextToStream(in.Text, in.Lang, in.Gender, in.AudioCode)

	fmt.Printf("%s:TtsRepley audio_data len=%d\n", ulog.CreateDateString(), len(audio_data))

	chnk := &pb.TtsReply{}
	chnk.AudioData = audio_data[:len(audio_data)]
	if err := gs.Send(chnk); err != nil {
		fmt.Println("TtsRequest send err=", err)
		return err
	}
	fmt.Printf("%s:Tts audio data send ok\n", ulog.CreateDateString())

	return nil
}

func (s *lang_server) TransRequest(ctx context.Context, in *pb.TransReq) (*pb.TransReply, error) {

	fmt.Printf("%s:TransRequest Text=%s\n", ulog.CreateDateString(), in.Text)

	trans_str, err := trans.TransRequest(in.SLn, in.TLn, in.Text)
	if err != nil {
		return &pb.TransReply{TransText: ""}, err
	}

	fmt.Printf("%s:TransRepley Text=%s\n", ulog.CreateDateString(), trans_str)

	return &pb.TransReply{TransText: trans_str}, nil
}

var ret *pb.StreamingRecognizeResponse = new(pb.StreamingRecognizeResponse)

func (s *lang_server) StreamingRecognize(gs pb.Lang_StreamingRecognizeServer) error {
	//	fmt.Printf("%s:StreamingRecognize0\n", ulog.CreateDateSting())

	ctx := context.Background()

	// [START speech_streaming_mic_recognize]
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	stream, err := client.StreamingRecognize(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			//			fmt.Printf("%s:StreamingRecognize1\n", ulog.CreateDateSting())
			in, err := gs.Recv()

			//			fmt.Printf("%s:StreamingRecognize11\n", ulog.CreateDateSting())

			if err == io.EOF {
				fmt.Println("is Eof!")
				// Nothing else to pipe, close the stream.
				if err := stream.CloseSend(); err != nil {
					log.Fatalf("Could not close stream: %v", err)
					fmt.Printf("Could not close stream: %v\n", err)
				}
				return
			}

			if err != nil {
				fmt.Printf("failed to recv: %v\n", err)
				//				return
				continue
			}

			switch x := in.StreamingRequest.(type) {
			case *(pb.StreamingRecognizeRequest_StreamingConfig):
				fmt.Printf("%s:Stt Config encode= %d, SampleRate=%d, Language=%s\n", ulog.CreateDateString(), x.StreamingConfig.Config.Encoding, x.StreamingConfig.Config.SampleRateHertz, x.StreamingConfig.Config.LanguageCode)
				// Send the initial configuration message.
				if err := stream.Send(&speechpb.StreamingRecognizeRequest{
					StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
						StreamingConfig: &speechpb.StreamingRecognitionConfig{
							Config: &speechpb.RecognitionConfig{
								Encoding:        speechpb.RecognitionConfig_AudioEncoding(x.StreamingConfig.Config.Encoding),
								SampleRateHertz: x.StreamingConfig.Config.SampleRateHertz,
								LanguageCode:    x.StreamingConfig.Config.LanguageCode,
							},
						},
					},
				}); err != nil {
					log.Fatal(err)
					fmt.Printf("stream send config err: %v\n", err)
				}
			case *(pb.StreamingRecognizeRequest_AudioContent):
				if err = stream.Send(&speechpb.StreamingRecognizeRequest{
					StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
						AudioContent: x.AudioContent[0:],
					},
				}); err != nil {
					log.Printf("Could not send audio: %v", err)
					fmt.Printf("Could not send audio err: %v\n", err)
				}
			case nil:
			default:
				//				return fmt.Errorf("StreamingRecognizeRequest.StreamingRequest has unexpected type %T", x)
				return
			}
		}
	}()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cannot stream results: %v", err)
		}
		if err := resp.Error; err != nil {
			// Workaround while the API doesn't give a more informative error.
			if err.Code == 3 || err.Code == 11 {
				log.Print("WARNING: Speech recognition request exceeded limit of 60 seconds.")
			}
			log.Fatalf("Could not recognize: %v", err)
		}

		//json for struct copy
		aj, _ := json.Marshal(resp)
		_ = json.Unmarshal(aj, ret)

		for _, result := range resp.Results {
			//		fmt.Printf("%s:Stt Result: %+v\n", ulog.CreateDateSting(), result)
			//			fmt.Printf("%s:Result Alternatives len= %d\n", ulog.CreateDateSting(), len(result.Alternatives))
			for _, alt := range result.Alternatives {
				fmt.Printf("%s:Stt text=\"%v\" (confidence=%3f)\n", ulog.CreateDateString(), alt.Transcript, alt.Confidence)

			}

		}

		gs.Send(ret)
	}

	return nil
}

var trans_str string
var trans_source string
var trans_target string
var tts_lang string
var tts_gender int32
var tts_audio_code int32

func (s *lang_server) StreamingRecTransTTs(gs pb.Lang_StreamingRecTransTTsServer) error {

	//	fmt.Printf("%s:StreamingRecTransTTs0\n", ulog.CreateDateSting())

	ctx := context.Background()

	// [START speech_streaming_mic_recognize]
	client, err := speech.NewClient(ctx)
	if err != nil {
		//	log.Fatal(err)
		fmt.Printf("%s: speech.NewClient(ctx) fail! %v\n\n", ulog.CreateDateString(), err)
		return nil
	}
	stream, err := client.StreamingRecognize(ctx)
	if err != nil {
		fmt.Printf("%s: client.StreamingRecognize(ctx) fail! %v\n\n", ulog.CreateDateString(), err)
		return nil
	}

	go func() {
		for {
			//			fmt.Printf("%s:StreamingRecognize1\n", ulog.CreateDateSting())
			in, err := gs.Recv()

			//			fmt.Printf("%s:StreamingRecognize11\n", ulog.CreateDateSting())

			if err == io.EOF {
				fmt.Printf("%s:speak is Eof!\n", ulog.CreateDateString())
				// Nothing else to pipe, close the stream.
				if err := stream.CloseSend(); err != nil {
					//log.Fatalf("Could not close stream: %v", err)
					fmt.Printf("%s:Could not close stream: %v\n", ulog.CreateDateString(), err)
				}
				return
			}

			if err != nil {
				fmt.Printf("%s:Stt failed to recv: %v\n\n", ulog.CreateDateString(), err)
				return
				//continue
			}

			switch x := in.StreamingRequest.(type) {
			case *(pb.StreamingRecgTransTtsRequest_StreamingConfig):
				fmt.Printf("%s:Stt Config encode= %d, SampleRate=%d, Language=%s\n", ulog.CreateDateString(), x.StreamingConfig.Config.Encoding, x.StreamingConfig.Config.SampleRateHertz, x.StreamingConfig.Config.LanguageCode)
				// Send the initial configuration message.

				if err := stream.Send(&speechpb.StreamingRecognizeRequest{
					StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
						StreamingConfig: &speechpb.StreamingRecognitionConfig{
							Config: &speechpb.RecognitionConfig{
								Encoding:        speechpb.RecognitionConfig_AudioEncoding(x.StreamingConfig.Config.Encoding),
								SampleRateHertz: x.StreamingConfig.Config.SampleRateHertz,
								LanguageCode:    x.StreamingConfig.Config.LanguageCode,
							},
						},
					},
				}); err != nil {
					fmt.Printf("%s:stream send config to google cloud, err: %v\n", ulog.CreateDateString(), err)
					return
					//log.Fatal(err)

				}

				trans_source = x.StreamingConfig.TransSoure
				trans_target = x.StreamingConfig.TransTarget
				tts_lang = x.StreamingConfig.TtsLn
				tts_gender = x.StreamingConfig.TtsGender
				tts_audio_code = x.StreamingConfig.TtsAudioCode

				fmt.Printf("%s:tr_s= %s, tr_t=%s, tts_ln=%s,tts_g=%d, tts_ac=%d\n", ulog.CreateDateString(), trans_source, trans_target, tts_lang, tts_gender, tts_audio_code)

			case *(pb.StreamingRecgTransTtsRequest_AudioContent):
				if err = stream.Send(&speechpb.StreamingRecognizeRequest{
					StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
						AudioContent: x.AudioContent[0:],
					},
				}); err != nil {
					//		log.Printf("Could not send audio: %v", err)
					fmt.Printf("%s:Could not send audio to google cloud, err: %v\n", ulog.CreateDateString(), err)
				}
			case nil:
			default:
				//				return fmt.Errorf("StreamingRecognizeRequest.StreamingRequest has unexpected type %T", x)
				return
			}
		}
	}()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("%s:Cannot stream results: %v\n\n", ulog.CreateDateString(), err)
			break
		}
		if err := resp.Error; err != nil {
			// Workaround while the API doesn't give a more informative error.
			if err.Code == 3 || err.Code == 11 {
				log.Print("WARNING: Speech recognition request exceeded limit of 60 seconds.")
			}
			fmt.Printf("%s:Could not recognize: %v\n\n", ulog.CreateDateString(), err)
			break
		}

		aj, _ := json.Marshal(resp)
		_ = json.Unmarshal(aj, ret)

		for _, result := range resp.Results {

			for _, alt := range result.Alternatives {
				fmt.Printf("%s:Stt text=\"%v\" (confidence=%3f)\n", ulog.CreateDateString(), alt.Transcript, alt.Confidence)
				trans_str = alt.Transcript

			}

		}

		if err := gs.Send(&pb.StreamingRecTranTTsResponse{StreamingResponse: &pb.StreamingRecTranTTsResponse_Recg{ret}}); err != nil {
			fmt.Printf("%s:stt text send fail! err=%v\n\n", ulog.CreateDateString(), err)
			return err
		}
	}

	//////************************trans begin*********************/////////////////
	//	fmt.Printf("%s:TransRequest Text=%s\n", ulog.CreateDateSting(), trans_str)

	if len(trans_str) == 0 {
		fmt.Printf("%s:Req Trans text is null\n\n", ulog.CreateDateString())
		return nil
	}

	tts_str, err := trans.TransRequest(trans_source, trans_target, trans_str)
	trans_str = ""

	if err != nil {
		//return &pb.TransReply{TransText: ""}, err
		gs.Send(&pb.StreamingRecTranTTsResponse{StreamingResponse: &pb.StreamingRecTranTTsResponse_TransText{[]byte("trans error!")}})
	}

	if len(tts_str) > 0 {
		fmt.Printf("%s:TransRepley Text=%s\n", ulog.CreateDateString(), tts_str)

		if err := gs.Send(&pb.StreamingRecTranTTsResponse{StreamingResponse: &pb.StreamingRecTranTTsResponse_TransText{[]byte(tts_str)}}); err != nil {
			fmt.Printf("%s:Trans text send fail! err=%v\n\n", ulog.CreateDateString(), err)
			return err
		}
	} else {
		fmt.Printf("%s:Trans text is null\n", ulog.CreateDateString())
	}

	/////*********************trans end***************************/////////////////

	//////************************tts begin*********************/////////////////
	if !strings.Contains(trans_target, "zh-CN") {
		//		fmt.Printf("%s:TtsRequest Text=%s\n", ulog.CreateDateSting(), tts_str)
		if len(tts_str) > 0 {
			audio_data, _ := tts.SynthesizeTextToStream(tts_str, tts_lang, tts_gender, tts_audio_code)

			fmt.Printf("%s:TtsRepley audio_data len=%d\n", ulog.CreateDateString(), len(audio_data))

			if err := gs.Send(&pb.StreamingRecTranTTsResponse{StreamingResponse: &pb.StreamingRecTranTTsResponse_AudioData{audio_data}}); err != nil {
				fmt.Printf("%s:Tts audio data send fail! err=%v\n\n", ulog.CreateDateString(), err)
				return err
			}

			fmt.Printf("%s:Tts audio data send ok\n\n", ulog.CreateDateString())
		} else {
			fmt.Printf("%s:Tts text is null\n\n", ulog.CreateDateString())
		}

	} else {
		fmt.Printf("%s:Tts lang is zh-CN! TTS no support!\n\n", ulog.CreateDateString())
	}

	/////*********************tts end***************************/////////////////

	return nil
}
