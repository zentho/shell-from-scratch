package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func main() {
	builtins := []string{"exit", "echo", "pwd", "cd", "type"}
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		cmdFields := strings.Fields(input)
		cmd := cmdFields[0]

		switch cmd {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(strings.Join(cmdFields[1:], " "))
		case "pwd":
			if path, err := os.Getwd(); err != nil {
				fmt.Fprintf(os.Stderr, "pwd: error getting current directory: %v\n", err)
			} else {
				fmt.Println(path)
			}
		case "cd":
			if len(cmdFields) < 2 {
				fmt.Fprintln(os.Stderr, "cd: missing operand")
				continue
			}
			dir := cmdFields[1]
			if dir == "~" {
				dir = os.Getenv("HOME")
			}
			if err := os.Chdir(dir); err != nil {
				fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", dir)
			}
		case "type":
			if len(cmdFields) < 2 {
				fmt.Fprintln(os.Stderr, "type: missing operand")
				continue
			}
			target := cmdFields[1]
			if slices.Contains(builtins, target) {
				fmt.Printf("%s is a shell builtin\n", target)
			} else if path, err := exec.LookPath(target); err == nil {
				fmt.Printf("%s is %s\n", target, path)
			} else {
				notFound(target)
			}
		default:
			execCmd(cmd, cmdFields[1:])
		}
	}
}

func execCmd(cmd string, args []string) {
	command := exec.Command(cmd, args...)
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout

	if err := command.Run(); err != nil {
		notFound(cmd)
	}
}

func notFound(cmd string) {
	fmt.Printf("%s: not found\n", cmd)
}