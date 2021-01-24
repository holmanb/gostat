package main
// #cgo CFLAGS: -g -std=c99 -pedantic -Wall -O0 -I. -I/usr/include
// #cgo LDFLAGS: -L/usr/include  -lX11  -L/usr/lib -lc
// #include <stdlib.h>
// #include <X11/Xlib.h>
import "C"

import (
	"unsafe"
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
