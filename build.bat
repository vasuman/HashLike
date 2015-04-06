@echo off

mkdir out
mkdir out\static

call go build -o out\server.exe .

xcopy .\client\hashlike.js .\out\static\ /Y

call emcc client\hc_worker.c --pre-js client\hc-worker.js -o out\static\hc-worker-emc.js -s EXPORTED_FUNCTIONS="['_get_nonce']"




