package main

import (
	"fmt"
	"os"
	"time"
	)

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
			c <- time.Now().Format(time.RFC1123)
			time.Sleep(time.Second)
		}
	}
}
