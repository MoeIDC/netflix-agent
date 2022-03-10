package main

import (
	"github.com/evsio0n/log"
	"netflix_agent/utils"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var err error

func init() {
	log.SetDebug(utils.IsDebug)
	log.IsShowLogCatagory(false)
	log.SetSyslog(utils.IsSyslog, "netflix-agent")
	log.IsShowDate(utils.LogHaseDate)
	err = utils.ChangeIPv6()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}

func main() {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGUSR1, os.Interrupt)
	go func() {
		sig := <-signalChannel
		for {
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt:
				log.Info("get a signal:" + sig.String() + " , stop the process")
				utils.FlushNAT()
				log.Info("flush nat iptables success, exiting...")
				os.Exit(0)
			}
		}
	}()
	detectBlock()
}

func detectBlock() {
	for true {
		if !utils.TestUnblock() {
			log.Warn("Unblock failing, changing IP...")
			err := utils.ChangeIPv6()
			if err != nil {
				log.Error(err.Error())
				os.Exit(1)
			}
		} else {
			log.Debug("Unblock OK")
		}
		time.Sleep(time.Second * 30)
	}
}
