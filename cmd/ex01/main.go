package main

import (
	"log"

	"github.com/vishvananda/netlink"
)

func main() {
	// encontrar interafde de rede qual o ip
	link, err := netlink.LinkByName("enp0s5")
	if err != nil {
		log.Fatal(err)
	}
	newAddr, err := netlink.ParseAddr("192.168.1.2/24")
	if err != nil {
		log.Fatal(err)
	}
	// adicionar ip na interface de rede
	if err := netlink.AddrAdd(link, newAddr); err != nil {
		log.Fatal(err)
	}

}