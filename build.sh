protoc --proto_path ../../../ -I=./proto --go_out=plugins=grpc:./proto proto/recordadder.proto
mv proto/github.com/brotherlogic/recordadder/proto/* ./proto
