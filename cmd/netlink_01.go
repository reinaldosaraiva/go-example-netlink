package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func main() {
	// Abrir um socket Netlink
	fd, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW, unix.NETLINK_ROUTE)
	if err != nil {
		fmt.Println("Erro ao abrir o socket: ", err)
		os.Exit(1)
	}
	defer unix.Close(fd)
	// Vincular o socket a um endereço
	sa := &unix.SockaddrNetlink{Family: unix.AF_NETLINK}
	err = unix.Bind(fd, sa)
	if err != nil {
		fmt.Println("Erro ao vincular o socket: ", err)
		os.Exit(1)
	}
	// criar um mensagem Netlink de exemplo
	msg := &unix.NlMsghdr{
		Len:   uint32(unix.SizeofNlMsghdr + len([]byte("Hello Kernel"))),
		Type:  unix.RTM_GETROUTE,
		Flags: unix.NLM_F_REQUEST | unix.NLM_F_DUMP,
		Seq:   1,
		Pid:   uint32(os.Getpid()),
	}
	// enviar a mensagem
	_, err = unix.Sendmsg(fd, []byte("Hello Kernel"), &msg, 0)
	if err != nil {
		fmt.Println("Erro ao enviar a mensagem: ", err)
		os.Exit(1)
	}
	// receber a mensagem
	buf := make([]byte, 8192)
	nr, _, err := unix.Recvfrom(fd, buf, 0)
	if err != nil {
		fmt.Println("Erro ao receber a mensagem: ", err)
		os.Exit(1)
	}
	// imprimir a mensagem
	fmt.Println("Resposta do Kernel: ", string(buf[:nr]))
}
