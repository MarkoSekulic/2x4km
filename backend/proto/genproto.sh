protoc\
  --go_out=./template-service --go_opt=paths=source_relative \
  --go-grpc_out=./template-service --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=./template-service \
  --grpc-gateway_opt=paths=source_relative \
  template-service.proto
protoc\
  --go_out=./users-service --go_opt=paths=source_relative \
  --go-grpc_out=./users-service --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=./users-service \
  --grpc-gateway_opt=paths=source_relative \
  users-service.proto