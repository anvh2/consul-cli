package integration

import (
	"context"
	"testing"

	pb "github.com/anvh2/consul-cli/grpc-gen/transfer"
	"github.com/stretchr/testify/assert"
	"gitlab.360live.vn/zalopay/zpi-pkg/testutils"
)

var clientTransfer = pb.NewTransferPointServiceClient(testutils.LocalConn(55215))

func TestTransferPoint(t *testing.T) {
	res, err := clientTransfer.TransferPoint(context.Background(), &pb.TransferRequest{
		ToID:   1,
		FromID: 2,
		Amount: 100,
	})

	assert.Nil(t, err)
	assert.Equal(t, int64(1), res.Code)
}
