package main

// These routines handle the asynchronous communication to the child.
// Basically the caller is expected to run
//   start_child
//   go read_child
//   send_cmd_to_child
// read_child takes a callback function that is called each time a line
// is read from the child.

// TODO better handling if the child stops working; maybe try and
// restart it?

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

var child *exec.Cmd
var send_to_child io.WriteCloser
var read_from_child io.ReadCloser

// Run the child, with filehandles to stdin/stdout
func start_child(cmd string) {
	log.Println("Starting", cmd)
	child = exec.Command(cmd)
	send_to_child, _ = child.StdinPipe()
	read_from_child, _ = child.StdoutPipe()
	child.Start()
}

type read_handler func(string)

// Read the child's output, and for each line we call a handler
// function
func read_child(handle read_handler) {
	scanner := bufio.NewScanner(read_from_child)

	for scanner.Scan() {
		line := scanner.Text()
		handle(line)
	}
	log.Fatal("Child no longer talks to me.  Aborted")
}

// Send a command to the child.
func send_cmd_to_child(cmd string) {
	// log.Println("Sending ", cmd)
	_, err := send_to_child.Write([]byte(cmd + "\n"))
	if err != nil {
		log.Println("Error sending command:", cmd, " -", err)
	}
}
