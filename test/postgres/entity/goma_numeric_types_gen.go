package entity

import "database/sql"

// NOTE: THIS FILE WAS PRODUCED BY THE
// GOMA CODE GENERATION TOOL (github.com/kyokomi/goma)
// DO NOT EDIT

// GomaNumericTypesEntity is generated goma_numeric_types table.
type GomaNumericTypesEntity struct {
	ID              int64   //`goma:"BIGINT:pk"`
	BoolColumns     bool    //`goma:"BOOL"`
	SmallintColumns int     //`goma:"SMALLINT"`
	IntColumns      int     //`goma:"INTEGER"`
	IntegerColumns  int     //`goma:"INTEGER"`
	SerialColumns   int     //`goma:"INTEGER"`
	DecimalColumns  string  //`goma:"NUMERIC"`
	NumericColumns  string  //`goma:"NUMERIC"`
	FloatColumns    float64 //`goma:"DOUBLE"`
}

// Scan GomaNumericTypesEntity all scan
func (e *GomaNumericTypesEntity) Scan(rows *sql.Rows) error {
	return rows.Scan(&e.ID, &e.BoolColumns, &e.SmallintColumns, &e.IntColumns, &e.IntegerColumns, &e.SerialColumns, &e.DecimalColumns, &e.NumericColumns, &e.FloatColumns)
}
