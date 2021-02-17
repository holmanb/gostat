package main

import (
	"log"
	"regexp"
	"strconv"
	"io/ioutil"
)

type some struct {
	path string
	last_total uint64
	rex_total regexp.Regexp
}
type full struct {
	path string
	last_total uint64
	rex_total regexp.Regexp
}

type psi_resource interface {
	SubStr(string) string
}

type Cpu struct {
	path string
	some some
}

type Mem struct {
	path string
	some some
	full full
}
type Io struct {
	path string
	some some
	full full
}

type Psi struct {
	c Cpu
	i Io
	m Mem
}

var rex_total = *regexp.MustCompile("total=(\\d+)\\n$")

func (c *Cpu) SubStr(s string)string{
	match := c.some.rex_total.FindStringSubmatch(s)
	if !(len(match) == 2){
		log.Fatalf("WARN: cpu parsing error found length: %s for string %s", len(match),s)
	}
	return match[1]
}
func (m *Mem) SubStr(s string)string{
	match := m.some.rex_total.FindStringSubmatch(s)
	if !(len(match) == 2){
		log.Fatalf("WARN: mem parsing error found length: %s for string %s", len(match),s)
	}
	return match[1]
}
func (i *Io) SubStr(s string)string{
	match := i.some.rex_total.FindStringSubmatch(s)
	if !(len(match) == 2){
		log.Fatalf("WARN: io parsing error found length: %s for string %s", len(match),s)
	}
	return match[1]
}


func (c *Cpu) Init(){
	c.some.rex_total = rex_total
	c.path = "/proc/pressure/cpu"
}

func (m *Mem) Init(){
	m.some.rex_total = rex_total
	m.path = "/proc/pressure/memory"
}

func (i *Io) Init(){
	i.some.rex_total = rex_total
	i.path = "/proc/pressure/io"
}

/*
* extracts total integer from the following line
* some avg10=0.00 avg60=0.00 avg300=0.00 total=66857405
 */
func extractInt(s string,r psi_resource)uint64 {
	match := r.SubStr(s)
	i, err := strconv.ParseUint(match, 10, 64)
	if err != nil {
		log.Fatal("WARN: parsing error", err)
	}
	return i
}

func (c *Cpu) Update(s chan<- string) {
	var diff uint64
	b, err := ioutil.ReadFile(c.path)
	if err != nil {
		log.Fatal(err)
	}
	str := string(b)
	curr_total := extractInt(str,c)

	// report a sane initial value
	if c.some.last_total == 0 {
		diff = 0
	} else{
		diff = curr_total - c.some.last_total
	}
	c.some.last_total = curr_total
	s <- strconv.FormatUint(diff, 10)
}
func (m *Mem) Update(c chan<- string) {
	var diff uint64
	b, err := ioutil.ReadFile(m.path)
	if err != nil {
		log.Fatal(err)
	}
	str := string(b)
	curr_total := extractInt(str,m)

	// report a sane initial value
	if m.some.last_total == 0 { diff = 0
	} else{
		diff = curr_total - m.some.last_total
	}
	m.some.last_total = curr_total
	c <- strconv.FormatUint(diff, 10)
}
func (i *Io) Update(c chan<- string) {
	var diff uint64
	b, err := ioutil.ReadFile(i.path)
	if err != nil {
		log.Fatal(err)
	}
	str := string(b)
	curr_total := extractInt(str,i)

	// report a sane initial value
	if i.some.last_total == 0 {
		diff = 0
	} else{
		diff = curr_total - i.some.last_total
	}
	i.some.last_total = curr_total
	c <- strconv.FormatUint(diff, 10)
}

func (p *Psi) Init() {
	p.c.Init()
	p.m.Init()
	p.i.Init()
}

func (p *Psi) Update(c,m,i chan string){
	p.c.Update(c)
	p.m.Update(m)
	p.i.Update(i)
}

func (c *Cpu) Close (){}
func (m *Mem) Close (){}
func (i *Io) Close (){}
