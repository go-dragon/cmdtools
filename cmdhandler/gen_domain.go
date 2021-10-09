package cmdhandler

import (
	"cmdtools/cmdhandler/tabledao"
	"cmdtools/core/dragon/conf"
	"cmdtools/tools"
	"log"
	"os"
	"strings"
)

const (
	entityDir     = "../domain/entity"
	mapperDir     = "../domain/mapper"
	repositoryDir = "../domain/repository"
	serviceDir    = "../domain/service"
)

// GenDomain gen domain files
// args : ./bin/dragon gen domain [table name]
func GenDomain(args []string) {
	tables := args[3:]

	// iterate tableNames, generate entity,repo,service
	for _, table := range tables {
		// analyze table columns
		tableFields := tabledao.DescTable(table)
		// gen all files
		// gen entity
		GenEntity(table, tableFields)
		// gen mapper sql
		GenMapper(table)
		// gen repository
		GenRepository(table)
		// gen service
		GenService(table)
	}
}

// GenEntity 生成实体
func GenEntity(table string, tableFields []tabledao.TableDesc) {
	// 判断文件是否存在，存在则不处理

	entityName := tools.CamelString(table) + "Entity"
	content := "package entity\n\n" +
		"type " + entityName + " struct {\n"
	// gen struct fields
	var structFields string
	for _, tableField := range tableFields {
		structFields += "	" + tools.CamelString(tableField.Field) + " " + tableField.Type
		if tableField.Extra == "auto_increment" {
			structFields += "`gorm:\"primaryKey;AUTO_INCREMENT\"`"
		}
		structFields += "\n"
	}
	structFields += "} \n"
	content += structFields + `
// TableName set orm table name
func (` + entityName + `) TableName() string {
	return "` + table + `"
}`

	entityFileName := conf.ExecDir + entityDir + "/" + table + "_entity.go"
	CreateFileIfNotExist(entityFileName, content)
}

// GenMapper 生成mapper，sql映射
func GenMapper(table string) {
	// 判断文件是否存在，存在则不处理
	mapperDirPrefixName := strings.ReplaceAll(table, "_", "")
	content := "package " + mapperDirPrefixName + "mapper\n\n" +
		"const (\n\tGetOne = `select * from " + table + " limit 1`\n)\n"
	mapperDirName := conf.ExecDir + mapperDir + "/" + mapperDirPrefixName + "mapper"
	os.Mkdir(mapperDirName, os.ModePerm) // 创建特定table mapper目录
	mapperFileName := mapperDirName + "/" + table + "_mapper.go"
	CreateFileIfNotExist(mapperFileName, content)
}

// GenRepository
func GenRepository(table string) {
	// 判断文件是否存在，存在则不处理
	tableCamelName := tools.CamelString(table)
	mapperDirPrefixName := strings.ReplaceAll(table, "_", "")
	content := "package repository\n\n" +
		"import (\n\t\"dragon/domain/entity\"\n\t\"dragon/domain/mapper/" + mapperDirPrefixName + "mapper" + "\"\n\t\"gorm.io/gorm\"\n)\n\n" +
		"type I" + tableCamelName + "Repository interface {GetOne() (*entity." + tableCamelName + "Entity" + ", error)}\n" +
		"type " + tableCamelName + "Repository struct {\n\tMysqlDB *gorm.DB\n}\n" +
		"func New" + tableCamelName + "Repository(db *gorm.DB) I" + tableCamelName + "Repository {\n\treturn &" + tableCamelName + "Repository{\n\t\tMysqlDB: db,\n\t}\n}\n" +
		"func (this *" + tableCamelName + "Repository" + ") GetOne() (*entity." + tableCamelName + "Entity, error) {\n\tvar data entity." + tableCamelName + "Entity\n\tres := this.MysqlDB.Raw(" + mapperDirPrefixName + "mapper.GetOne).Scan(&data)\n\treturn &data, res.Error\n}"
	repoFileName := conf.ExecDir + repositoryDir + "/" + table + "_repository.go"
	CreateFileIfNotExist(repoFileName, content)
}

// GenService 生成服务
func GenService(table string) {
	// 判断文件是否存在，存在则不处理
	tableCamelName := tools.CamelString(table)
	content := "package service\n\nimport (\n\t\"dragon/domain/entity\"\n\t\"dragon/domain/repository\"\n\t\"gorm.io/gorm\"\n)\n" +
		"type I" + tableCamelName + "Service interface {GetOne() (*entity." + tableCamelName + "Entity" + ", error)}\n" +
		"type " + tableCamelName + "Service struct {\n\t\t" + tableCamelName + "Repository repository.I" + tableCamelName + "Repository\n}\n" +
		"func New" + tableCamelName + "Service(db *gorm.DB) I" + tableCamelName + "Service {\n\treturn &" + tableCamelName + "Service{\n\t\t" + tableCamelName + "Repository: repository.New" + tableCamelName + "Repository(db),\n\t}\n}\n" +
		"func (this *" + tableCamelName + "Service" + ") GetOne() (*entity." + tableCamelName + "Entity, error) {\n\tvar data *entity." + tableCamelName + "Entity\n\tdata,err := this." + tableCamelName + "Repository.GetOne()\n\treturn data, err\n}"

	repoFileName := conf.ExecDir + serviceDir + "/" + table + "_service.go"
	CreateFileIfNotExist(repoFileName, content)
}

// CreateFileIfNotExist 如果文件不存在，则创建文件
func CreateFileIfNotExist(fileName string, content string) {
	fileName = conf.FmtSlash(fileName) // 区分windows系统
	fileInfo, err := os.Stat(fileName)
	if fileInfo != nil {
		// 文件存在，不改变
		return
	}
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Fatalln(err)
	}
	file.WriteString(content)
}
