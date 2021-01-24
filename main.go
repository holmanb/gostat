package main

// #cgo CFLAGS: -g -std=c99 -pedantic -Wall -O0 -I. -I/usr/include
// #cgo LDFLAGS: -L/usr/include  -lX11  -L/usr/lib -lc
// #include <stdlib.h>
// #include <X11/Xlib.h>
import "C"

import (
	"unsafe"
	"fmt"
	"os"
	)


type Display struct {
	display *C.Display
}

func (d *Display) XOpenDisplay() {
	d.display = C.XOpenDisplay(nil)
	if d.display == nil {
		panic("Can't open display")
	}
}

func (d Display) Update(s string){
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	w := C.XDefaultRootWindow(d.display)
	C.XStoreName(d.display, w, cstr)
	C.XFlush(d.display)
}

func (d Display) XCloseDisplay(){
	C.XCloseDisplay(d.display)
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
		s := os.Args[1]
		d.Update(s)
	} else {
		// TODO: implemenet status bar
		fmt.Println("not implemented")
	}

}
