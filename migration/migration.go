package migration

import (
	"github.com/onetooler/bistory-backend/container"
	"github.com/onetooler/bistory-backend/model"
	"github.com/onetooler/bistory-backend/repository"
)

func Init(container container.Container) {
	if container.GetConfig().Database.Migration {
		createDatabase(container.GetRepository())
	}
	if container.GetConfig().Extension.MasterGenerator {
		createMasterData(container.GetRepository())
	}
}

func createDatabase(db repository.Repository) {
	_ = db.DropTableIfExists(&model.Account{})
	_ = db.AutoMigrate(&model.Account{})
}

func createMasterData(db repository.Repository) {
	adminAccount, _ := model.NewAccountWithPasswordEncrypt("test", "test@example.com", "test", model.AuthorityAdmin)
	db.Create(adminAccount)
}
