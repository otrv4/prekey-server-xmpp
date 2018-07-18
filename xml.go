package main

import (
	"bytes"
	"encoding/xml"
)

func xmlToString(x interface{}) string {
	enc, _ := xml.Marshal(x)
	return string(enc)
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
