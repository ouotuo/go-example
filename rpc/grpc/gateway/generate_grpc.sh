#protoc -I . proto/echo.proto --go_out=plugins=grpc:proto
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:. \
  proto/echo.proto
