package ctl

import (
	"errors"
	"gorm.io/gorm"
)

func (r *Controller) IsExist(err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
