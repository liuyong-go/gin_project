package modelsgenerate

import (
	"fmt"
	"io"

	"gitee.com/liuyongchina/go-library/libs/ycommon"
	"gorm.io/gorm"

	"os"
	"strings"
)

type Generator struct {
	DB         *gorm.DB
	ModelsPath string
	Schema     string
}

func NewGenerator(db *gorm.DB, path string, schema string) *Generator {
	return &Generator{DB: db, ModelsPath: path, Schema: schema}
}
func (ge *Generator) Genertate(tableNames ...string) {
	tableNamesStr := ""
	for _, name := range tableNames {
		if tableNamesStr != "" {
			tableNamesStr += ","
		}
		tableNamesStr += "'" + name + "'"
	}
	tables := ge.getTables(tableNamesStr) //生成所有表信息
	fmt.Println("table info", tables)
	//tables := getTables("admin_info","video_info") //生成指定表信息，可变参数可传入过个表名
	for _, table := range tables {
		fields := ge.getFields(table.Name)
		ge.generateModel(table, fields)
	}
}

//获取表信息
func (ge *Generator) getTables(tableNames string) []Table {
	fmt.Println("tables", tableNames)
	var tables []Table
	fmt.Println(ge.DB.Statement.Schema)
	dbName := ge.Schema
	if tableNames == "" {
		ge.DB.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='" + dbName + "';").Find(&tables)
	} else {
		ge.DB.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE TABLE_NAME IN (" + tableNames + ") AND table_schema='" + dbName + "';").Find(&tables)
	}
	return tables
}

//获取所有字段信息
func (ge *Generator) getFields(tableName string) []Field {
	var fields []Field
	ge.DB.Raw("show FULL COLUMNS from " + tableName + ";").Find(&fields)
	return fields
}

//生成Model
func (ge *Generator) generateModel(table Table, fields []Field) {
	content := "package models\n\n"
	//表注释
	if len(table.Comment) > 0 {
		content += "// " + table.Comment + "\n"
	}
	content += "type " + ycommon.CamelCase(table.Name) + " struct {\n"
	//生成字段
	var fieldName = ""
	for _, field := range fields {
		if field.Field == "id" {
			fieldName = "ID"
		} else {
			fieldName = ycommon.CamelCase(field.Field)
		}

		fieldJson := getFieldJson(field)
		fieldType := getFiledType(field)
		fieldComment := getFieldComment(field)
		content += "	" + fieldName + " " + fieldType + " `" + fieldJson + "` " + fieldComment + "\n"
	}
	content += "}"
	content += "\n"
	content += "func New" + ycommon.CamelCase(table.Name) + "() (*" + ycommon.CamelCase(table.Name) + ") {\n"
	content += "	return &" + ycommon.CamelCase(table.Name) + "{}\n"
	content += "}"
	content += "\n"
	content += "func (*" + ycommon.CamelCase(table.Name) + ") TableName() string {\n"
	content += "	return \"" + table.Name + "\"\n"
	content += "}"
	filename := ge.ModelsPath + ycommon.CamelCase(table.Name) + ".go"
	var f *os.File
	var err error
	if checkFileIsExist(filename) {
		fmt.Println(ycommon.CamelCase(table.Name) + " 已存在，需删除才能重新生成...")
		return
	} else {
		f, err = os.Create(filename)
		if err != nil {
			panic(err)
		}
	}
	defer f.Close()
	_, err = io.WriteString(f, content)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(ycommon.CamelCase(table.Name) + " 已生成...")
	}
}

//获取字段类型
func getFiledType(field Field) string {
	typeArr := strings.Split(field.Type, "(")

	switch typeArr[0] {
	case "int":
		return "int"
	case "integer":
		return "int"
	case "mediumint":
		return "int"
	case "bit":
		return "int"
	case "year":
		return "int"
	case "smallint":
		return "int"
	case "tinyint":
		return "int"
	case "bigint":
		return "int64"
	case "decimal":
		return "float32"
	case "double":
		return "float32"
	case "float":
		return "float32"
	case "real":
		return "float32"
	case "numeric":
		return "float32"
	case "timestamp":
		return "time.Time"
	case "datetime":
		return "time.Time"
	case "time":
		return "time.Time"
	default:
		return "string"
	}
}

//获取字段json描述
func getFieldJson(field Field) string {
	return `json:"` + field.Field + `"`
}

//获取字段说明
func getFieldComment(field Field) string {
	if len(field.Comment) > 0 {
		return "// " + field.Comment
	}
	return ""
}

//检查文件是否存在
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
