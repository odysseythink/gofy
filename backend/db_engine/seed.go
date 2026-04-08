package dbengine

import (
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

// Seed populates initial data if tables are empty.
func Seed() {
	db := Instance().DB

	// Check if setup already exists
	var count int64
	db.Model(&models.DifySetup{}).Count(&count)
	if count > 0 {
		return
	}

	mlog.Info("seeding initial data...")

	// No initial data needed — setup wizard handles first-time config
	mlog.Info("database seeding completed")
}
