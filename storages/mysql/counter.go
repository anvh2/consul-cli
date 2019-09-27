package mysql

import (
	pb "github.com/anvh2/consul-cli/grpc-gen/counter"
	"github.com/jinzhu/gorm"

	// include mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// CounterDb ...
type CounterDb struct {
	db *gorm.DB
}

// NewCounterDb -
func NewCounterDb(db *gorm.DB) *CounterDb {
	db.AutoMigrate(&pb.PointData{})

	return &CounterDb{
		db: db,
	}
}

// IncreasePoint -
func (db *CounterDb) IncreasePoint(item *pb.PointData) error {
	return db.db.Create(item).Error
}
