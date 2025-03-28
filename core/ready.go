package core

func PreReady(config *AppConfig) error {
	// Executar as Migrations
	if err := config.GormStore.AutoMigrate(&User{}, &Role{}, &Permission{}); err != nil {
		return err
	}
	// Executar as Seeds
	if config.Super != nil {
		if err := config.SeedUserAdmin(); err != nil {
			return err
		}
	}
	if err := config.SeedPermissions(&Permissions); err != nil {
		return err
	}
	return nil
}

func PosReady(config *AppConfig) error {
	return nil
}
