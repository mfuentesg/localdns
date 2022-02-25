all: dns-pb

dns-pb:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pb/dns.proto

.PHONY: clean
clean:
	@rm pb/dns*.pb.go