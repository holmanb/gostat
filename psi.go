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

type resource interface {
	SubStr(string) string
}

type cpu struct {
	path string
	some some
}

type mem struct {
	path string
	some some
	full full
}
type io struct {
	path string
	some some
	full full
}

type Psi struct {
	c cpu
	i io
	m mem
}

func (c *cpu) SubStr(s string)string{
	match := c.some.rex_total.FindStringSubmatch(s)
	if !(len(match) == 2){
		log.Fatalf("WARN: cpu parsing error found length: %s for string %s", len(match),s)
	}
	return match[1]
}
func (m *mem) SubStr(s string)string{
	match := m.some.rex_total.FindStringSubmatch(s)
	if !(len(match) == 2){
		log.Fatalf("WARN: mem parsing error found length: %s for string %s", len(match),s)
	}
	return match[1]
}
func (i *io) SubStr(s string)string{
	match := i.some.rex_total.FindStringSubmatch(s)
	if !(len(match) == 2){
		log.Fatalf("WARN: io parsing error found length: %s for string %s", len(match),s)
	}
	return match[1]
}


func (c *cpu) cpu_init(rex_total regexp.Regexp, path string){
	c.some.rex_total = rex_total
	if path == "" {
		c.path = "/proc/pressure/cpu"
	} else {
		c.path = path
	}
}

func (m *mem) mem_init(rex_total regexp.Regexp, path string){
	m.some.rex_total = rex_total
	if path == "" {
		m.path = "/proc/pressure/memory"
	} else {
		m.path = path
	}
}
func (i *io) io_init(rex_total regexp.Regexp, path string){
	i.some.rex_total = rex_total
	if path == "" {
		i.path = "/proc/pressure/io"
	} else {
		i.path = path
	}
}

/*
* extracts total integer from the following line
* some avg10=0.00 avg60=0.00 avg300=0.00 total=66857405
 */
func extractInt(s string,r resource)uint64 {
	match := r.SubStr(s)
	i, err := strconv.ParseUint(match, 10, 64)
	if err != nil {
		log.Fatal("WARN: parsing error", err)
	}
	return i
}

func (c *cpu) cpu_read() string{
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
	return strconv.FormatUint(diff, 10)
}
func (m *mem) mem_read() string{
	var diff uint64
	b, err := ioutil.ReadFile(m.path)
	if err != nil {
		log.Fatal(err)
	}
	str := string(b)
	curr_total := extractInt(str,m)

	// report a sane initial value
	if m.some.last_total == 0 {
		diff = 0
	} else{
		diff = curr_total - m.some.last_total
	}
	m.some.last_total = curr_total
	return strconv.FormatUint(diff, 10)
}
func (i *io) io_read() string{
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
	return strconv.FormatUint(diff, 10)
}



func (p *Psi) Psi_init() {
	rex_total := regexp.MustCompile("total=(\\d+)\\n$")
	p.c.cpu_init(*rex_total, "")
	p.m.mem_init(*rex_total, "")
	p.i.io_init(*rex_total, "")
}
func (p *Psi) Get_psi(c,m,i chan string){
	c <- p.c.cpu_read()
	m <- p.m.mem_read()
	i <- p.i.io_read()
}

