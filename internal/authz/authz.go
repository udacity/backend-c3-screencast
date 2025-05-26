package authz

import (
	"fmt"
	"strings"

	"eventure/internal/models"
)

func CanPerform(user models.User, action string, resource string) bool {
	required := fmt.Sprintf("%s:%s", resource, action)

	for _, p := range user.Role.Permissions {
		if strings.EqualFold(p.Name, required) {
			return true
		}
	}
	return false
}
