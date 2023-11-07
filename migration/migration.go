package migration

import (
	"github.com/onetooler/bistory-backend/container"
	"github.com/onetooler/bistory-backend/infrastructure"
	"github.com/onetooler/bistory-backend/model"
)

func Init(container container.Container) {
	if container.GetConfig().Database.Migration {
		createDatabase(container.GetRepository())
	}
	if container.GetConfig().Extension.MasterGenerator {
		createMasterData(container.GetRepository())
	}
}

func createDatabase(db infrastructure.Repository) {
	_ = db.DropTableIfExists(&model.Account{})
	_ = db.AutoMigrate(&model.Account{})
}

func createMasterData(db infrastructure.Repository) {
	adminAccount, _ := model.NewAccountWithPasswordEncrypt("test", "test@example.com", "test", model.AuthorityAdmin)
	db.Create(adminAccount)
}
