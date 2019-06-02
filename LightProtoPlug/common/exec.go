package common

import (
	"fmt"
	"os/exec"
)

const PBPLUGCOMMAND = "protoc -I "

func ExecPbPlug(filePath string, fileName string) string {
	fmt.Println("proto file path:" + filePath)
	fmt.Println("proto file name:" + fileName)
	path, file := buildPath(filePath)
	execCmd := PBPLUGCOMMAND + path + " --go_out=" + path + " " + file
	//execCmd := PBPLUGCOMMAND + path + " --go_out="+path+" " + file
	var cmd *exec.Cmd
	cmd = exec.Command("/bin/sh", "-c", execCmd)
	_, err := cmd.Output()
	_ = cmd.Start()
	if err != nil {
		panic("exec protoc command error")
	}
	return fileName + ".pb.go"
}

func buildPath(fileName string) (string, string) {
	index := 0
	for i, r := range fileName {
		if r == '/' {
			index = i
		}
	}
	return fileName[:index], fileName[index+1:]
}
