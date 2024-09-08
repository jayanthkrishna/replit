package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func test() {
	cmd1 := exec.Command("ls", "./")

	// The `Output` method executes the command and
	// collects the output, returning its value
	out, _ := cmd1.Output()

	fmt.Println("Out1 : ", string(out))

	cmd := exec.Command("grep", "main")

	// Create a new pipe, which gives us a reader/writer pair
	reader, writer := io.Pipe()

	// assign the reader to Stdin for the command
	cmd.Stdin = reader
	// the output is printed to the console
	cmd.Stdout = os.Stdout
	go func() {
		writer.Write(out)
		writer.Close()
	}()

	if err := cmd.Run(); err != nil {
		fmt.Println("could not run command: ", err)
	}
}
