PROTO_DIR=proto
GEN_DIR=gen/go

.PHONY: all proto crdb obs clean format

all: proto crdb obs

proto:
	@echo "🚀 Generating protobuf files..."
	buf generate --template buf.gen.yaml

crdb: proto
	@echo "🔨 Building CRDB..."
	cd CRDB && go build -o ../bin/crdb

obs: proto
	@echo "🔨 Building OBS..."
	cd OBS && go build -o ../bin/obs

run-crdb: crdb
	@echo "🏃‍♂️ Running CRDB..."
	./bin/crdb

run-obs: obs
	@echo "🏃‍♂️ Running OBS..."
	./bin/obs

clean:
	@echo "🧹 Cleaning up..."
	rm -rf $(GEN_DIR)
	rm -rf bin/crdb bin/obs

fmt:
	@echo "🎨 Formatting Go files..."
	find . -name '*.go' -not -path './gen/*' -exec gofmt -s -w {} +
