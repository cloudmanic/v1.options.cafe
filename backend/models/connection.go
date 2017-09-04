package models

import (
  "github.com/jinzhu/gorm"
)

type Connection struct {
  Conn *gorm.DB
} 

/* End File */