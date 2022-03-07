package utils

import (
	"fmt"
	"os/exec"
)

func FlushInterface() error {
	_, err := exec.Command(fmt.Sprintf("/sbin/ip -6 addr flush dev %s", GetConfig().GetString("net.interface.name"))).Output()
	return err
}

func ChangeInterfaceIP() error {
	err := FlushInterface()
	if err != nil {
		return err
	}
	_, err = exec.Command(fmt.Sprintf("/sbin/ip -6 addr add %s dev %s", getRandomIPv6().String(), GetConfig().GetString("net.interface.name"))).Output()
	return err
}
