package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"encoding/json"
	"io/ioutil"

	"github.com/kyokomi/goma"
)

var sampleDataMap = map[reflect.Type]string{
	reflect.TypeOf(int(1)):         "1",
	reflect.TypeOf(float32(32.1)):  "32.1",
	reflect.TypeOf(float64(64.1)):  "64.1",
	reflect.TypeOf(int64(64)):      "64",
	reflect.TypeOf(string("1111")): "'1111'",
	reflect.TypeOf([]uint8{}):      "'abcdefghijk'",
	reflect.TypeOf(time.Now()):     "'2006/01/02 13:40:00'",
}

var driverImports = map[string]string{
	"mysql":    `_ "github.com/go-sql-driver/mysql"`,
	"postgres": `_ "github.com/lib/pq"`,
}

func generate(pkg string, opt goma.Options) {
	log.SetFlags(log.Llongfile)

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	opt.CurrentDir = currentDir

	// xorm reverse mysql root:@/test?charset=utf8 templates/goxorm
	orm, err := xorm.NewEngine(opt.Driver, opt.Source())
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	tables, err := orm.DBMetas()
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	helperData := HelperTemplateData{}
	helperData.PkgName = pkg
	helperData.DriverImport = driverImports[opt.Driver]
	helperData.Options = opt.Tuples()

	helperData.DaoImport = opt.DaoImportPath()
	helperData.DaoPkgName = opt.DaoPkgName()

	var daoList []DaoTemplateData

	// sql, dao generate

	daoRootPath := filepath.Join(opt.CurrentDir, opt.DaoRootDir)
	entityRootPath := filepath.Join(opt.CurrentDir, opt.EntityRootDir)
	sqlRootPath := filepath.Join(opt.CurrentDir, opt.SQLRootDir)

	for _, table := range tables {
		// create templateData
		data := newTemplateData(table, opt)

		// dao template
		if err := data.execDaoTemplate(daoRootPath); err != nil {
			log.Fatalln(err)
		}

		// entity template
		if err := data.execEntityTemplate(entityRootPath); err != nil {
			log.Fatalln(err)
		}

		// sql template
		if err := data.Table.execTableTemplate(sqlRootPath); err != nil {
			log.Fatalln(err)
		}

		daoList = append(daoList, data)
	}

	// asset generate
	assetData := AssetTemplateData{}
	assetData.DaoPkgName = opt.DaoPkgName()
	if err := assetData.execAssetTemplate(daoRootPath); err != nil {
		log.Fatalln(err)
	}

	// queryargs generate
	queryArgsData := QueryArgsTemplateData{}
	queryArgsData.DaoPkgName = opt.DaoPkgName()
	queryArgsData.SQLRootDir = opt.SQLRootDir
	if err := queryArgsData.execQueryArgsTemplate(daoRootPath); err != nil {
		log.Fatalln(err)
	}

	// helper generate

	helperData.DaoList = daoList

	if err := helperData.execHelperTemplate(currentDir); err != nil {
		log.Fatalln(err)
	}

	// config generate

	data, err := json.MarshalIndent(opt, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}
	if err := ioutil.WriteFile(opt.ConfigPath(), data, 0644); err != nil {
		log.Fatalln(err)
	}
}

func newTemplateData(table *core.Table, opt goma.Options) DaoTemplateData {
	imports := newImports(table.Columns())
	columns := newColumns(table.Columns())

	data := DaoTemplateData{}
	data.Name = lintName(strings.Title(table.Name) + "Dao")
	data.MemberName = "s" + lintName(strings.Title(table.Name))
	data.EntityName = lintName(strings.Title(table.Name) + "Entity")
	data.DaoPkgName = opt.DaoPkgName()
	data.EntityPkgName = opt.EntityPkgName()
	data.EntityImport = opt.EntityImportPath()
	data.Table = TableTemplateData{
		Name:      table.Name,
		TitleName: lintName(strings.Title(table.Name)),
		Columns:   columns,
	}
	data.Imports = imports.slice()

	return data
}

func newImports(columns []*core.Column) set {
	importsMap := make(set, 0)
	for _, c := range columns {
		typ := core.SQLType2Type(c.SQLType)
		importsMap.add(typ.PkgPath())
	}
	return importsMap
}

func newColumns(columns []*core.Column) []ColumnTemplateData {
	var results []ColumnTemplateData
	for _, c := range columns {

		typ := core.SQLType2Type(c.SQLType)
		typeName := typ.Name()
		if typ.PkgPath() != "" {
			typeName = typ.PkgPath() + "." + typ.Name()
		}

		if typeName == "" {
			typeName = typ.String()
		}

		primaryKey := ""
		if c.IsPrimaryKey {
			primaryKey = ":pk"
		}

		typeLength := ""
		if c.SQLType.DefaultLength > 0 {
			typeLength = fmt.Sprintf("(%d)", c.SQLType.DefaultLength)
		}
		typeDetail := fmt.Sprintf("`goma:\"%s"+typeLength+primaryKey+"\"`", c.SQLType.Name)

		column := ColumnTemplateData{
			Name:         c.Name,
			TitleName:    lintName(strings.Title(c.Name)),
			TypeName:     typeName,
			TypeDetail:   typeDetail,
			IsPrimaryKey: c.IsPrimaryKey,
			Sample:       sampleDataMap[typ],
		}
		results = append(results, column)
	}
	return results
}
