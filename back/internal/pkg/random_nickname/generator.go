package random_nickname

import (
	"fmt"
	"math/rand"
	"time"
)

var nameList = []string{
	"星河", "清风", "白夜", "落尘", "孤影",
	"长歌", "惊鸿", "浮生", "青山", "流云",
}

// GenerateNickname 生成随机昵称：名字且数字
func GenerateNickname() string {
	rand.Seed(time.Now().UnixNano())

	// 随机名字
	name := nameList[rand.Intn(len(nameList))]

	// 随机4位数字
	number := rand.Intn(90) + 10

	return fmt.Sprintf("%s且%d", name, number)
}
