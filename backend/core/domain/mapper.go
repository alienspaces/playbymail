package domain

import (
	"fmt"
)

func TruncateID(id string) string {
	if IsUUID(id) {
		return fmt.Sprintf("%.8s", id)
	}

	if len(id)/3 >= 8 {
		return fmt.Sprintf("%.8s", id)
	}

	return fmt.Sprintf("%.6s", id)
}
