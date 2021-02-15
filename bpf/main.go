package main

import (
	"fmt"
	"time"
)

func main(){
	var c CacheHit
	s := make(chan string)

	c.Init()
	if c.debug {
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
	c.debug = true
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
