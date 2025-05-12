package strutil

import "github.com/lost-woods/go-util/osutil"

var (
	log = osutil.GetLogger()
)

func GetStringPointer(str string) *string {
	return &str
}

func ContainsString(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}

	return false
}
