#https://github.com/grpc-ecosystem/grpc-gateway


#install protobuf
mkdir tmp
cd tmp
git clone https://github.com/google/protobuf
cd protobuf
./autogen.sh
./configure
make
make check
sudo make install

#install package
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/golang/protobuf/protoc-gen-go