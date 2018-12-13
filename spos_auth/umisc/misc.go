// misc
package umisc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

func CreateDateString() string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second())

}

func GetTimeSub(t string) {

	//	time1 := "2015-03-20 08:50:29"
	//	time2 := "2015-03-21 09:04:25"

	//	t1, err := time.Parse("2006-01-02 15:04:05", time1)
	//	t2, err := time.Parse("2006-01-02 15:04:05", time2)
	//	if err == nil && t1.Before(t2) {

	//		fmt.Println("true")
	//	}

	t3, _ := time.Parse("2006-01-02 15:04:05", t)

	fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n",
		t3.Year(),
		t3.Month(),
		t3.Day(),
		t3.Hour(),
		t3.Minute(),
		t3.Second())

	fmt.Printf("time du=%f\n", time.Now().Sub(t3).Hours())

	//	now := time.Now()

	end_time := time.Now()
	var dur_time time.Duration = end_time.Sub(t3)
	var elapsed_min float64 = dur_time.Minutes()
	var elapsed_sec float64 = dur_time.Seconds()
	var elapsed_nano int64 = dur_time.Nanoseconds()
	fmt.Printf("elasped %f minutes or \nelapsed %f seconds or \nelapsed %d nanoseconds\n",
		elapsed_min, elapsed_sec, elapsed_nano)

}

func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n+1:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//整形转换成字节
func Int16ToBytesBigEndian(n int16) []byte {

	x := int16(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//整形转换成字节
func Int16ToBytesLitEndian(n int16) []byte {

	x := int16(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}

//4字节转换成整形
func Bytes4ToIntOfBigEndian(b []byte) int32 {

	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int32(x)
}

//4字节转换成整形
func Bytes4ToIntOfLitEndian(b []byte) int32 {

	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)

	return int32(x)
}

//4字节转换成整形
func Bytes8ToIntOfBigEndian(b []byte) int64 {

	bytesBuffer := bytes.NewBuffer(b)

	var x int64
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int64(x)
}

//4字节转换成整形
func Bytes8ToIntOfLitEndian(b []byte) int64 {

	bytesBuffer := bytes.NewBuffer(b)

	var x int64
	binary.Read(bytesBuffer, binary.LittleEndian, &x)

	return int64(x)
}

//2字节转换成整形
func Bytes2ToIntBigEndian(b []byte) int16 {

	bytesBuffer := bytes.NewBuffer(b)

	var x int16
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int16(x)
}

//2字节转换成整形
func Bytes2ToIntLitEndian(b []byte) int16 {

	bytesBuffer := bytes.NewBuffer(b)

	var x int16
	binary.Read(bytesBuffer, binary.LittleEndian, &x)

	return int16(x)
}

//1字节转换成整形
func Bytes1ToInt(b []byte) int8 {

	bytesBuffer := bytes.NewBuffer(b)

	var x int8
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int8(x)
}

//整形转换成字节
func IntToBytesBigEndian(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//整形转换成字节
func IntToBytesLittleEndian(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}

//4字节转换成整形
func Bytes4ToInt(b []byte) int {

	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

func SendMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}
