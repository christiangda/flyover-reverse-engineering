#!/bin/bash
# protoc --go_out=../pkg/mps/. *.proto

DST_DIR=../.
SRC_DIR=.
protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/ResourceManifest.proto
