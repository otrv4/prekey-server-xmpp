package main

import (
	"fmt"

	xco "github.com/sheenobu/go-xco"
)

func unknownIQ(stanza *xco.Iq) (ret interface{}, iqtype string, ignore bool) {
	fmt.Printf("Unknown IQ: %s\n", stanza.Content)
	return nil, "", false
}

type iqFunction func(*xco.Iq) (interface{}, string, bool)

var knownIQs = map[string]iqFunction{}

func registerKnownIQ(stanzaType, fullName string, f iqFunction) {
	knownIQs[stanzaType+" "+fullName] = f
}

func getIQHandler(stanzaType, namespace, local string) iqFunction {
	f, ok := knownIQs[fmt.Sprintf("%s %s %s", stanzaType, namespace, local)]
	if ok {
		return f
	}
	return unknownIQ
}

func processIQ(stanza *xco.Iq) (ret interface{}, iqtype string, ignore bool) {
	if nspace, local, ok := tryDecodeXML([]byte(stanza.Content)); ok {
		return getIQHandler(stanza.Type, nspace, local)(stanza)
	}
	return nil, "", false
}

func discoInfoIQ(_ *xco.Iq) (interface{}, string, bool) {
	return DiscoveryInfoQuery{
		Identities: []DiscoveryIdentity{
			{
				Category: "auth",
				Type:     "otr-prekey",
				Name:     "OTR Prekey Server",
			},
		},
		Features: []DiscoveryFeature{
			{Var: "http://jabber.org/protocol/disco#info"},
			{Var: "http://jabber.org/protocol/disco#items"},
			{Var: "http://jabber.org/protocol/otrv4-prekey-server"},
		},
	}, "", false
}

func discoItemsIQ(_ *xco.Iq) (interface{}, string, bool) {
	return DiscoveryItemsQuery{
		DiscoveryItems: []DiscoveryItem{
			{
				Jid:  *xmppName,
				Node: "fingerprint",
				Name: *prekeyServerFingerprint,
			},
		},
	}, "", false
}

func init() {
	registerKnownIQ("get", "http://jabber.org/protocol/disco#info query", discoInfoIQ)
	registerKnownIQ("get", "http://jabber.org/protocol/disco#items query", discoItemsIQ)
}
