rm bin/elaina
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOEXE=linux go build -o bin/elaina
scp bin/elaina deepin:/root/myproj/elaina/bin/
scp Dockerfile deepin:/root/myproj/elaina/