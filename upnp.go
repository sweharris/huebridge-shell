package main

// This file is a mashup of https://github.com/pborges/huejack and
// https://github.com/mdempsky/huejack
// I like the way mdempsky uses structures to create the XML, and
// the host, but I prefer pborge's use of net.ListenMulticastUDP()

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	upnp_multicast_address = "239.255.255.250:1900"
	upnpPath               = "/upnp/setup.xml"
)

func upnpSetup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("[UPNP] Setup URL called from ", r.RemoteAddr)
	type Root struct {
		XMLName struct{} `xml:"urn:schemas-upnp-org:device-1-0 root"`

		Major int `xml:"specVersion>major"`
		Minor int `xml:"specVersion>minor"`

		URLBase string

		DeviceType   string `xml:"device>deviceType"`
		FriendlyName string `xml:"device>friendlyName"`
		Manufacturer string `xml:"device>manufacturer"`
		ModelName    string `xml:"device>modelName"`
		ModelNumber  string `xml:"device>modelNumber"`
		UDN          string `xml:"device>UDN"`
	}

	x := Root{
		Major: 1,
		Minor: 0,

		URLBase: "http://" + r.Host + "/",

		DeviceType:   "urn:schemas-upnp-org:device:Basic:1",
		FriendlyName: "huebridge",
		Manufacturer: "Royal Philips Electronics",
		ModelName:    "Philips hue bridge 2012",
		ModelNumber:  "929000226503",
		UDN:          "uuid:f6543a06-800d-48ba-8d8f-bc2949eddc33",
	}

	w.Header().Set("Content-Type", "application/xml")
	xmlstr, _ := xml.Marshal(x)
	io.WriteString(w, string(xmlstr))
	io.WriteString(w, "\n")
}

// http://www.burgestrand.se/hue-api/api/discovery/

func upnpResponder(port int) {
	addr, err := net.ResolveUDPAddr("udp", upnp_multicast_address)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenMulticastUDP("udp", nil, addr)
	l.SetReadBuffer(1500)

	defer l.Close()
	log.Println("[UPNP] listening")

	var b [1500]byte
	for {
		n, raddr, err := l.ReadFromUDP(b[:])
		if err != nil {
			log.Fatal("[UPNP] ReadFromUDP failed:", err)
		}
		req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(b[:n])))
		if err != nil {
			log.Println("[UPNP] ReadRequest failed from", raddr, ":", err)
			continue
		}

		// TODO(mdempsky): Is this overly strict?  The "UPnP-basic-Basic-v1-Device.pdf"
		// spec suggests "Basic:1.0" should be acceptable too instead of "basic:1".
		// Regardless, this is sufficient for Echo.
		// TODO(mdempsky): According to draft-cai-ssdp-v1-03, we should also respond to
		// "St: ssdp:all" requests.
		if req.Method != "M-SEARCH" || req.URL.Path != "*" ||
			req.Header.Get("Man") != `"ssdp:discover"` ||
			req.Header.Get("St") != `urn:schemas-upnp-org:device:basic:1` {
			continue
		}
		log.Println("[UPNP] basic device discovery request from", raddr)

		upnpAnswer(port, raddr)
	}
}

func upnpAnswer(port int, raddr *net.UDPAddr) {
	c, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	url := fmt.Sprintf("http://%s:%d"+upnpPath, c.LocalAddr().(*net.UDPAddr).IP, port)

	// TODO(mdempsky): Better way to format HTTP responses?
	var buf bytes.Buffer
	buf.WriteString("HTTP/1.1 200 OK\r\n")
	http.Header{
		// According to draft-frystyk-http-extensions-03, we
		// MUST include a no-cache="Ext" directive in the Cache-Control
		// header field, but then Echo ignores our response.
		"Cache-Control": {`max-age=300`},
		"Ext":           {``},
		"Location":      {url},
		"Opt":           {`"http://schemas.upnp.org/upnp/1/0/"; ns=01`},
		"St":            {`urn:schemas-upnp-org:device:basic:1`},
		"Usn":           {`uuid:f6543a06-800d-48ba-8d8f-bc2949eddc33`},
	}.Write(&buf)
	buf.WriteString("\r\n")

	// log.Println("Sending response: ",string(buf.Bytes()))
	_, err = c.Write(buf.Bytes())
	if err != nil {
		log.Println("Error writing UPnP response:", err)
	}
}
