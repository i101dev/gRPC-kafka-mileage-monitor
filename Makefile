obu:
	@go build -o bin/obu ./obu
	@./bin/obu

gateway:
	@go build -o bin/gateway gateway/main.go
	@./bin/gateway
	
receiver:
	@go build -o bin/receiver ./data_receiver
	@./bin/receiver

calculator:
	@go build -o bin/calculator ./distance_calculator
	@./bin/calculator

agg:
	@go build -o bin/aggregator ./aggregator
	@./bin/aggregator

proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto


.PHONY: obu
.PHONY: agg
.PHONY: gateway
.PHONY: receiver
.PHONY: calculator

.PHONY: proto

# PATH="${PATH}:${HOME}/go/bin"
# docker run --name prometheus -d -p 127.0.0.1:9090:9090 prom/prometheus