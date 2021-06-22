package common_tool

import (
	"strings"
)

func GetEtcdFmtKey(stringList ...string) string {
	key := ""
	if len(stringList) == 0 {
		return key
	}
	key = strings.Join(stringList, "|")
	return key
}
