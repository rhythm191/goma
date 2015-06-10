package goma

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05.999999"
)

var (
	gomaParseTime = time.UTC
)

// ParseTime will configure goma to convert dates of the given location
func ParseTime(location *time.Location) {
	gomaParseTime = location
}

// Goma is sql.DB access wrapper.
type Goma struct {
	*sql.DB
	options Options
}

// QueryArgs sql query args
type QueryArgs map[string]interface{}

// Open is create goma client.
// - database open
func Open(configPath string) (*sql.DB, error) {
	opts, err := NewOptions(configPath)
	if err != nil {
		return nil, err
	}
	return OpenOptions(opts)
}

// OpenOptions is create goma client.
// - database open
func OpenOptions(options Options) (*sql.DB, error) {
	return sql.Open(options.Driver, options.Source())
}

// Close sql.DB close.
func (d *Goma) Close() error {
	d.debugPrintln("goma close")

	err := d.DB.Close()

	return err
}

// MySQLGenerateQuery generate bind args query
func MySQLGenerateQuery(queryString string, args QueryArgs) string {
	if len(args) <= 0 {
		return queryString
	}

	for key, val := range args {
		re := regexp.MustCompile(`\/\* ` + key + ` \*\/.*`)

		replaceWord := ""
		switch v := val.(type) {
		default:
			switch reflect.TypeOf(val).Kind() {
			case reflect.String:
				replaceWord = "'" + reflect.ValueOf(v).String() + "'"
			}
		case int:
			replaceWord = strconv.Itoa(v)
		case bool:
			if val.(bool) {
				replaceWord = "true"
			} else {
				replaceWord = "false"
			}
		case float32:
			replaceWord = strconv.FormatFloat(float64(v), 'f', 3, 32)
		case float64:
			replaceWord = strconv.FormatFloat(v, 'f', 3, 64)
		case int64:
			replaceWord = strconv.FormatInt(v, 10)
		case string:
			replaceWord = "'" + v + "'"
		case []uint8:
			replaceWord = "'" + string(v) + "'"
			//		case Time:
			//			replaceWord = "'" + val.(Time).Time.Format("15:04:05") + "'"
			//		case Date:
			//			replaceWord = "'" + time.Time(val.(Date)).Format("2006-01-02") + "'"
			//		case Timestamp:
			//			replaceWord = "'" + time.Time(val.(Timestamp)).Format("2006-01-02 15:04:05.999999999") + "'"
			//		case mysql.NullTime:
			//			replaceWord = "'" + val.(mysql.NullTime).Time.Format("2006-01-02 15:04:05.999999999") + "'"
		case time.Time:
			fmt.Println(val)
			if v.IsZero() {
				replaceWord = "'0000-00-00'"
			} else {
				replaceWord = "'" + v.In(gomaParseTime).Format(timeFormat) + "'"
			}
		}
		queryString = re.ReplaceAllString(queryString, replaceWord)
	}

	return queryString
}

// PostgresGenerateQuery generate bind args query
func PostgresGenerateQuery(queryString string, args QueryArgs) string {
	if len(args) <= 0 {
		return queryString
	}

	for key, val := range args {
		re := regexp.MustCompile(`\/\* ` + key + ` \*\/.*`)

		replaceWord := ""
		switch v := val.(type) {
		case int:
			replaceWord = strconv.Itoa(v)
		case bool:
			replaceWord = strconv.FormatBool(v)
		case float32:
			replaceWord = strconv.FormatFloat(float64(v), 'f', 3, 32)
		case float64:
			replaceWord = strconv.FormatFloat(v, 'f', 3, 64)
		case int64:
			replaceWord = strconv.FormatInt(v, 10)
		case string:
			replaceWord = "'" + v + "'"
		case []uint8:
			replaceWord = "'" + string(v) + "'"
			//		case Time:
			//			replaceWord = "'" + val.(Time).Time.Format("15:04:05") + "'"
			//		case Date:
			//			replaceWord = "'" + time.Time(val.(Date)).Format("2006-01-02") + "'"
			//		case Timestamp:
			//			replaceWord = "'" + time.Time(val.(Timestamp)).Format("2006-01-02 15:04:05.999999999") + "'"
			//		case mysql.NullTime:
			//			replaceWord = "'" + val.(mysql.NullTime).Time.Format("2006-01-02 15:04:05.999999999") + "'"
		case time.Time:
			if v.IsZero() {
				replaceWord = "'0000-00-00'"
			} else {
				replaceWord = "'" + v.In(gomaParseTime).Format(timeFormat) + "'"
			}
		}
		queryString = re.ReplaceAllString(queryString, replaceWord)
	}

	return queryString
}

func (d *Goma) debugPrintln(v ...interface{}) {
	if d.options.Debug {
		fmt.Println(v...)
	}
}
