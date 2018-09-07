package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"

	xco "github.com/sheenobu/go-xco"
)

func main() {
	flag.Parse()

	if e := validateArguments(); e != nil {
		fmt.Printf("%v\n", e)
		return
	}

	opts := xco.Options{
		Name:         *xmppName,
		SharedSecret: *xmppSharedSecret,
		Address:      joinHostPort(*xmppIP, *xmppPort),
	}

	c, err := xco.NewComponent(opts)
	if err != nil {
		fmt.Printf("error when connecting to server: %v\n", err)
		return
	}

	c.MessageHandler = messageHandler
	c.IqHandler = iqHandler

	err = c.Run()
	if err != nil {
		fmt.Printf("error when running component: %v\n", err)
		return
	}
}

func joinHostPort(ip string, port uint) string {
	return net.JoinHostPort(ip, fmt.Sprintf("%d", port))
}

func getTCPAddr(ip string, port uint) *net.TCPAddr {
	addr, _ := net.ResolveTCPAddr("tcp", joinHostPort(ip, port))
	return addr
}

func getPrekeyResponseFromRealServer(u string, data []byte) ([]byte, error) {
	addr := getTCPAddr(*rawIP, *rawPort)
	con, e := net.DialTCP(addr.Network(), nil, addr)
	if e != nil {
		return nil, e
	}
	defer con.Close()

	toSend := []byte{}
	toSend = appendShort(toSend, uint16(len(u)))
	toSend = append(toSend, []byte(u)...)
	toSend = appendShort(toSend, uint16(len(data)))
	toSend = append(toSend, data...)
	if _, e = con.Write(toSend); e != nil {
		return nil, e
	}
	con.CloseWrite()
	res, e := ioutil.ReadAll(con)
	if e != nil {
		return nil, e
	}
	if len(res) == 0 {
		return nil, errors.New("server closed connection without sending any data - this probably happened because of a malformed request")
	}

	res2, ss, ok := extractShort(res)
	if !ok || uint16(len(res2)) != ss {
		return nil, errors.New("unexpected length of data received")
	}
	return res2, nil
}
