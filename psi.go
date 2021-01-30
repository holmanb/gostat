package main

import (
	"bufio"
	"fmt"
	"os"
	"log"
	"regexp"
	"strconv"
)

type some struct {
	path string
	last_total uint64
	rex_total regexp.Regexp
}
type full struct {
	path string
	last_total uint64
	f *os.File
	s *bufio.Scanner
	rex_total regexp.Regexp
}

type resource interface {
	SubStr(string) string
}

type cpu struct {
	path string
	some some
	f *os.File
	s *bufio.Scanner
}

type mem struct {
	path string
	some some
	full full
	f *os.File
	s *bufio.Scanner
}
type io struct {
	some some
	full full
	path string
	f *os.File
	s *bufio.Scanner
}

type Psi struct {
	c cpu
	i io
	m mem
}

func (c *cpu) SubStr(s string)string{
	match := c.some.rex_total.FindStringSubmatch(s)
	if !(len(match) == 2){
		fmt.Printf("WARN: parsing error found length: %s in string %s", len(match),s)
		return ""
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


	file, err := os.Open(c.path)
	if err != nil {
		log.Fatal(err)
	}
	c.f = file
	c.s = bufio.NewScanner(file)
}

/*
* extracts total integer from the following line
* some avg10=0.00 avg60=0.00 avg300=0.00 total=66857405
 */
func extractInt(s string,r resource)uint64 {
	match := r.SubStr(s)
	i, err := strconv.ParseUint(match, 10, 64)
	if err != nil {
		fmt.Println("WARN: parsing error in ParseUint()", err)
	}
	return i
}

func (c *cpu) cpu_read() string{
	var diff uint64
	_, err := c.f.Seek(0, os.SEEK_SET)
	if err != nil {
		log.Fatal(err)
	}
	c.s.Scan()
	str := c.s.Text()
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

func (c *cpu) cpu_free(){
	c.f.Close()
}


func (p *Psi) Psi_init() {
	rex_total := regexp.MustCompile("total=(\\d+)$")
	p.c.cpu_init(*rex_total, "")
}
func (p *Psi) Get_psi(c chan string){
	c <- p.c.cpu_read()
}

func (p *Psi) Psi_close() {
	p.c.cpu_free()
}
