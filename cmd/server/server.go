package main

import "fmt"
import eadmin "github.com/sholiday/espresense-admin"

func main() {
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
