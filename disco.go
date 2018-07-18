package main

import "encoding/xml"

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
