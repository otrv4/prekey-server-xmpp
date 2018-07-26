package main

import (
	"errors"
	"flag"
	"regexp"
)

var (
	rawPort                 = flag.Uint("raw-port", 3242, "Port to connect to the raw server on")
	rawIP                   = flag.String("raw-address", "localhost", "Address to connect to the raw server on")
	xmppPort                = flag.Uint("xmpp-port", 53200, "Port to connect to the XMPP server on")
	xmppIP                  = flag.String("xmpp-address", "localhost", "Address to connect to the XMPP server on")
	xmppSharedSecret        = flag.String("shared-secret", "changeme", "Shared secret for authenticating to the XMPP server")
	xmppName                = flag.String("name", "changeme", "Name for the XMPP component. Usually something like prekeys.example.org")
	prekeyServerFingerprint = flag.String("fingerprint", "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", "Fingerprint for the prekey server. This is expected to be 56 bytes expressed in hexadecimal - thus 112 digits")
)

func isFingerprint(name string) bool {
	res, _ := regexp.MatchString("^[0-9A-Fa-f]{112}$", name)
	return res
}

func validateArguments() error {
	if !isFingerprint(*prekeyServerFingerprint) {
		return errors.New("fingerprint provided is not valid")
	}
	if *xmppName == "" || *xmppName == "changeme" {
		return errors.New("invalid xmpp name given")
	}
	if *xmppSharedSecret == "" {
		return errors.New("invalid shared secret given")
	}
	return nil
}
