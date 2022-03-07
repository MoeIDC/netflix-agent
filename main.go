package main

import (
	"fmt"
	"github.com/evsio0n/log"
	"netflix_agent/utils"
	"os"
	"os/user"
	"time"
)

var err error

func init() {
	//check if user is run as root
	if u, _ := user.Current(); u.Name != "root" {
		err = fmt.Errorf("You must run this program as root")
	} else {
		err = utils.ChangeIPv6()
	}
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
