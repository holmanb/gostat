package main

import (
	"fmt"
	"os"
	"time"
	)

func get_disk_space() string {
	return "FULL"
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
		fmt.Println("running buildin update loop")
		c := make(chan string)
		go func(c chan string) {
			for i := range c{
				d.Update(i)
			}
		}(c)
		for {
			t := time.Now().Format(time.RFC1123)
			d := get_disk_space()
			s := fmt.Sprintf("[Disk : %s]  %s",d, t)
			c <- s
			time.Sleep(time.Second)
		}
	}
}
