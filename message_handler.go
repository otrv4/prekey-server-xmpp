package main

import (
	"fmt"
	"strings"

	xco "github.com/sheenobu/go-xco"
)

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
func validDomain(d string) bool {
	if *validDomains == "" {
		return true
	}
	ds := strings.Split(*validDomains, ",")
	return contains(ds, d)
}

func messageHandler(c *xco.Component, m *xco.Message) error {
	if !validDomain(m.Header.From.DomainPart) {
		fmt.Printf("Attempt to add data from unauthorized domain: %v@%v\n", m.Header.From.LocalPart, m.Header.From.DomainPart)
		return nil
	}

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
