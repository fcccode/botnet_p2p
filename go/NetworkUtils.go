package main

import (
	"net/http"
	"io/ioutil"
	"net"
	"log"
)

func getRemoteIP() (string, error) {
	var remoteIP = "127.0.0.1"
	response, err := http.Get("https://api.ipify.org")
	if err != nil {
		return remoteIP, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return remoteIP, err
		}
		remoteIP = string(bodyBytes)
	}
	return remoteIP, nil
}

func getLocalIPs() ([]string, error) {
	ips := make([]string, 0, 8)
	interfaces, err := net.Interfaces()
	if err != nil {
		return ips, err
	}
	for _, i := range interfaces {
		addresses, err := i.Addrs()
		if err != nil {
			return ips, err
		}
		for _, addr := range addresses {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.To4() != nil {
				ips = append(ips, ip.To4().String())
			}
		}
	}
	return ips, nil
}

func checkNAT() (bool, error) {
	remoteIP, err := getRemoteIP()
	if err != nil {
		return true, err
	}
	log.Println("Remote IP:", remoteIP)
	localIPs, err := getLocalIPs()

	for _, localIP := range localIPs {
		if localIP == remoteIP {
			return false, nil
		}
	}

	return true, nil
}