package main

import "time"

type Time struct {
}

func (t *Time) Init(){}

func (t *Time) Update(s chan string){
	s <- time.Now().Format(time.RFC1123)
}
func (t *Time) Close(){}
