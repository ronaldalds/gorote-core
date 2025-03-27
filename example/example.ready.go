package example

import (
	"fmt"

	"github.com/ronaldalds/gorote-core/core"
)

func PreReady(config *core.AppConfig) error {
	// Executar as Migrations
	if err := config.GormStore.AutoMigrate(); err != nil {
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

func PosReady(config *core.AppConfig) error {
	return nil
}
