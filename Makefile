services=user video social interaction
cur_dir=$(CURDIR)
module=github.com/ACaiCat/tiktok-go
idl_dir=$(cur_dir)/idl
cmd_dir=$(cur_dir)/cmd

.PHONY: kitex-gen
kitex-gen:
	@for service in $(services); do \
  		cd $(cur_dir); \
  		mkdir -p $(cmd_dir)/$$service && cd $(cmd_dir)/$$service ;\
		kitex -gen-path ./../../kitex_gen -service $$service $(idl_dir)/rpc/$$service.thrift; \
	done
	@go mod tidy


.PHONY: hz-gen-api
hz-gen-api:
	@for service in $(services); do \
  		hz update -I idl -idl $(idl_dir)/api/$$service.thrift; \
  		echo '[INFO] Code Generation is Done!'; \
	done
	@go mod tidy

.PHONY: fmt
fmt:
	@golangci-lint fmt

.PHONY: lint
lint:
	@golangci-lint run
