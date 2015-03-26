#!/bin/bash

mkdir -p out/static

go build -o out/server.exe .

cp client/*.js client/*.html out/static/
cp -r templates out


