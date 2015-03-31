#!/bin/bash

mkdir -p out/static

go build -o out/server.exe .

emcc  client/hc_worker.c --pre-js client/hc-worker.js -o out/static/hc-worker-emc.js -s EXPORTED_FUNCTIONS="['_get_nonce']"
cp client/hashlike.js out/static/



