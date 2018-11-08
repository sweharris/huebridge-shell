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
// (Well, Reachable is always True, but it's still part of State)
// And "name" is an internal structure that doesn't export... I hope!
type Light_State struct {
	On        bool `json:"on"`
	Bri       int  `json:"bri"`
	Reachable bool `json:"reachable"`
	name      string
}

// Except for State, we expect these to be constant for every light.
type Light struct {
	State            Light_State `json:"state"`
	Type             string      `json:"type"`
	Name             string      `json:"name"`
	ModelID          string      `json:"modelid"`
	ManufacturerName string      `json:"manufacturername"`
	ProductName      string      `json:"productname"`
	UniqueID         string      `json:"uniqueid"`
	SWVersion        string      `json:"swversion"`
}

// We will index lights based on their name
var Lights map[string]Light_State

func NewLight(name string) Light_State {
	return Light_State{
		On:        false,
		Bri:       0,
		Reachable: true,
		name:      name,
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
	return Light{
		State:            state,
		Type:             "Dimmable light",
		Name:             state.name,
		ModelID:          "LWB014",
		ManufacturerName: "Philips",
		ProductName:      "Hue white lamp",
		UniqueID:         name,
		SWVersion:        "1",
	}
}
