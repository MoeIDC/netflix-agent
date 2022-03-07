package utils

import (
	"crypto/rand"
	"github.com/coreos/go-iptables/iptables"
	"github.com/evsio0n/log"
	"math/big"
	"net"
	"os"
)

var ip6t *iptables.IPTables

var IsDebug = GetConfig().GetBool("log.debug")
var HaseDate = GetConfig().GetBool("log.date.show")

func init() {
	log.SetDebug(IsDebug)
	log.IsShowDate(HaseDate)

	if IsDebug {
		log.Info("Debug mode")
		log.Info("IPv6 start" + GetConfig().GetString("net.ipv6.start"))
		log.Info("IPv6 end" + GetConfig().GetString("net.ipv6.end"))
		log.Info("Net interface name" + GetConfig().GetString("net.interface.name"))
	}
	var err error
	ip6t, err = iptables.NewWithProtocol(iptables.ProtocolIPv6)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func FlushNAT() error {
	return ip6t.ClearChain("nat", "POSTROUTING")
}

func ChangeIPtables() error {
	err := FlushNAT()
	if err != nil {
		return err
	}
	err = ip6t.Append("nat",
		"POSTROUTING",
		"-o",
		GetConfig().GetString("net.interface.name"),
		"-j",
		"SNAT",
		"--to-source",
		getRandomIPv6().String(),
	)
	if err != nil {
		return err
	}
	return nil
}

func ChangeIPv6() error {
	if GetConfig().GetString("mode") == "iptables" {
		err := FlushNAT()
		if err != nil {
			return err
		}
		err = ChangeIPtables()
		return err
	} else {
		err := FlushInterface()
		if err != nil {
			return err
		}
		err = ChangeInterfaceIP()
		return err
	}
}

func ip6toInt(IPv6Address net.IP) *big.Int {
	IPv6Int := big.NewInt(0)

	// from http://golang.org/pkg/net/#pkg-constants
	// IPv6len = 16
	IPv6Int.SetBytes(IPv6Address.To16())
	return IPv6Int
}

func getRandomIPv6() net.IP {
	start := net.ParseIP(GetConfig().GetString("net.ipv6.start"))
	end := net.ParseIP(GetConfig().GetString("net.ipv6.end"))
	startInt := ip6toInt(start)
	endInt := ip6toInt(end)

	offset, err := rand.Int(rand.Reader, big.NewInt(0).Sub(endInt, startInt))
	if err != nil {
		log.Info(err)
	}

	v6Int := big.NewInt(0).Add(startInt, offset)
	var v6Addr net.IP = v6Int.Bytes()

	return v6Addr
}
