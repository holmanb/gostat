package main

import (
	"fmt"
	"os"
	"time"
	"github.com/holmanb/gostat/bpf"
	)

type Resource interface {
	Init()
	Update(s chan string)
	Close()
}

type resource struct {
	Resource
	c chan string
	format string
	enabled bool
}

var resources = []resource {
	resource {
		enabled: true,
		format: "Cache Hit Rate %4v ",
		Resource: &bpf.CacheHit{},
		c: make(chan string),
	},
	resource {
		enabled: true,
		format: "| Mem: %4v ",
		Resource: &Mem{},
		c: make(chan string),
	},
	resource {
		enabled: true,
		format: "Cpu: %4v ",
		Resource: &Cpu{},
		c: make(chan string),
	},
	resource {
		enabled: true,
		format: "Io: %4v ",
		Resource: &Io{},
		c: make(chan string),
	},
	resource {
		enabled: true,
		format: "| %v ",
		Resource: &Time{},
		c: make(chan string),
	},
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
	var output string
	c := make(chan string)

	for _,resource := range resources {
		resource.Init()
	}
	go func(c chan string) {
		debug := false
		for i := range c {
			if debug {
				fmt.Println(i)
			} else {
				d.Update(i)
			}
		}
	}(c)
	for {
		for _,resource := range resources {
			go resource.Update(resource.c)
		}
		for _,resource := range resources {
			output += fmt.Sprintf(resource.format, <-resource.c)
		}
		c <- output
		output = ""
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
