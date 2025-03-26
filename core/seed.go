package core

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"gorm.io/gorm"
)

func (s *AppConfig) SeedUserAdmin() error {
	var user User
	err := s.GormStore.Where("username = ?", s.Super.SuperUser).First(&user).Error
	if err == nil {
		return fmt.Errorf("admin already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check admin existence: %s", err.Error())
	}
	hashPassword, err := HashPassword(s.Super.SuperPass)
	if err != nil {
		return fmt.Errorf("failed to create admin: %s", err.Error())
	}
	admin := &User{
		FirstName:   s.Super.SuperName,
		LastName:    "Admin",
		Username:    s.Super.SuperUser,
		Email:       s.Super.SuperEmail,
		Password:    hashPassword,
		Active:      true,
		IsSuperUser: true,
		Phone1:      s.Super.SuperPhone,
	}
	if err := s.GormStore.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create user: %s", err.Error())
	}
	log.Println("admin created successfully")
	return nil
}

func (s *AppConfig) SeedPermissions(permissions any) error {
	v := reflect.ValueOf(permissions)
	if v.Kind() == reflect.Ptr {
		v = v.Elem() // Dereferencia o ponteiro para obter o valor subjacente
	}

	// Verifica se o valor subjacente é uma struct
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, got %v", v.Kind())
	}
	t := v.Type()
	for i := range t.NumField() {
		field := t.Field(i)
		valueTag := field.Tag.Get("value")
		descriptionTag := field.Tag.Get("description")
		var item Permission
		if valueTag != "" && field.Type.Kind() == reflect.String {
			v.Field(i).SetString(valueTag)
			err := s.GormStore.Where("code = ?", valueTag).First(&item).Error
			if err == nil {
				log.Printf("permission with code '%s' already exists \n", v.Field(i).String())
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed to check permission existence: %s", err.Error())
			}

			permission := &Permission{
				Name:        t.Field(i).Name,
				Code:        valueTag,
				Description: descriptionTag,
			}
			if err := s.GormStore.Create(&permission).Error; err != nil {
				return fmt.Errorf("failed to create permission: %s", err.Error())
			}
		}
	}
	log.Println("permissions created successfully")
	return nil
}
