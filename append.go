// append.go
// this is a demo for https://www.upwork.com/jobs/
// Axel Gonzalez

package main

import (
	"fmt"
	"time"
	"os"
	//"fmt"
)


func loop(f *os.File, s string, t string, z int) {

	//d := date.Format(s)

	for {

		now := time.Now()

		x := fmt.Sprintf("%s | %s\n", now.Format(t), s)

		f.WriteString(x)
		fmt.Print(x)

		time.Sleep(time.Duration(z) * time.Second)
	}
}


func main() {

	f1, _ := os.OpenFile("log1", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	f2, _ := os.OpenFile("log2", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	f3, _ := os.OpenFile("log3", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	f4, _ := os.OpenFile("log4", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	defer f1.Close()
	defer f2.Close()
	defer f3.Close()
	defer f4.Close()

	go loop(f1, "this is log 1", "Jan 2, 2006 at 3:04:05 pm (MST)", 1)
	go loop(f2, "this is log 2", "2006-01-02T15:04:05Z", 2)
	go loop(f3, "this is log 3", "Jan 2, 2006 at 3:04:05 pm (MST)", 3)
	go loop(f4, "this is log 4", "2006-01-02T15:04:05Z", 4)


	for {
		time.Sleep(time.Duration(3)*time.Second)
	}
}

