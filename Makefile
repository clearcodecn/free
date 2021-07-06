bindata:
	@rm -rf ./assets
	@mkdir -p ./assets/template
	@mkdir -p ./assets/static
	@go-bindata -pkg template -o ./assets/template/template.go  ./template/*
	@go-bindata -pkg static -o ./assets/static/static.go ./static/...

build:
	@rm -rf bin
	@mkdir -p bin
	@go build -o bin/yunsmsv2 cmd/server/main.go
	@go build -o bin/generator cmd/generator/generator.go

ip2region:
	@rm -rf ./cache/ip2region.db
	@wget https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region.db -O ./cache/ip2region.db