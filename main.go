package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net"

	xco "github.com/sheenobu/go-xco"
)

func main() {
	flag.Parse()
	opts := xco.Options{
		Name:         *xmppName,
		SharedSecret: *xmppSharedSecret,
		Address:      net.JoinHostPort(*xmppIP, fmt.Sprintf("%d", *xmppPort)),
	}

	c, err := xco.NewComponent(opts)
	if err != nil {
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

		return nil
	}

	c.IqHandler = func(_ *xco.Component, m *xco.Iq) error {
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

		c.Send(&resp)

		return nil
	}

	e := c.Run()
	if e != nil {
		panic(e)
	}
}

func xmlToString(x interface{}) string {
	enc, _ := xml.Marshal(x)
	return string(enc)

}

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

func tryDecodeXML(data []byte) (nspace, local string, ok bool) {
	token, _ := xml.NewDecoder(bytes.NewBuffer(data)).Token()
	if token == nil {
		return "", "", false
	}

	startElem, ok := token.(xml.StartElement)
	if !ok {
		return "", "", false
	}

	return startElem.Name.Space, startElem.Name.Local, true
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

// DiscoveryInfoQuery contains the deserialized information about a service discovery info query
// See: XEP-0030, Section 3
type DiscoveryInfoQuery struct {
	XMLName    xml.Name            `xml:"http://jabber.org/protocol/disco#info query"`
	Node       string              `xml:"node,omitempty"`
	Identities []DiscoveryIdentity `xml:"identity,omitempty"`
	Features   []DiscoveryFeature  `xml:"feature,omitempty"`
}

// DiscoveryIdentity contains identity information for a specific discovery
type DiscoveryIdentity struct {
	XMLName  xml.Name `xml:"http://jabber.org/protocol/disco#info identity"`
	Lang     string   `xml:"lang,attr,omitempty"`
	Category string   `xml:"category,attr"`
	Type     string   `xml:"type,attr"`
	Name     string   `xml:"name,attr"`
}

// DiscoveryFeature contains information about a specific discovery feature
type DiscoveryFeature struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#info feature"`
	Var     string   `xml:"var,attr"`
}

// DiscoveryItemsQuery contains a query for discovery items
type DiscoveryItemsQuery struct {
	XMLName        xml.Name        `xml:"http://jabber.org/protocol/disco#items query"`
	DiscoveryItems []DiscoveryItem `xml:"item,omitempty"`
}

// DiscoveryItem contains one discovery item
type DiscoveryItem struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#items item"`
	Jid     string   `xml:"jid,attr"`
	Name    string   `xml:"name,attr"`
	Node    string   `xml:"node, attr"`
}

func appendShort(l []byte, r uint16) []byte {
	return append(l, byte(r>>8), byte(r))
}

func extractShort(d []byte) ([]byte, uint16, bool) {
	if len(d) < 2 {
		return nil, 0, false
	}

	return d[2:], uint16(d[0])<<8 |
		uint16(d[1]), true
}

func getPrekeyResponseFromRealServer(u string, data []byte) []byte {
	addr, _ := net.ResolveTCPAddr("tcp", net.JoinHostPort(*rawIP, fmt.Sprintf("%d", *rawPort)))
	con, _ := net.DialTCP(addr.Network(), nil, addr)
	defer con.Close()

	toSend := []byte{}
	toSend = appendShort(toSend, uint16(len(u)))
	toSend = append(toSend, []byte(u)...)
	toSend = appendShort(toSend, uint16(len(data)))
	toSend = append(toSend, data...)
	con.Write(toSend)
	con.CloseWrite()
	res, _ := ioutil.ReadAll(con)
	res2, ss, _ := extractShort(res)
	if uint16(len(res2)) != ss {
		fmt.Printf("Unexpected length of data received\n")
		return nil
	}
	return res2
}
