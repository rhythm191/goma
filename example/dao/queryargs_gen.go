package dao

// NOTE: THIS FILE WAS PRODUCED BY THE
// GOMA CODE GENERATION TOOL (github.com/kyokomi/goma)
// DO NOT EDIT

import (
	"path/filepath"

	"github.com/kyokomi/goma"
)

type queryArgSettings struct {
	rootDir string
	sqlFile bool
}

var settings queryArgSettings

// SetupQueryArgs setup query file path
func SetupQueryArgs(rootDir string, sqlFile bool) {
	settings.rootDir = rootDir
	settings.sqlFile = sqlFile
}

// GenerateQuery exported queryArgs
func GenerateQuery(tableName string, queryName string, args map[string]interface{}) (string, []interface{}, error) {
	return queryArgs(tableName, queryName, args)
}

func queryArgs(tableName string, queryName string, args map[string]interface{}) (string, []interface{}, error) {
	return settings.queryArgs(tableName, queryName, args)
}

func (s queryArgSettings) queryArgs(tableName string, queryName string, args map[string]interface{}) (string, []interface{}, error) {
	filePath := createSqlFilePath(s.rootDir, tableName, queryName)
	return goma.MySQLGenerateQuery(assetSQL(filePath), args)
}

func assetSQL(filePath string) string {
	var data []byte
	var err error
	if settings.sqlFile {
		data, err = AssetFile(filePath)
	} else {
		data, err = Asset(filePath)
	}
	if err != nil {
		// Asset was not found.
	}
	return string(data)
}

func createSqlFilePath(rootDir string, tableName string, queryName string) string {
	return filepath.Join(rootDir, "sql", tableName, queryName+".sql")
}
