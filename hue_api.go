package main

// This is based on the mdempsky code, but modified to handle the light
// structures I prefer

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func hue_ListenAndServe() error {
	laddr := net.TCPAddr{Port: 5000}
	l, err := net.ListenTCP("tcp4", &laddr)
	if err != nil {
		return err
	}

	router := httprouter.New()
	router.GET(upnpPath, upnpSetup)

	// router.GET("/api/:userId", hue_getLightsList)
	router.PUT("/api/:userId/lights/:lightId/state", hue_setLightState)
	router.GET("/api/:userId/lights/:lightId", hue_getLightInfo)
	router.GET("/api/:userId/lights", hue_getLightsList)

	go upnpResponder(l.Addr().(*net.TCPAddr).Port)
	return http.Serve(l, requestLogger(router))
}

// Handler:
//      state is the state of the "light" after the handler function
//  if error is set to true echo will reply with "sorry the device is not responding"
type Handler func(key, val int)

func requestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Println("[WEB]", r.RemoteAddr, r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}

func hue_setLightState(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer r.Body.Close()
	var req struct {
		On  *bool `json:"on"`
		Bri *int  `json:"bri"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	// log.Printf("[DEVICE] req = %#v", req)

	lightID := p.ByName("lightId")
	thislight, ok := Lights[lightID]
	if !ok {
		log.Printf("device %v missing", lightID)
		return
	}

	onoff := ""
	if req.On != nil {
		if *req.On {
			onoff = "ON"
		} else {
			onoff = "OFF"
		}
	}
	bristr := ""
	if req.Bri != nil {
		bristr = strconv.Itoa(*req.Bri)
	}
	log.Println("[API] Update for", lightID, "from", r.RemoteAddr, "setting", onoff, bristr)
	send_cmd_to_child("LIGHT#" + thislight.name + "#" + onoff + "#" + bristr)

	// Let's assume the send to child always worked.  If it
	// didn't then the result will eventually update a few seconds
	// later
	//
	// So we create a result that looks like the current status
	// of the light.  We'll report the current state, which
	// might have been updated from the child; let's give it
	// a second
	time.Sleep(time.Second / 10)

	// A bit of a kludge...
	this_light := Lights[lightID]
	// I have no idea how to make this work properly with golang
	// structures, so I'm just gonna fake up the JSON and send
	// it directly
	res := `[ {"success":{"/lights/` + lightID + `/state/bri":` + strconv.Itoa(this_light.Bri) + `}}, {"success":{"/lights/` + lightID + `/state/on":` + strconv.FormatBool(this_light.On) + `}} ]`
	// log.Println(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(res))
}

func hue_getLightsList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("[API] Asking for all lights from ", r.RemoteAddr)
	// Built a light hierarchy based on all the light_status entries
	hue_lights := make(map[string]Light)
	for id, light := range Lights {
		hue_lights[id] = light_from_state(id, light)
	}
	sendJSON(w, hue_lights)
}

func hue_getLightInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	lightID := p.ByName("lightId")
	log.Println("[API] Asking for", lightID, "from", r.RemoteAddr)
	l, ok := Lights[lightID]
	if !ok {
		log.Printf("device %v missing", lightID)
		return
	}

	light := light_from_state(lightID, l)

	sendJSON(w, light)
}

func sendJSON(w http.ResponseWriter, val interface{}) {
	w.Header().Set("Content-Type", "application/json")

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(val); err != nil {
		log.Fatal("[WEB] error encoding json: ", err)
	}
	// log.Print("sending JSON response: ", buf.String())

	w.Write(buf.Bytes())
}

// It looks like the API calls don't like spaces in the index name
// so we'll just strip them out
func name_to_index(name string) string {
	return strings.Replace(name, " ", "", -1)
}
