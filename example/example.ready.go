package example

func PreReady(config *AppConfig) error {
	// Executar as Migrations
	// if err := config.GormStore.AutoMigrate(); err != nil {
	// 	return err
	// }
	// Executar as Seeds
	// if err := config.SeedPermissions(&Permissions); err != nil {
	// 	return err
	// }
	return nil
}

func PosReady(config *AppConfig) error {
	return nil
}
