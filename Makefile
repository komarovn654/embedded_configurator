PROJECTNAME=$("embedded configurator")

GOTESTPACKAGES=./config/ ./generator/ ./stm32f4xx/ ./stm32f4xx/pll_config/
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
