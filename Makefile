.PHONY: all update clean test

all: main.zip

main.zip: main
	build-lambda-zip -output main.zip main

main: cmd/rssnotify/main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -ldflags '-d -s -w' -a -tags rssnotify -installsuffix rssnotify -o main cmd/rssnotify/main.go

update: main.zip
	aws lambda update-function-code --function-name rssnotify --zip-file fileb://main.zip

clean:
	rm -f main.zip main

test:
	go test -race -count=1 .
