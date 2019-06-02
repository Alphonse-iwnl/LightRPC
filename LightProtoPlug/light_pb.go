package main

import (
	"fmt"
	"os"
	"strings"
)
import . "rpc/LightProtoPlug/common"
import (
	. "rpc/LightProtoPlug/protocol"
)

func main() {
	// load command
	filePath, fileName, opt := CheckArgs()
	if len(strings.Split(filePath, "/")) == 1 {
		dirPath, _ := os.Getwd()
		filePath = dirPath + "/" + filePath
	}
	// func()
	// read pb-> check 'service' struct in .pb->
	// load 'service' info to mem
	Services := AnalyzeFileContent(filePath)
	if Services == nil {
		return
	}
	// exec origin pb-plug ->create pb.go
	pbFileName := ExecPbPlug(filePath, fileName)
	// print code to file with service info
	if !TypeServiceCode(opt, Services, pbFileName) {
		fmt.Println("Type service code error.")
		return
	}
	fmt.Println("Generate Success.")
}
