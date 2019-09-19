package main

// This file handles the state configuration for each bulb.  The main
// routines just track Light_State in the Lights variable.  When the
// Hues API needs to present the data then we put that into a Light
// structure, with some hard coded defaults.
//
//
// func NewLight() -- returns a new Light_State structure that can be
//                    used to be added to the Lights map
// func get_or_create_light() -- if a light exists in the map, return
//                               that, otherwise create a new one
// func light_from_state() -- takes the Light_State structure and puts it
//                            into the Light structure for the Hue API

import (
	"log"
)

// Data structures for light modelling
//
// These structures are based on data from my white bulbs and cut down
// to what I think is the minimum necessary

// These are the values we need to track per light.
type Light_State struct {
	On   bool
	Bri  int
	name string
}

// This is what we return as an API call.  Most of this is constant!
type Light struct {
	State struct {
		On        bool   `json:"on"`
		Bri       int    `json:"bri"`
		Alert     string `json:"alert"`
		Mode      string `json:"mode"`
		Reachable bool   `json:"reachable"`
	} `json:"state"`
	SWUpdate struct {
		State       string `json:"state"`
		LastInstall string `json:"lastinstall"`
	} `json:"swupdate"`
	Type             string `json:"type"`
	Name             string `json:"name"`
	ModelID          string `json:"modelid"`
	ManufacturerName string `json:"manufacturername"`
	ProductName      string `json:"productname"`
	Capabilities     struct {
		Certified bool `json:"certified"`
		Control   struct {
			MinDim int `json:"mindimlevel"`
			MaxLum int `json:"maxlumen"`
		} `json:"control"`
		Streaming struct {
			Renderer bool `json:"renderer"`
			Proxy    bool `json:"proxy"`
		} `json:"streaming"`
	} `json:"capabilities"`
	Config struct {
		Archetype string `json:"archetype"`
		Function  string `json:"function"`
		Direction string `json:"direction"`
		Startup   struct {
			Mode       string `json:"mode"`
			Configured bool   `json:"configured"`
		} `json:"startup"`
	} `json:"config"`
	UniqueID   string `json:"uniqueid"`
	SWConfigID string `json:"swconfigid"`
	SWVersion  string `json:"swversion"`
	ProductID  string `json:"productid"`
}

// We will index lights based on their name
var Lights map[string]Light_State

func NewLight(name string) Light_State {
	return Light_State{
		On:   false,
		Bri:  0,
		name: name,
	}
}

func get_or_create_light(index, name string) Light_State {
	// Find the current light, create a new one if needed
	light, ok := Lights[index]
	if !ok {
		log.Println("Creating new light:", name)
		light = NewLight(name)
	}
	return light
}

func light_from_state(name string, state Light_State) Light {
	l := &Light{}

	// These are the only things that change.  Everything else is constant
	l.State.On = state.On
	l.State.Bri = state.Bri
	l.Name = state.name
	l.UniqueID = name

	l.Type = "Dimmable light"
	l.ModelID = "LWB014"
	l.ManufacturerName = "Philips"
	l.ProductName = "Hue white lamp"
	l.SWVersion = "1"
	l.SWConfigID = "69806BE9"
	l.ProductID = "Philips-LWB014-1-A19DLv4"

	l.State.Alert = "none"
	l.State.Mode = "homeautomation"
	l.State.Reachable = true

	l.SWUpdate.State = "noupdates"
	l.SWUpdate.LastInstall = "2018-12-14T00:08:29"

	l.Capabilities.Certified = true
	l.Capabilities.Control.MinDim = 5000
	l.Capabilities.Control.MaxLum = 840
	l.Capabilities.Streaming.Renderer = false
	l.Capabilities.Streaming.Proxy = false

	l.Config.Archetype = "classicbulb"
	l.Config.Function = "functional"
	l.Config.Direction = "omnidirectional"
	l.Config.Startup.Mode = "powerfail"
	l.Config.Startup.Configured = true

	return *l
}
