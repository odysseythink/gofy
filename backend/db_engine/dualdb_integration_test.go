package dbengine

import (
	"os"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Dual-DB integration test. Skipped when DSN env vars are not set.
//
// Env vars (set by docker-compose.test.yml workflow):
//
//	GOFY_TEST_MYSQL_DSN — e.g. "root:test@tcp(127.0.0.1:33306)/gofy?charset=utf8mb4&parseTime=True&loc=Local"
//	GOFY_TEST_PG_DSN    — e.g. "host=127.0.0.1 port=55432 user=postgres password=test dbname=gofy sslmode=disable"
//
// Exercises the critical bool/tinyint round-trip that motivated the refactor.

type boolProbe struct {
	ID      string `gorm:"column:id;primaryKey;size:36"`
	Enabled bool   `gorm:"column:enabled;not null;default:false"`
}

func (boolProbe) TableName() string { return "gofy_dualdb_bool_probe" }

func runBoolRoundTrip(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.Migrator().DropTable(&boolProbe{}); err != nil {
		t.Fatalf("drop: %v", err)
	}
	if err := db.AutoMigrate(&boolProbe{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	defer db.Migrator().DropTable(&boolProbe{})

	rows := []boolProbe{{ID: "a", Enabled: true}, {ID: "b", Enabled: false}}
	if err := db.Create(&rows).Error; err != nil {
		t.Fatalf("insert: %v", err)
	}

	var got []boolProbe
	if err := db.Order("id").Find(&got).Error; err != nil {
		t.Fatalf("select: %v", err)
	}
	if len(got) != 2 || !got[0].Enabled || got[1].Enabled {
		t.Fatalf("round-trip mismatch: %+v", got)
	}

	var trueCount int64
	if err := db.Model(&boolProbe{}).Where("enabled = ?", true).Count(&trueCount).Error; err != nil {
		t.Fatalf("where bool: %v", err)
	}
	if trueCount != 1 {
		t.Fatalf("expected 1 true row, got %d", trueCount)
	}
}

func TestBoolRoundTripMySQL(t *testing.T) {
	dsn := os.Getenv("GOFY_TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("GOFY_TEST_MYSQL_DSN not set")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open mysql: %v", err)
	}
	runBoolRoundTrip(t, db)
}

func TestBoolRoundTripPostgres(t *testing.T) {
	dsn := os.Getenv("GOFY_TEST_PG_DSN")
	if dsn == "" {
		t.Skip("GOFY_TEST_PG_DSN not set")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	runBoolRoundTrip(t, db)
}
