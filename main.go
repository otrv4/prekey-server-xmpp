package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"

	xco "github.com/sheenobu/go-xco"
)

func main() {
	flag.Parse()

	// TODO: validate fingerprint
	// TODO: validate xmppName
	// TODO: validate shared secret

	opts := xco.Options{
		Name:         *xmppName,
		SharedSecret: *xmppSharedSecret,
		Address:      joinHostPort(*xmppIP, *xmppPort),
	}

	c, err := xco.NewComponent(opts)
	if err != nil {
		// TODO: print error properly
		panic(err)
	}

	c.MessageHandler = func(_ *xco.Component, m *xco.Message) error {
		res := getPrekeyResponseFromRealServer(fmt.Sprintf("%s@%s", m.Header.From.LocalPart, m.Header.From.DomainPart), []byte(m.Body))
		resp := xco.Message{
			Header: xco.Header{
				From: m.To,
				To:   m.From,
				ID:   m.ID,
			},
			Type:    m.Type,
			Body:    string(res),
			XMLName: m.XMLName,
		}

		c.Send(&resp)
		// TODO: print error properly

		return nil
	}

	c.IqHandler = func(_ *xco.Component, m *xco.Iq) error {
		ret, _, _ := processIQ(m)
		// TODO: print error properly
		resp := xco.Iq{
			Header: xco.Header{
				From: m.To,
				To:   m.From,
				ID:   m.ID,
			},
			Type:    "result",
			Content: xmlToString(ret),
			XMLName: m.XMLName,
		}

		c.Send(&resp)
		// TODO: print error properly

		return nil
	}

	e := c.Run()
	if e != nil {
		// TODO: print error properly
		panic(e)
	}
}

func joinHostPort(ip string, port uint) string {
	return net.JoinHostPort(ip, fmt.Sprintf("%d", port))
}

func getTCPAddr(ip string, port uint) *net.TCPAddr {
	addr, _ := net.ResolveTCPAddr("tcp", joinHostPort(ip, port))
	return addr
}

func getPrekeyResponseFromRealServer(u string, data []byte) []byte {
	addr := getTCPAddr(*rawIP, *rawPort)
	con, _ := net.DialTCP(addr.Network(), nil, addr)
	// TODO: print error properly
	defer con.Close()

	toSend := []byte{}
	toSend = appendShort(toSend, uint16(len(u)))
	toSend = append(toSend, []byte(u)...)
	toSend = appendShort(toSend, uint16(len(data)))
	toSend = append(toSend, data...)
	con.Write(toSend)
	// TODO: print error properly
	con.CloseWrite()
	res, _ := ioutil.ReadAll(con)
	// TODO: print error properly
	res2, ss, _ := extractShort(res)
	if uint16(len(res2)) != ss {
		fmt.Printf("Unexpected length of data received\n")
		return nil
	}
	return res2
}
