package go_micro_srv_user

import(
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

// Set the id as a uuid.NewV4()
func (model *User) BeforeCreate(scope *gorm.Scope) error {
	newUuid := uuid.NewV4()
	return scope.SetColumn("Id", newUuid.String())
}