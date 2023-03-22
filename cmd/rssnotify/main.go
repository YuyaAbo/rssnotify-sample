package main

import (
	"github.com/YuyaAbo/rssnotify-sample"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	//err := rssnotify.Run()
	//if err != nil {
	//	log.Println(err)
	//}
	lambda.Start(rssnotify.Run)
}
