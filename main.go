package main

import (
	"github.com/evsio0n/log"
	"netflix_agent/utils"
	"os"
	"time"
)

var err error

func init() {
	log.SetDebug(true)
	log.IsShowDate(true)
	//check if user is run as root
	if os.Geteuid() != 0 {
		log.Error("You must run this program as root")
		os.Exit(1)
	}
	err = utils.ChangeIPv6()
}

func main() {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	go detectBlock()

}

func detectBlock() {
	for true {
		if !utils.TestUnblock() {
			log.Info("Unblock failing, changing IP...")
			err := utils.ChangeIPv6()
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
		} else {
			log.Info("Unblock OK")
		}
		time.Sleep(time.Second * 30)
	}
}
