.PHONY: proto

location := C:\Go\src
proto:
	protoc \
		-I=./proto \
		-I=${GOPATH}\pkg\mod\github.com\grpc-ecosystem\grpc-gateway@v1.16.0\third_party\googleapis \
		--go_out=plugins=grpc:. \
		--grpc-gateway_out=. \
		--swagger_out=logtostderr=true:./doc ./proto/*.proto