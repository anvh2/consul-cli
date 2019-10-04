package rpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"google.golang.org/grpc/credentials"
)

type ServerNameCheckCreds struct {
	mu       sync.Mutex
	sn       string
	expected string
}

func (c *ServerNameCheckCreds) ServerHandshake(rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	if _, err := io.WriteString(rawConn, c.sn); err != nil {
		fmt.Printf("Failed to write the server name %s to the client %v", c.sn, err)
		return nil, nil, err
	}
	return rawConn, nil, nil
}
func (c *ServerNameCheckCreds) ClientHandshake(ctx context.Context, addr string, rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	b := make([]byte, len(c.expected))
	errCh := make(chan error, 1)
	go func() {
		_, err := rawConn.Read(b)
		errCh <- err
	}()
	select {
	case err := <-errCh:
		if err != nil {
			fmt.Printf("Failed to read the server name from the server %v", err)
			return nil, nil, err
		}
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	}
	if c.expected != string(b) {
		fmt.Printf("Read the server name %s want %s", string(b), c.expected)
		return nil, nil, errors.New("received unexpected server name")
	}
	return rawConn, nil, nil
}
func (c *ServerNameCheckCreds) Info() credentials.ProtocolInfo {
	c.mu.Lock()
	defer c.mu.Unlock()
	return credentials.ProtocolInfo{}
}
func (c *ServerNameCheckCreds) Clone() credentials.TransportCredentials {
	c.mu.Lock()
	defer c.mu.Unlock()
	return &ServerNameCheckCreds{
		expected: c.expected,
	}
}
func (c *ServerNameCheckCreds) OverrideServerName(s string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.expected = s
	return nil
}
