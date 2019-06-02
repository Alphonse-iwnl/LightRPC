package protocol

import "rpc/LightProtoPlug/common"

type Typer interface {
	OutputServiceCode(mode string) bool
}

func TypeServiceCode(mode string, services []common.MethodInfo, fileName string) bool {
	typer := NewHttpTyper(fileName, common.ServiceName, services)
	return typer.OutputServiceCode(mode)
}

func buildFileName(fileName string) string {
	index := 0
	for i, r := range fileName {
		if r == '/' {
			index = i
		}
	}
	return "./" + fileName[index+1:]
}
