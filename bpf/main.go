package main

import (
	"fmt"
	"time"
	//bpf "github.com/holmanb/gostat/bpf"
)

func main(){
	var c CacheHit
	s := make(chan string)

	c.Init()
	if c.Debug {
		fmt.Printf("%24v%24v%24v%24v%24v%24v%24v%24v\n",
			"mpa",
			"mbd",
			"apcl",
			"apd",
			"total",
			"misses",
			"hits",
			"ratio")
	} else {
		fmt.Printf("%24v%24v%24v\n", "HITS", "MISSES", "RATIO")
	}
	c.Debug = true
	go func(s chan string) {
		for {
			c.Update(s)
			time.Sleep(time.Second)
		}
	}(s)
	for {
		fmt.Println(<-s)
	}
	c.Close()
}
