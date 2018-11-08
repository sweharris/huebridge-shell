package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// This handles the LIGHT status update sent from the child process
//    name on/off brightness
func set_light_status(status []string) {
	name := status[0]

	power := false
	if strings.ToUpper(status[1]) == "ON" {
		power = true
	}

	bri, err := strconv.Atoi(status[2])
	if err != nil {
		bri = 0
	}

	name_index := name_to_index(name)
	light := get_or_create_light(name_index, name)

	if light.On != power {
		log.Println("Got light update for", name, "setting Power", power)
		light.On = power
	}

	if light.Bri != bri {
		log.Println("Got light update for", name, "setting Brightness", bri)
		light.Bri = bri
	}

	Lights[name_index] = light
}

// This takes a list of lights as reported by the LIST status update
// and ensures that this list matches the Lights map.  Basically we
// create any missing entries with defaults, then delete ones that are
// no longer relevant

func update_light_list(lights []string) {
	// We'll make the list of lights into a map so we can quickly compare
	list := make(map[string]bool)

	for _, v := range lights {
		name_index := name_to_index(v)
		list[name_index] = true
		Lights[name_index] = get_or_create_light(name_index, v)
	}

	for v := range Lights {
		_, ok := list[v]
		if !ok {
			log.Println("Deleting", v)
			delete(Lights, v)
		}
	}
}

// Input to this function will be output from the child.  The line
// either needs to be
//    LIST#name1#name2 ...
// or
//    LIGHT#name#on/off#brightness
//
// Note that names can have spaces in them; the separator is #

func on_read(s string) {
	input := strings.Split(s, "#")
	cmd := strings.ToUpper(input[0])
	if cmd == "LIGHT" {
		set_light_status(input[1:])
	} else if cmd == "LIST" {
		update_light_list(input[1:])
	} else {
		log.Println("Ignoring", s)
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Syntax: ", os.Args[0], " child_program")
	}

	// Create the new lights map
	Lights = make(map[string]Light_State)

	// Set the child process running to update the light status
	start_child(os.Args[1])
	defer child.Wait()
	go read_child(on_read)

	// And now start the Hue listener
	log.Fatal(hue_ListenAndServe())

	stdin := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter text: ")
		text, _ := stdin.ReadString('\n')
		send_cmd_to_child(text)
	}
}
