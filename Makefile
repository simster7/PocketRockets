

.PHONY: codegen
codegen:
	protoc \
		--go_out=plugins=grpc:`pwd`/backend \
        --plugin="protoc-gen-ts=./frontend/node_modules/.bin/protoc-gen-ts" \
        --js_out="import_style=commonjs,binary:`pwd`/frontend" \
        --ts_out="service=true:`pwd`/frontend" \
        ./api/v1/apis.proto