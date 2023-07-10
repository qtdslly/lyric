package main

import (
  "log"
  "syscall"
  "time"
  "fmt"
)

func main() {
	fmt.Println("Welcome to ZSCloud system")
	fmt.Println()
	fmt.Println("Documentation : https://help.aukoo.cn")
	fmt.Println("Support       : https://support.aukoo.cn")
	fmt.Println()

	p := time.Now()
	d := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",p.Year(),p.Month(),p.Day(),p.Hour(),p.Minute(),p.Second())
	fmt.Println("Last login : " + d)

	const bin = "/perm/home/figlet/figlet"
	if err := syscall.Exec(bin, []string{bin, "zscloud"}, nil); err != nil {
		fmt.Println(bin + " : run error")
		log.Fatal(err)
	}
}

