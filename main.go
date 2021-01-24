package main

import (
	"fmt"
	"os"
	)

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
		// TODO: implement status bar
		fmt.Println("not implemented")
	}

}
