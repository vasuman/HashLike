@echo off

mkdir out
mkdir out\static

go build -o out\server.exe .

emcc  client\hc_worker.c --pre-js client\hc-worker.js -o out\static\hc-worker-emc.js -s EXPORTED_FUNCTIONS="['_get_nonce']"

copy .\client\hashlike.js .\out\static\hashlike.js



