generate_grpc_code:
	protoc \
	--go_out=currencyConvertor \
	--go_opt=paths=source_relative \
	--go-grpc_out=currencyConvertor \
	--go-grpc_opt=paths=source_relative \
	currencyConvertor.proto