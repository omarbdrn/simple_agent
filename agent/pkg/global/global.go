package global

import (
	"github.com/omarbdrn/simple_agent/pkg/models"
)

var current_share models.Share

func SetCurrentShare(share models.Share) {
	current_share = share
}

func GetCurrentShare() models.Share {
	return current_share
}

func IsInArray(target string, array []string) bool {
	for _, element := range array {
		if element == target {
			return true
		}
	}
	return false
}
