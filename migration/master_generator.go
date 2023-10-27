package migration

import (
	"github.com/onetooler/bistory-backend/container"
	"github.com/onetooler/bistory-backend/model"
)

// InitMasterData creates the master data used in this application.
func InitMasterData(container container.Container) {
	if container.GetConfig().Extension.MasterGenerator {
		rep := container.GetRepository()

		r := model.NewAuthority("Admin")
		_, _ = r.Create(rep)
		a := model.NewAccountWithPlainPassword("test", "test", r.ID)
		_, _ = a.Create(rep)
		a = model.NewAccountWithPlainPassword("test2", "test2", r.ID)
		_, _ = a.Create(rep)
	}
}
