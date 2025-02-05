package migrations

// Migrations should NEVER use types from other packages. Types can change
// and then migrations run on a _new_ database will fail or behave unexpectedly.
// Instead of importing types, always re-create the type in the migration, as
// is done here, even though the same type is defined in pkg/api

import (
	"strings"

	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/internal/kafka/internal/api/dbapi"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const oldIngressPrefix = "elb.mk."

func resetOldIngressControllerRoutes() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20210927170000",
		Migrate: func(tx *gorm.DB) error {
			var kafkas []dbapi.KafkaRequest
			routesUpdate := map[string]interface{}{
				"routes": nil,
			}
			err := tx.Model(&dbapi.KafkaRequest{}).
				Select("id", "routes", "routes_created").
				Where("status = ?", "ready").
				Where("routes_created = ? ", "true").
				Scan(&kafkas).
				Error

			if err != nil {
				return err
			}

			for _, kafka := range kafkas {
				hasOldIngress, err := hasOldIngress(&kafka)
				if err != nil {
					return err
				}
				if !hasOldIngress {
					continue
				}

				err = tx.Model(&kafka).Updates(routesUpdate).Error
				if err != nil {
					return err
				}
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	}
}

func hasOldIngress(kafka *dbapi.KafkaRequest) (bool, error) {
	routes, err := kafka.GetRoutes()
	if err != nil {
		return false, err
	}

	for _, route := range routes {
		if strings.HasPrefix(route.Router, oldIngressPrefix) {
			return true, nil
		}
	}

	return false, nil
}
