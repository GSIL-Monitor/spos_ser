/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	//	"bufio"
	"errors"
	//	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io"
	"log"
	"os"
	//"time"
	//	"encoding/base64"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/casbin/casbin"

	pb "spos_proto/uni"
	pb_auth "spos_proto/uni/auth"
	pb_hello "spos_proto/uni/hello"
	pb_voice "spos_proto/uni/voice"
)

const (
	//	address     = "18.191.92.76:61051"
	//address = "localhost:61055"
	//address = "196.168.99.100:61052"
	//	address     = "139.224.221.196:61052"
	address     = "192.168.10.82:32001" //ambassador
	defaultName = "hello world"

	/////////*************hello test***********//
	hello_test = false

	/////////*************misc test****************//
	db_test     = false
	spos_ssl    = true // false true
	token_test  = true
	casbin_test = false //true

	/////////**************lang test*****************//
	lang_test       = true
	lang_login_test = true
	lang_hello_test = true

	tts_test   = false //false
	trans_test = true

	/////////***************auth test******************//
	auth_test = false //true

	/////////***************chat test******************//
	voice_test       = false
	voice_hello_test = true
)

// customCredential
type customCredential struct {
	auth string
}

func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	fmt.Println("c.auth:" + c.auth)
	//	encodeString := base64.StdEncoding.EncodeToString([]byte(c.auth))
	//	fmt.Println("base64:" + encodeString)

	//	return map[string]string{
	//		"authorization": "Basic " + encodeString,
	//		//		"appkey": "i am key",
	//	}, nil

	return map[string]string{
		"authorization": "Bearer " + c.auth,
	}, nil
}

func (c customCredential) RequireTransportSecurity() bool {
	if spos_ssl {
		return true
	}

	return false
}

var (
	// Initialize gorm DB
	//	gormDB, _ = gorm.Open("sqlite3", "sample.db")
	gdb *gorm.DB
	//gormDB, _ = gorm.Open("mysql", "root:Unitone@/user_manager?charset=utf8&parseTime=True&loc=Local")

	// Initialize Auth with configuration
//	Auth = auth.New(&auth.Config{
//		DB: gormDB,
//	})
)

func OpenUnDb() error {

	fmt.Println("WORDPRESS_DB_HOST:" + os.Getenv("WORDPRESS_DB_HOST"))
	fmt.Println("WORDPRESS_DB_PORT:" + os.Getenv("WORDPRESS_DB_PORT"))
	fmt.Println("WORDPRESS_DB_PASSWORD:" + os.Getenv("WORDPRESS_DB_PASSWORD"))

	//	connect_str := "root:" + os.Getenv("WORDPRESS_DB_PASSWORD") + "@tcp(" + os.Getenv("WORDPRESS_DB_HOST") + ":" + os.Getenv("WORDPRESS_DB_PORT") + ")/uni_spos?charset=utf8"
	//	connect_str := "root:" + os.Getenv("WORDPRESS_DB_PASSWORD") + "@tcp(172.17.0.14:3306)/uni_spos?charset=utf8&parseTime=True&loc=Local"
	connect_str := "root:Unitone@2018" + "@tcp(192.168.20.107:3306)/uni_spos?charset=utf8&parseTime=True&loc=Local"
	fmt.Printf("spos auth connect_str=%s\n", connect_str)
	//	db, err := gorm.Open("mysql", "root:Unitone@/uni_spos?charset=utf8&parseTime=True&loc=Local")
	db, err := gorm.Open("mysql", connect_str)
	fmt.Printf("spos auth err=%v\n", err)

	//		db, err := gorm.Open("mysql", "root:Unitone@tcp(127.0.0.1:3306)/user_manager?charset=utf8&parseTime=True&loc=Local")
	//		db, err := gorm.Open("postgres", "host=10.110.130.48 port=5432 user=postgres dbname=uni_spos password=Unitone")
	if err != nil {
		//		panic("db connection failure")
		return errors.New("db open error")
	}

	fmt.Printf("spos auth db open ok!\n")
	gdb = db

	return nil
}

func main() {
	//	var conn *grpc.ClientConn
	var err error

	var c pb.LangClient
	var c_auth pb_auth.AuthClient
	var c_voice pb_voice.ChatClient
	var c_hello pb_hello.GreeterClient

	var opts []grpc.DialOption

	if db_test {
		if OpenUnDb() != nil {
			fmt.Printf("spos auth db open error!\n")
			panic("db connection failure")
		}
	}

	if casbin_test {

		// setup casbin auth rules
		e := &casbin.Enforcer{}
		e.InitWithFile("./basic_model.conf", "./basic_policy.csv")
		//		authEnforcer, err := casbin.NewEnforcerSafe("./basic_model.conf", "./basic_policy.csv")
		//		if err != nil {
		//			fmt.Printf("NewEnforcerSafe: err=%v\n", err)
		//			log.Fatal(err)
		//		}
		//		if e.Enforce(sub, obj, act) == true {
		//			// permit alice to read data1
		//		} else {
		//			// deny the request, show an error
		//		}

		return
	}

	Encryptor := BcryptEncryptorNew(&Config{})
	EncryptedPassword, err := Encryptor.Digest("unitone123")
	fmt.Println("greeter start...,EncryptedPassword=", EncryptedPassword)
	if err := Encryptor.Compare(EncryptedPassword, "unitone123"); err == nil {
		fmt.Println("password DefaultAuthorizeHandler password ok")
	} else {
		fmt.Println("password DefaultAuthorizeHandler password error")

	}

	address_new := os.Getenv("ADDRESS")
	if len(address_new) == 0 {
		address_new = address
	}

	fmt.Println("greeter start...,port=", address_new)

	if spos_ssl {

		fmt.Println("greeter with ssl")
		// TLS
		creds, err := credentials.NewClientTLSFromFile("./gubstech.crt", "gubstech.com")
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
		//	conn, err = grpc.Dial(address_new, grpc.WithTransportCredentials(creds))

	} else {
		fmt.Println("greeter with no ssl")
		// Set up a connection to the server.
		opts = append(opts, grpc.WithInsecure())
		//	conn, err = grpc.Dial(address_new, grpc.WithInsecure())

	}

	if token_test {
		fmt.Println("greeter with token")
		fmt.Println("SPOS_AUTH:" + os.Getenv("SPOS_AUTH"))
		opts = append(opts, grpc.WithPerRPCCredentials(&customCredential{auth: os.Getenv("SPOS_AUTH")}))
	}

	conn, err := grpc.Dial(address_new, opts...)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	//	c := pb.NewGreeterClient(conn)

	name := defaultName

	if auth_test {
		fmt.Println("\ngreeter with auth_test")
		c_auth = pb_auth.NewAuthClient(conn)

		fmt.Println("\ngreeter auth_test")
		login_resp, err := c_auth.LoginRequest(context.Background(), &pb_auth.LoginReq{Username: "123", Passwd: "456"})
		if err != nil {
			//log.Fatalf("could not Login: %v", err)
			log.Printf("could not Login: %v", err)
		}

		for {
			c, err := login_resp.Recv()
			if err != nil {
				if err == io.EOF {
					//	log.Printf("Transfer of %d bytes successful", len(blob))
					log.Printf("Transfer bytes successful final")
					break
				}

				panic(err)
			}

			fmt.Println("c.Ret=", c.Ret)
			fmt.Println("c.Token=", c.Token)

		}

	}

	if voice_test {
		c_voice = pb_voice.NewChatClient(conn)
		fmt.Println("\ngreeter with voice_test")

		if voice_hello_test {
			fmt.Println("\nvoice hello_test")
			//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			//	defer cancel()
			//			for i := 0; i < 5; i++ {

			//			fmt.Println("test=", i)

			//				go func() {
			//			fmt.Println("go test=", i)

			//			r, err := c_voice.SayHello(context.Background(), &pb_voice.HelloRequest{Name: name})
			//			if err != nil {
			//				log.Fatalf("could not voice greet: %v", err)
			//			}
			//			log.Printf("voice Greeting: %s\n", r.Message)

			r, err := c_voice.Login(context.Background(), &pb_voice.ClientLoginReq{Sid: "U1542097798", Password: "1234567"})
			if err != nil {
				log.Fatalf("could not voice login: %v", err)
			}
			log.Printf("voice login token: %s\n", r.ResToken)

			//				}()
			//			}

		}

	}

	if lang_test {
		//		c = pb.NewSposClient(conn)
		c = pb.NewLangClient(conn)
		fmt.Println("greeter with lang_test")

		if lang_login_test {
			fmt.Println("\nlang login test")
			r, err := c.Login(context.Background(), &pb.ClientLoginReq{Sid: "U1542097798", Password: "1234567"})
			if err != nil {
				log.Fatalf("could not voice login: %v", err)
			}
			log.Printf("voice login token: %s\n", r.ResToken)
		}

		if lang_hello_test {
			fmt.Println("\nlang hello_test")
			//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			//	defer cancel()
			r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
			if err != nil {
				log.Fatalf("could not lang greet: %v", err)
			}
			log.Printf("Greeting: %s\n", r.Message)
		}

		if trans_test {
			fmt.Println("\nlang trans_test")
			r, err := c.TransRequest(context.Background(), &pb.TransReq{SLn: "en", TLn: "ru", Text: name})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			log.Printf("Translate text: %s", r.TransText)

		}
	}

	//	if tts_test {
	//		fmt.Println("greeter tts_test")
	//		audio_stream, err := c.TtsRequest(context.Background(), &pb.TtsReq{Text: name, Lang: "en-US", Gender: 2, AudioCode: 2})
	//		if err != nil {
	//			log.Fatalf("could not greet: %v", err)
	//		}

	//		if PathExists("spos.mp3") {
	//			err := os.Remove("spos.mp3")
	//			if err != nil {
	//				fmt.Println("file remove Error!")
	//				fmt.Printf("%s\n", err)
	//			} else {
	//				fmt.Println("file remove ok!")
	//			}
	//		}
	//		fo, err := os.OpenFile("spos.mp3", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	//		if err != nil {
	//			panic("open file failed")
	//		}
	//		defer fo.Close()

	//		for {
	//			c, err := audio_stream.Recv()
	//			if err != nil {
	//				if err == io.EOF {
	//					//	log.Printf("Transfer of %d bytes successful", len(blob))
	//					log.Printf("Transfer bytes successful final")
	//					break
	//				}

	//				panic(err)
	//			}

	//			if _, err := fo.Write(c.AudioData[0:]); err != nil {
	//				panic(err)
	//			}
	//		}
	//	}

	if hello_test {
		fmt.Println("\ngreeter hello_test")
		//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		//	defer cancel()
		r, err := c_hello.SayHello(context.Background(), &pb_hello.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s\n", r.Message)
	}

	//	log.Printf("audio data len: %d", len(r))

	//	if voice_test {
	//		var sid, did, msg string
	//		flag.StringVar(&sid, "sid", "U1536825106", "source id")
	//		flag.StringVar(&did, "did", "U1536825121", "dest id")
	//		flag.StringVar(&msg, "msg", "hello!", "send msg")
	//		flag.Parse()
	//		log.Println("source id:", sid)
	//		log.Println("dest id:", did)
	//		log.Println("send msg:", msg)

	//		// get context
	//		ctx := context.Background()

	//		r, err := c_voice.Login(ctx, &pb_voice.ClientLoginReq{Sid: sid, Name: "patli2", Mac: "22.11.33.44.55.66"})
	//		if err != nil {
	//			log.Fatalf("could not greet: %v", err)
	//		}
	//		log.Printf("Greeting: %v\n", r.Ret)

	//		// create stream
	//		stream, err := c_voice.BidStream(ctx)
	//		if err != nil {
	//			log.Printf("create stream fail: [%v]\n", err)
	//		}

	//		if err := stream.Send(&pb_voice.ClientSendMsgToServer{Mtype: pb_voice.ClientSendMsgToServer_TEXT, Sid: sid, Did: "0000", Msg: []byte("Hi!I am comming!")}); err != nil {
	//			return
	//		}

	//		// start goroutine input client command
	//		go func() {
	//			log.Println("please input command...")
	//			in_com := bufio.NewReader(os.Stdin)
	//			for {
	//				// get command, enter end
	//				in_str, _ := in_com.ReadString('\n')
	//				// send command to server
	//				if err := stream.Send(&pb_voice.ClientSendMsgToServer{Mtype: pb_voice.ClientSendMsgToServer_TEXT, Sid: sid, Did: did, Msg: []byte(in_str)}); err != nil {
	//					return
	//				}
	//			}
	//		}()
	//		for {
	//			// recv data stream from server
	//			res, err := stream.Recv()
	//			if err == io.EOF {
	//				log.Println("recv server end signal")
	//				break // end
	//			}
	//			if err != nil {
	//				// TODO: recv error
	//				log.Println("recv data error:", err)
	//			}
	//			// print server msg
	//			log.Printf("from id=%s msg id=%s; msg num=%d; msg status=%d; datas: %s; ", res.Sid, res.MsgId, res.MsgNum, res.MsgStatus, string(res.Msg))
	//		}
	//	}
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
