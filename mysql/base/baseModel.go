package base

import (
	"gorm.io/gorm"
	"rabbit/mysql/db"
)

type BaseModel struct {
}

func (b *BaseModel) GetDb() *gorm.DB {
	return db.GetDb()
}
