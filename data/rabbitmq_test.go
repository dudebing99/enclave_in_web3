package data

import (
	"fmt"
	"testing"
)

func TestPubSub(test *testing.T) {
	// TODO: 初始化配置文件
	mq := MustGetRabbitMq("zany_swap")

	mq.PublishJson("hello1")
	mq.PublishJson("hello2")
	mq.PublishJson("hello3")

	channel, _ := mq.Consume()
	for msg := range channel {
		fmt.Println(string(msg.Body))
	}
}
