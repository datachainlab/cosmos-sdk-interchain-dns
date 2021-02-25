.PHONY: build
build:
	go build -mod readonly -o build/simd ./simapp/simd

.PHONY: test
test:
	go test -v ./x/...
	go test ./simapp/...

.PHONY: proto-gen
proto-gen:
	@echo "Generating Protobuf files"
	docker run -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen sh ./scripts/protocgen.sh
