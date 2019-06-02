package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type MethodInfo struct {
	Protocol   string
	MethodName string
	InputName  string
	OutputName string
}

var ServiceName string
var Methods []MethodInfo
var IgnoreType []string

func init() {
	IgnoreType = []string{"float", "int32", "int64", "uint32", "uint64", "bool", "string"}
}

func AnalyzeFileContent(filePath string) []MethodInfo{
	messages, serviceStr := LoadFile(filePath)
	if checkServiceStrInvalid(serviceStr) && _checkMessage(messages) {
		return Methods
	} else {
		panic("No service info or .proto file format error")
	}
}

func _isIgnoreType(_type string) bool {
	for _, iType := range IgnoreType {
		if _type == iType {
			return true
		}
	}
	return false
}

// _checkMessage ignore base type/check message is valid
func _checkMessage(messages []string) bool {
	if len(Methods) == 0 {
		return false
	}
	for _, _method := range Methods {
		inputFlag := false
		outputFlag := false
		// input or output type is base type
		if _isIgnoreType(_method.InputName) {
			inputFlag = true
		}
		if _isIgnoreType(_method.OutputName) {
			outputFlag = true
		}
		for _, msg := range messages {
			msg = strings.TrimSpace(msg)
			if strings.Compare(msg, _method.InputName) == 0 {
				inputFlag = true
			}
			if strings.Compare(msg, _method.OutputName) == 0 {
				outputFlag = true
			}
		}
		if inputFlag && outputFlag == false {
			return false
		}
	}
	return true
}

// lines:protocol {message:input} methodName({message:output})
func _checkArgs(content string) bool {
	lines := strings.Split(content, ";")
	if len(lines) == 0 {
		return false
	}
	for _, item := range lines {
		if item == "\n"{
			continue
		}
		item = strings.Trim(item, "\n")
		item = strings.TrimSpace(item)
		args := strings.Split(item, " ")
		if len(args) != 4 {
			return false
		}
		NameInput := strings.Split(args[1], "(")
		if NameInput[1][len(NameInput[1])-1:len(NameInput[1])] != ")" {
			return false
		}
		var info MethodInfo
		info.Protocol = args[0]
		OutputName := strings.Trim(args[3],"(")
		OutputName = strings.Trim(args[3],")")
		info.OutputName = OutputName[1:]
		info.MethodName = NameInput[0]
		info.InputName = NameInput[1][:len(NameInput[1])-1]
		Methods = append(Methods, info)
	}
	return true
}

func checkServiceStrInvalid(str string) bool {
	var start int
	var end int
	for i := 0; i < len(str); i++ {
		if str[i] == '{' {
			start = i
			ServiceName = str[:start]
			ServiceName = strings.TrimSpace(ServiceName)
		}
		if str[i] == '}' {
			end = i
		}
	}
	content := str[start+1 : end]
	return _checkArgs(content)
}

//
func LoadFile(filePath string) ([]string, string) {
	// var firstImportPosition int64
	var firstServicePosition int64
	var lastServicePosition int64
	var serviceContent string
	var messageContent []string

	protoFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("open file:%s error:%v", filePath, err)
		panic("open file error.")
	}
	defer func() {
		_ = protoFile.Close()
	}()
	fileData, err := ioutil.ReadAll(protoFile)
	if err != nil {
		fmt.Printf("load file:%s content error:%v", filePath, err)
		panic("read file error.")
	}

	fileStr := string(fileData)

	// get service content
	for i := 0; i < len(fileStr); i++ {
		if fileStr[i] == 's' && fileStr[i:i+7] == "service" {
			firstServicePosition = int64(i + 7 + 1)
			break
		}
	}
	for i := firstServicePosition; i < int64(len(fileStr)); i++ {
		if fileStr[i] == '}' {
			lastServicePosition = i + 1
		}
	}
	if firstServicePosition == 0 {
		serviceContent = ""
	} else {
		serviceContent = fileStr[firstServicePosition:lastServicePosition]
	}

	offset := 1
	for {
		i := offset
		for ; i < len(fileStr); i++ {
			if len(fileStr)-i < 8 {
				i = len(fileStr)
				break
			}
			if fileStr[i] == 'm' && fileStr[i:i+8] == "message " {
				firstOffset := i + 8
				lastOffset := firstOffset
				for fileStr[lastOffset] != '{' {
					lastOffset++
				}
				message := fileStr[firstOffset:lastOffset]
				messageContent = append(messageContent, message)
				offset = lastOffset
				break
			}
		}
		if i == len(fileStr) {
			break
		}
	}

	return messageContent, serviceContent
}
