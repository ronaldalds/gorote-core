package core

import (
	"fmt"
)

func PreReady(config *AppConfig) error {
	// Executar as Migrations
	if err := config.GormStore.AutoMigrate(&User{}, &Role{}, &Permission{}); err != nil {
		return err
	}
	// Executar as Seeds
	if config.Super != nil {
		if err := config.SeedUserAdmin(); err != nil {
			fmt.Println(err.Error())
		}
	}
	if err := config.SeedPermissions(&Permissions); err != nil {
		fmt.Println(err.Error())
	}
	return nil
}

func PosReady(config *AppConfig) error {
	return nil
}
