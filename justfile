default:
    just fmt

fmt:
    go fmt

run:
    go run .

build:
    go build -o bin/etnograbber

buildwin:
    GOOS=windows GOARCH=amd64 go build -o bin/etnograbber.exe

buildserver:
    GOOS=linux GOARCH=amd64 go build -o bin/linux/etnograbber

buildall:
   just build
   just buildwin
   just buildserver
