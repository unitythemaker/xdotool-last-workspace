package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

var currentWorkspace string
var lastWorkspace string
var lastWorkspaceDir = "xdotool-last-workspace"
var lastWorkspaceFilePath string

func main() {
	cachePath, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	lastWorkspaceDir = cachePath + "/" + lastWorkspaceDir
	if os.Getenv("LAST_WORKSPACE_DIR") != "" {
		lastWorkspaceDir = os.Getenv("LAST_WORKSPACE_DIR")
	}
	lastWorkspaceFilePath = lastWorkspaceDir + "/last_workspace"
	// create the directory if it doesn't exist
	os.MkdirAll(lastWorkspaceDir, os.ModePerm)

	args := os.Args[1:]
	if len(args) > 0 && args[0] == "-d" {
		startDaemon()
		return
	}
	switchToLastWorkspace()
}

func startDaemon() {
	currentWorkspace = getCurrentWorkspace()
	lastWorkspace = currentWorkspace
	for {
		currentWorkspace = getCurrentWorkspace()
		if currentWorkspace != lastWorkspace {
			fmt.Println(currentWorkspace)
			writeLastWorkspace(lastWorkspace)
			lastWorkspace = currentWorkspace
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func switchToLastWorkspace() {
	lastWorkspace := readLastWorkspace()
	setCurrentWorkspace(lastWorkspace)
}

func getCurrentWorkspace() string {
	out, err := exec.Command("xdotool", "get_desktop").Output()
	if err != nil {
		log.Fatalln("WARN: Error while executing xdotool command:", err)
	}
	return string(out[0])
}

func setCurrentWorkspace(workspace string) {
	exec.Command("xdotool", "set_desktop", workspace).Run()
}

func readLastWorkspace() string {
	file, err := os.Open(lastWorkspaceFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	var lastWorkspace string
	fmt.Fscanf(file, "%s", &lastWorkspace)
	return lastWorkspace
}

func writeLastWorkspace(workspace string) {
	file, err := os.Create(lastWorkspaceFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	_, err = file.WriteString(workspace)
	if err != nil {
		panic(err)
	}
}
