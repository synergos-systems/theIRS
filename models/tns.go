// Code generated by https://github.com/gocomply/xsd2go; DO NOT EDIT.
// Models for http://schemas.xmlsoap.org/soap/envelope/
package tns

import (
	"encoding/xml"
)

// Element
type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`

	Header *Header `xml:"Header"`

	Body Body `xml:"Body"`
}

// Element
type Header struct {
	XMLName xml.Name `xml:"Header"`
}

// Element
type Body struct {
	XMLName xml.Name `xml:"Body"`
}

// XSD ComplexType declarations

type Fault struct {
	XMLName xml.Name

	Faultcode string `xml:"faultcode"`

	Faultstring string `xml:"faultstring"`

	Faultactor string `xml:"faultactor"`

	Detail *Detail `xml:"detail"`
}

type Detail struct {
	XMLName xml.Name
}

// XSD SimpleType declarations

type EncodingStyle string
