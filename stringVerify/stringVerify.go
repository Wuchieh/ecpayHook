package stringVerify

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type BannedWords struct {
	Words []string `json:"bannedWords"`
}

var bannedWordsList []string

func init() {
	// 读取 banText.json
	banText, err := os.ReadFile("banText.json")
	if err != nil {
		log.Fatal(err)
	}

	// 解析 JSON
	var bannedWords BannedWords
	if err := json.Unmarshal(banText, &bannedWords); err != nil {
		log.Fatal(err)
	}
	bannedWordsList = bannedWords.Words
}

// StringVerify 輸出false代表有包含敏感字詞
// 輸出true代表沒有包含敏感字詞，通過了驗證
func StringVerify(input ...string) bool {
	// 檢查是否包含敏感字
	for _, bannedWord := range bannedWordsList {
		for _, s := range input {
			if strings.Contains(s, bannedWord) {
				return false
			}
		}
	}
	return true
}

func _() {
	input := "這是一條訊息啦。"
	if StringVerify(input) {
		log.Println("不包含敏感詞彙")
	} else {
		log.Println("包含敏感詞彙")
	}
}
