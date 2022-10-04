package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
)

func Migrations() []*gormigrate.Migration {
	var migrations = []*gormigrate.Migration{
		m202209070900_initial_schema(),
		m202209251530_add_autoallowips_column(),
		m202209251532_add_alias_column(),
		m202229251530_add_alias_column_constraint(),
		m202210040828_add_derpmap_colum(),
	}
	return migrations
}
