rsrc -manifest main.manifest -ico ./logo.ico -o main.syso

go build -ldflags="-H windowsgui -s -w"
upx translator.exe
