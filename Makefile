services=user video social interaction
cur_dir=$(CURDIR)
module=github.com/ACaiCat/tiktok-go
idl_dir=$(cur_dir)/idl
cmd_dir=$(cur_dir)/cmd

.PHONY: hz-gen-api
hz-gen-api:
	@for service in $(services); do \
  		hz update -I idl -idl $(idl_dir)/$$service/$$service.thrift; \
  		echo '[INFO] Code Generation is Done!'; \
	done
	@go mod tidy

.PHONY: fmt
fmt:
	@golangci-lint fmt

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: test
test:
	@go test ./... -gcflags="all=-N -l"
