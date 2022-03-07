package main

import (
	"log"
	"time"

	"netflix_agent/utils"
)

func init() {
	err := utils.ChangeIPv6()
	if err != nil {
		log.Panic(err)
	}
}

func main() {

	go detectBlock()

}

func detectBlock() {
	for true {
		if !utils.TestUnblock() {
			log.Println("Unblock failing, changing IP...")
			err := utils.ChangeIPv6()
			if err != nil {
				log.Panic(err)
			}
		} else {
			log.Println("Unblock OK")
		}
		time.Sleep(time.Second * 30)
	}
}
