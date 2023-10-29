package migration

import (
	"github.com/onetooler/bistory-backend/container"
	"github.com/onetooler/bistory-backend/model"
)

// InitMasterData creates the master data used in this application.
func InitMasterData(container container.Container) {
	if container.GetConfig().Extension.MasterGenerator {
		repp := container.GetRepository()

		adminAccount, _ := model.NewAccountWithPasswordEncrypt("test", "test@example.com", "test", model.AuthorityAdmin)
		repp.Create(adminAccount)
	}
}
