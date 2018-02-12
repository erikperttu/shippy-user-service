package auth

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

// Set the id as a uuid.NewV4()
func (model *User) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	return scope.SetColumn("Id", uuid.String())
}
