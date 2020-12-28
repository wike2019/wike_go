cd proto

protoc --go_out=plugins=grpc:../services/  ./Models.proto
protoc --go_out=plugins=grpc:../services/  ./User.proto
protoc-go-inject-tag   -input=../services/Models.pb.go
cd ..