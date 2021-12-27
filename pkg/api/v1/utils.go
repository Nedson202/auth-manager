package v1

import "github.com/gofrs/uuid"

func getUUID() string {
	u4, _ := uuid.NewV4()
	return u4.String()
}
