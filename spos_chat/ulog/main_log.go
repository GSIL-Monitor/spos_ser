// main_log
package ulog

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func MakeMultDir(dir string) {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		fmt.Printf("Error MakeMultDir %s\n", err)
	} else {
		//		fmt.Print("MakeMultDir Create Directory OK!\n")
	}
}

func CreateDateString() string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second())
}

func CreateDateWithNanoString() string {
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d_%d",
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second(),
		time.Now().Nanosecond())
}

var Ul *Logger

func CreateLog() {

	MakeMultDir("usc_log")
	lf, err := Open("./usc_log/log.txt", 1024, 100*1024*1024, 64)
	if err != nil {
		os.Exit(1)
	}
	//	defer lf.Close()

	//log file
	Ul = New(lf, "", log.LstdFlags|log.Lshortfile, LevelDebug)
}

func closeLog() {

}

func CreateDataFlie(filename string) *os.File {

	var f *os.File

	var err1 error

	if checkFileIsExist(filename) { //if file exit
		f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //open the file
		fmt.Println("file exit\n")
	} else {
		f, err1 = os.Create(filename) //ctreate file
		//		fmt.Println("file no exit\n")
	}
	check(err1)

	return f
}

func CloseFile(f *os.File) {
	f.Close()
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func WriteDataToFile(buf []byte, f *os.File) {

	f.Seek(0, 2)
	//	n2, err3 := f.Write(buf) //
	_, err3 := f.Write(buf) //
	check(err3)
	//	fmt.Printf("writed num=%d\n", n2)
	f.Sync()
}

func WriteStringToFile(str string, f *os.File) {

	f.Seek(0, 2)
	//	n2, err3 := f.WriteString(str) //
	_, err3 := f.WriteString(str) //
	check(err3)
	//	fmt.Printf("writed num=%d\n", n2)
	f.Sync()
}

func CreateRand() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int()
}
