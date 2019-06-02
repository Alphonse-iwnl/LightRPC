package server

import (
	"context"
	"github.com/golang/protobuf/proto"
	"math/rand"
	. "rpc/LightRpc/proto"
	"strconv"
)

type ExecService struct {
}

func (es *ExecService) ExecService(ctx context.Context, in *Test01) (*Test, error) {
	output := &Test{}

	output.Message = in.Message
	output.Id = proto.String(strconv.Itoa(int(rand.Int63())))
	return output, nil
}
