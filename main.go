package main

import (
	"fmt"
	"os"
	"time"
	)

func get_disk_space(s chan string){
	s <- "FULL"
}
func get_time(s chan string){
	s <- time.Now().Format(time.RFC1123)
}

/*
to implement
============
resources:
memory /proc/meminfo
cpu util
storage space

temperature:
by drive
battery
*/


func main(){
	d := Display{}
	d.XOpenDisplay()
	defer d.XCloseDisplay()

	if len(os.Args) > 2 {
		fmt.Printf("usage:\n\t%s - to use buildin status bar or\n\t%s <string> - to manually set status\n",
			os.Args[0],
			os.Args[0])
		os.Exit(1)
	} else if len(os.Args) == 2 {
		fmt.Println("updating status with user defined string:", os.Args[1])
		s := os.Args[1]
		d.Update(s)
	} else {
		fmt.Println("running update loop")
		c := make(chan string)
		tm := make(chan string)
		ds := make(chan string)
		go func(c chan string) {
			for i := range c{
				d.Update(i)
			}
		}(c)
		for {
			go get_time(tm)
			go get_disk_space(ds)

			s := fmt.Sprintf("[Disk : %s]  %s",<-ds, <-tm)
			c <- s
			time.Sleep(time.Second)
		}
	}
}
