package utils

import (
	"fmt"
	"time"
)

func GetRedisKey() string {
	today := time.Now().Format("2006-01-02")
	return fmt.Sprintf("datatest_ingestion_count:%v", today)
}
