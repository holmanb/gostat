package main

import (
	"fmt"
	"os"
	"time"
	"github.com/holmanb/gostat/bpf"
	)

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

func loop(d *Display){
	fmt.Println("running update loop")
	psi := Psi {}
	bpf_cache := bpf.CacheHit {}
	c := make(chan string)
	tm := make(chan string)
	ch := make(chan string)
	psi_cpu := make(chan string)
	psi_mem := make(chan string)
	psi_io := make(chan string)
	psi.Psi_init()
	bpf_cache.Init()
	go func(c chan string) {
		debug := true
		for i := range c{
			if debug {
				fmt.Println(i)
			} else {
				d.Update(i)
			}
		}
	}(c)
	for {
		go get_time(tm)
		go psi.Get_psi(psi_cpu, psi_mem, psi_io)
		go bpf_cache.Update(ch)
		c <- fmt.Sprintf("Page Cache %5v |Pressure Stats cpu:%6s mem:%6s io:%6s | %s",<-ch, <-psi_cpu, <-psi_mem, <-psi_io, <-tm)
		time.Sleep(time.Second)
	}
}

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
		loop(&d)
	}
}
