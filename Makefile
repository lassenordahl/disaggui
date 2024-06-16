PROTO_DIR=proto
GEN_DIR=gen/go

.PHONY: all proto crdb obs bucket clean format

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

bucket:
	@echo "🔨 Building Bucket..."
	cd bucket && go build -o ../bin/bucket main.go

run-crdb: crdb
	@echo "🏃‍♂️ Running CRDB..."
	./bin/crdb

run-obs: obs
	@echo "🍕 Starting OBS DB and Server..."
	./bin/obs

run-bucket: bucket
	@echo "🏃‍♂️ Running Bucket..."
	chmod +x ./bin/bucket
	./bin/bucket

clean:
	@echo "🧹 Cleaning up..."
	rm -rf $(GEN_DIR)
	rm -rf bin/crdb bin/obs

fmt:
	@echo "🎨 Formatting Go files..."
	find . -name '*.go' -not -path './gen/*' -exec gofmt -s -w {} +
