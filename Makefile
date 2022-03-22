all: dns-pb

dns-pb:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pb/dns.proto

test:
	@go test -short -bench=".*" -benchmem -race -covermode=atomic $$(go list ./... | grep -v pb)

.PHONY: clean
clean:
	@rm pb/dns*.pb.go
