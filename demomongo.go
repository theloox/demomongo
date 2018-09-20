// demomongo.go
// this is a demo for https://www.upwork.com/jobs/
// Axel Gonzalez

package main

import (
	"strings"
	//"strconv"
	"os"
	"fmt"
	//"os"
	"time"
	"bufio"
	"gopkg.in/mgo.v2"
)

type Record interface {
	ilog_time() time.Time
	ilog_msg() string
	ifile_name() string
	ilog_format() string
}

type Rec struct {
	Log_time time.Time
	Log_msg string
	File_name string
	Log_format string
}

func (r Rec) ilog_time() time.Time {
	return r.Log_time
}

func (r Rec) ilog_msg() string {
	return r.Log_msg
}

func (r Rec) ifile_name() string {
	return r.File_name
}

func (r Rec) ilog_format() string {
	return r.Log_format
}


func getrecord(s string, fn string, typ string) Rec {
	var ret Rec
	var ts time.Time
	var err error

	sa := strings.Split(s, " | ")

	d := sa[0]

	if typ == "1" {
		ts, err = time.Parse("Jan 2, 2006 at 3:04:05 pm (MST)", d)
	} else {
		ts, err = time.Parse("2006-01-02T15:04:05Z", d)
	}
	//fmt.Println(ts)

	if err == nil {
		ret = Rec{Log_time: ts, Log_msg: sa[1], File_name: fn, Log_format: typ}
	}

	return ret
}

func loop(in interface{Record}, c chan Record) {
	var last time.Time

	fmt.Println("name:", in.ifile_name(), "type: ", in.ilog_format())

	file, _ := os.Open(in.ifile_name())
	defer file.Close()

	for {
		// get all file
		// this would be more eficient on bsd/linux/osx with opening a tail subprocess

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {

			line := scanner.Text()
	   		//fmt.Println(line)

			z := getrecord(line, in.ifile_name(), in.ilog_format())
 	   		//fmt.Println(z)

			if z.Log_time.IsZero() {
				continue
			}

			if z.Log_time.Unix() > last.Unix() {
				c <- z

				last = z.Log_time
			}
    		}

		file.Seek(0, 0)
		fmt.Println("end ", in.ifile_name())

		//hardcoded sleep
		time.Sleep(time.Duration(10) * time.Second)
	}
}

func main() {

	args := make([]string, len(os.Args) - 1, len(os.Args))
	copy(args, os.Args[1:])

	c := make(chan Record)

	//sanity

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "panic: usage: %s <file1> <file2> ...\n", os.Args[0])
		os.Exit(2)
	}

	for _, fin := range args {

		fn := fin
		//fmt.Println(fn)
		if fn[0] == '/' {
			fmt.Fprintf(os.Stderr, "panic: no absolute paths\n")
			os.Exit(2)
		}

		if _, err := os.Stat(fn); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "panic: file doesn't exist %s\n", fn)
			os.Exit(2)
		}

	}

	// db
	session, err := mgo.Dial("localhost/task_en")
	db := session.DB("task_en").C("log")

	if err != nil {
		fmt.Fprintf(os.Stderr, "panic: db connection failed %s\n", err)
		os.Exit(2)
	}


	for _, fn := range args {

		var ft string

		f, _ := os.OpenFile(fn, os.O_RDONLY, 0600)

		data := make([]byte, 16)

		count, _ := f.Read(data)
		f.Close()

		if count > 7 && (data[5] == ',' || data[6] == ',') {
			ft = "1"
		} else if count > 11 && data[10] == 'T' {
			ft = "2"
		} else {
			fmt.Fprintf(os.Stderr, "panic: unknow file type %s\n", fn)
			os.Exit(2)
		}

		go loop(Rec{File_name: fn, Log_format: ft}, c)
	}

	for {
		x := <- c
		fmt.Println("new record:", x)
		//fmt.Printf("%+v\n", x)

		// now insert
		if db != nil {

		}
		err = db.Insert(&Rec{x.ilog_time(), x.ilog_msg(), x.ifile_name(), x.ilog_format()})
		if err != nil {
			fmt.Fprintf(os.Stderr, "panic: Can't insert in db\n")
		}
	}
}
