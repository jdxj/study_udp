package main

import (
	"fmt"
	"net"

	"golang.org/x/net/ipv4"
)

func main() {
	//1. 得到一个interface
	en4, err := net.InterfaceByName("en4")
	if err != nil {
		fmt.Println(err)
	}
	group := net.IPv4(224, 0, 0, 250)

	//2. bind一个本地地址
	c, err := net.ListenPacket("udp4", "0.0.0.0:1024")
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()

	//3.
	p := ipv4.NewPacketConn(c)
	if err := p.JoinGroup(en4, &net.UDPAddr{IP: group}); err != nil {
		fmt.Println(err)
	}

	//4.更多的控制
	if err := p.SetControlMessage(ipv4.FlagDst, true); err != nil {
		fmt.Println(err)
	}

	//5.接收消息
	b := make([]byte, 1500)
	for {
		n, cm, src, err := p.ReadFrom(b)
		if err != nil {
			fmt.Println(err)
		}

		if cm.Dst.IsMulticast() {
			if cm.Dst.Equal(group) {
				fmt.Printf("received: %s from <%s>\n", b[:n], src)
				n, err = p.WriteTo([]byte("world"), cm, src)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Unknown group")
				continue
			}
		}
	}
}
