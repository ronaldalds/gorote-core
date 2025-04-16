package example

func (config *AppConfig) PreReady() error {
	// Executar as Migrations
	// if err := config.GormStore.AutoMigrate(); err != nil {
	// 	return err
	// }
	// Executar as Seeds
	// if err := config.SavePermissions(
	// 	PermissionExampleCreate,
	// 	PermissionExampleView,
	// 	PermissionExampleUpdate,
	// ); err != nil {
	// 	return err
	// }
	return nil
}

func (s *Service) PosReady() error {
	return nil
}
