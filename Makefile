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

rev-version: VERSION ?= "1.0.3"
rev-version:
	@echo "🔧 Setting version to $(VERSION) in package.json"
	sed -i '' 's/"version": "[^"]*"/"version": $(VERSION)/' ui/package.json
	@echo "🔨 Building the React app..."
	cd ui && npm run build
	@echo "📦 Cleaning the previous obsbundle"
	rm -rf obsbundle
	@echo "📦 Copying built files to bundles/v$(VERSION)"
	mkdir -p bundles/v$(VERSION)
	cp -r ui/dist/* bundles/v$(VERSION)/
	@echo "📦 Creating zip archive of the built files"
	cd bundles && zip -r v$(VERSION).zip v$(VERSION)
	@echo "📦 Cleaning up..."
	rm -rf bundles/v$(VERSION)

