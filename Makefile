genpb:
	protoc -I/usr/local/include -Iidl \
		-I$$GOPATH/src \
		-I$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		-I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--plugin=protoc-gen-gogo=$$GOPATH/bin/protoc-gen-gogo \
		--gogo_out=plugins=grpc:grpc-gen/counter \
		idl/counter.proto

	protoc -I/usr/local/include -Iidl \
		-I$$GOPATH/src \
		-I$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		-I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--plugin=protoc-gen-gogo=$$GOPATH/bin/protoc-gen-gogo \
		--gogo_out=plugins=grpc:grpc-gen/transfer \
		idl/transfer.proto

	protoc -I/usr/local/include -Iidl \
		-I$$GOPATH/src \
		-I$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		-I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--plugin=protoc-gen-gogo=$$GOPATH/bin/protoc-gen-gogo \
		--gogo_out=plugins=grpc:grpc-gen/user \
		idl/user.proto

	protoc -I/usr/local/include -Iidl \
		-I$$GOPATH/src \
		-I$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		-I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--plugin=protoc-gen-gogo=$$GOPATH/bin/protoc-gen-gogo \
		--gogo_out=plugins=grpc:grpc-gen/echo \
		idl/echo.proto

		protoc -I/usr/local/include -Iidl \
		-I$$GOPATH/src \
		-I$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		-I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--plugin=protoc-gen-gogo=$$GOPATH/bin/protoc-gen-gogo \
		--gogo_out=plugins=grpc:grpc-gen/load-balancer \
		idl/load_balancer.proto
