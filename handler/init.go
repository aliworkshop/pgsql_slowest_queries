package handler

import "gorm.io/gorm"

type handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) Interface {
	h := &handler{db: db}
	return h
}
