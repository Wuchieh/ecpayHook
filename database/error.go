package database

import (
	"fmt"
	"github.com/lib/pq"
	"log"
	"regexp"
)

func pgError(err error) error {
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return err
	}
	if pqErr.Code == "23505" { // 唯一值衝突錯誤代碼
		match := regexp.MustCompile(`Key \((.*?)\)=\((.*?)\) already exists.`).FindStringSubmatch(pqErr.Detail)
		if len(match) == 3 && match[1] == "username" { // 提取冲突的列和值
			return fmt.Errorf("usernameAlreadyExists")
		} else if len(match) == 3 && match[1] == "\"lineID\"" {
			return fmt.Errorf("lineForLineError")
		} else {
			log.Println("pqErr.Detail", pqErr.Detail, match)
		}
	} else {
		log.Println("例外錯誤", pqErr.Code)
	}
	return err
}
