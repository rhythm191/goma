package entity

import (
	"database/sql"

	"time"
)

// NOTE: THIS FILE WAS PRODUCED BY THE
// GOMA CODE GENERATION TOOL (github.com/kyokomi/goma)
// DO NOT EDIT

// GomaDateTypes is generated goma_date_types table.
type GomaDateTypes struct {
	ID                 int64      `goma:"size:20:pk"`
	DatetimeColumns    time.Time  `goma:""`
	TimestampColumns   time.Time  `goma:""`
	NilDatetimeColumns *time.Time `goma:""`
}

// Scan GomaDateTypes all scan
func (e *GomaDateTypes) Scan(rows *sql.Rows) error {
	err := rows.Scan(&e.ID, &e.DatetimeColumns, &e.TimestampColumns, &e.NilDatetimeColumns)

	e.DatetimeColumns = e.DatetimeColumns.In(time.Local)

	e.TimestampColumns = e.TimestampColumns.In(time.Local)

	if e.NilDatetimeColumns != nil {
		_NilDatetimeColumns := e.NilDatetimeColumns.In(time.Local)
		e.NilDatetimeColumns = &_NilDatetimeColumns
	}

	return err
}
