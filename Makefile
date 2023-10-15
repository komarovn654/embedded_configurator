PROJECTNAME=$("embedded configurator")

GOTESTPACKAGES=./config ./config/common_paths ./config/pll_config ./generator/pll_generator ./targets/stm32 ./targets/stm32/stm32_target_pll ./utils/config_parser

GOFILES=$(wildcard *.go)

MAKEFLAGS += silent

compile: clean
	go build $(GOFILES)

run:
	go run $(GOFILES) -config config -config_path .
	
test:
	go test $(GOTESTPACKAGES)

clean:
	go clean
	rm *.h
