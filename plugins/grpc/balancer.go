// Referencer: google.golang.org/grpc/grpclb/grpc_lb_v1

package rpc

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/anvh2/consul-cli/plugins/consul"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/grpclb/grpc_lb_v1"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// GrpcServiceWithLB -
type GrpcServiceWithLB struct {
	remoteBalancer *RemoteBalancer
	lbAddr         string
	lbServer       *grpc.Server
	backends       []*BaseGrpcService
	beName         string
	// beIPs          []net.IP
	bePorts []int

	lbListener  net.Listener
	beListenter []net.Listener
}

// NewServerWithExternalLoadBalancer -
func NewServerWithExternalLoadBalancer(grpcRegister GrpcRegister, targetName string, numServices int, ports []int) (*GrpcServiceWithLB, error) {
	if len(ports) != numServices {
		return nil, errors.New("Invalid number of services")
	}
	grpcServiceWithLB := &GrpcServiceWithLB{}
	grpcServiceWithLB.bePorts = ports
	grpcServiceWithLB.beName = targetName

	for i := 0; i < numServices; i++ {
		beService := NewGrpcServer(grpcRegister)
		grpcServiceWithLB.backends = append(grpcServiceWithLB.backends, beService)
	}

	return grpcServiceWithLB, nil
}

// Run -
// TODO: this function is leak gouroutins. When shutdown load balancer server, backends service still running
func (s *GrpcServiceWithLB) Run() error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// start base services
	go func() {
		for i, service := range s.backends {
			id, _ := uuid.NewV4()
			idstr := fmt.Sprintf("name-%s", id.String())
			config := consul.Config{
				ID:      idstr,
				Name:    "NameService",
				Tags:    []string{"DEV"},
				Address: "127.0.0.1",
				Port:    s.bePorts[i],
			}
			err := service.RegisterWithConsul(&config)
			if err != nil {
				return
			}

			service.AddShutdownHook(func() {
				service.DeRegisterFromConsul(idstr)
			})

			service.RegisterHealthCheck()

			go service.Run(s.bePorts[i])
		}
		<-sigs
	}()

	// start load balancer server
	go func() {
		var err error
		port := 55220
		s.lbListener, err = net.Listen("tcp", ":"+strconv.Itoa(port))
		// set load balancer address
		s.lbAddr = net.JoinHostPort("localhost", strconv.Itoa(s.lbListener.Addr().(*net.TCPAddr).Port))

		if err != nil {
			err = fmt.Errorf("failed to create the listener for the load balancer %v", err)
			fmt.Println(err)
			return
		}

		s.lbServer = grpc.NewServer()

		// TODO: register lb with consul
		id, _ := uuid.NewV4()
		idstr := fmt.Sprintf("lbname-%s", id.String())
		config := consul.Config{
			ID:      idstr,
			Name:    "LBNameService",
			Tags:    []string{"DEV"},
			Address: "127.0.0.1",
			Port:    port,
		}
		consul.Register(&config)
		// TODO: register healcheck
		grpc_health_v1.RegisterHealthServer(s.lbServer, &consul.HealthImpl{})

		s.remoteBalancer = NewRemoteBalancerServer(s.beName)
		grpc_lb_v1.RegisterLoadBalancerServer(s.lbServer, s.remoteBalancer)
		fmt.Println("Load Balancer Server is running on port: ", port)

		go func() {
			s.lbServer.Serve(s.lbListener)
		}()
		<-sigs
	}()

	<-sigs
	return nil
}
