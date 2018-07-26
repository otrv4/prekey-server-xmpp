package main

import (
	"fmt"

	xco "github.com/sheenobu/go-xco"
)

func messageHandler(c *xco.Component, m *xco.Message) error {
	res, e := getPrekeyResponseFromRealServer(fmt.Sprintf("%s@%s", m.Header.From.LocalPart, m.Header.From.DomainPart), []byte(m.Body))
	if e != nil {
		fmt.Printf("encountered error when communicating with the raw server: %v\n", e)
		return nil
	}

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

	e = c.Send(&resp)
	if e != nil {
		fmt.Printf("encountered error when sending response: %v\n", e)
	}

	return nil
}
