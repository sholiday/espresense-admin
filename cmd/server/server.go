package main

import (
	"fmt"
	eadmin "github.com/sholiday/espresense-admin"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	c, err := eadmin.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	a, err := eadmin.NewWebApp(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = a.Run()
	if err != nil {
		fmt.Println(err)
	}
}
