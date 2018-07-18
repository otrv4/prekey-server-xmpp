package main

import (
	"fmt"

	xco "github.com/sheenobu/go-xco"
)

func iqHandler(c *xco.Component, m *xco.Iq) error {
	ret, _, _ := processIQ(m)
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

	e := c.Send(&resp)
	if e != nil {
		fmt.Printf("encountered error when sending response: %v\n", e)
	}

	return nil
}
