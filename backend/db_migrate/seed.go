package dbmigrate

import (
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
)

// Seed populates initial data if tables are empty.
func Seed() {
	db := dbengine.Instance().DB

	// Check if setup already exists
	var count int64
	db.Model(&models.GofySetup{}).Count(&count)
	if count > 0 {
		return
	}

	mlog.Info("seeding initial data...")

	// No initial data needed — setup wizard handles first-time config
	mlog.Info("database seeding completed")
}
