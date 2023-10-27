package migration

import (
	"github.com/onetooler/bistory-backend/container"
	"github.com/onetooler/bistory-backend/model"
)

// CreateDatabase creates the tables used in this application.
func CreateDatabase(container container.Container) {
	if container.GetConfig().Database.Migration {
		db := container.GetRepository()

		_ = db.DropTableIfExists(&model.Account{})
		_ = db.DropTableIfExists(&model.Authority{})

		_ = db.AutoMigrate(&model.Account{})
		_ = db.AutoMigrate(&model.Authority{})
	}
}
