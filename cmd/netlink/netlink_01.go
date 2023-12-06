package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

func main() {
	// Abrir um socket Netlink
	fd, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW, unix.NETLINK_ROUTE)
	if err != nil {
		log.Fatalf("Erro ao abrir o socket: %v", err)
	}
	defer unix.Close(fd)

	// Vincular o socket
	la := &unix.SockaddrNetlink{Family: unix.AF_NETLINK}
	if err := unix.Bind(fd, la); err != nil {
		log.Fatalf("Erro ao vincular o socket: %v", err)
	}

	// Preparar a mensagem Netlink
	pid := uint32(os.Getpid())
	log.Printf("PID: %d", pid)
	nlmsg := unix.NlMsghdr{
		Len:   unix.NLMSG_HDRLEN + uint32(len("Hello, Kernel!")),
		Type:  unix.RTM_GETROUTE,
		Flags: unix.NLM_F_REQUEST | unix.NLM_F_ACK,
		Seq:   1,
		Pid:   pid,
	}

	// Preparar o buffer da mensagem
	payload := []byte("Hello, Kernel!")
	buf := make([]byte, nlmsg.Len)
	binary.LittleEndian.PutUint32(buf[0:4], nlmsg.Len)
	binary.LittleEndian.PutUint16(buf[4:6], nlmsg.Type)
	binary.LittleEndian.PutUint16(buf[6:8], nlmsg.Flags)
	binary.LittleEndian.PutUint32(buf[8:12], nlmsg.Seq)
	binary.LittleEndian.PutUint32(buf[12:16], nlmsg.Pid)
	copy(buf[unix.NLMSG_HDRLEN:], payload)

	// Enviar a mensagem
	if err := unix.Sendto(fd, buf, 0, la); err != nil {
		log.Fatalf("Erro ao enviar a mensagem: %v", err)
	}

	// Receber e imprimir a resposta do kernel
	resp := make([]byte, 4096)
	nr, _, err := unix.Recvfrom(fd, resp, 0)
	if err != nil {
		log.Fatalf("Erro ao receber a resposta: %v", err)
	}
	fmt.Printf("Resposta do Kernel: %s\n", resp[:nr])
}
