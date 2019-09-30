package integration

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/anvh2/consul-cli/grpc-gen/transfer"
	"github.com/stretchr/testify/assert"
	"gitlab.360live.vn/zalopay/zpi-pkg/testutils"
)

var client = pb.NewTransferPointServiceClient(testutils.LocalConn(55211))

func TestTransferPoint(t *testing.T) {
	res, err := client.TransferPoint(context.Background(), &pb.TransferRequest{
		ToID:   1,
		FromID: 0,
		Amount: 100,
	})

	assert.Nil(t, err)
	fmt.Println(res)
}
