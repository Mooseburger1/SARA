#See https://www.npmjs.com/package/ts-protoc-gen for docs on protoc-gen-ts

PROTOC_GEN_TS_PATH="${SARAPATH}/proto/node_modules/.bin/protoc-gen-ts"
OUT_DIR="${SARAPATH}/frontend"
protoc --plugin="protoc-gen-ts=${PROTOC_GEN_TS_PATH}" --js_out="import_style=commonjs,binary:${OUT_DIR}" --ts_out="service=grpc-web:${OUT_DIR}" --proto_path=$SARAPATH $SARAPATH/proto/sara.proto