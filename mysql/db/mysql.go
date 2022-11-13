package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"sync"
	gl "rabbit/globals"
)

type MysqlPool struct {
	db  *gorm.DB
	dns string
}

func (mp *MysqlPool) GetDb() *gorm.DB {
	return mp.db
}

var once sync.Once
var instance *MysqlPool

func GetDb() *gorm.DB {
	return Instance().GetDb()
}

func Instance() *MysqlPool {
	once.Do(func() {
		instance = &MysqlPool{}
		instance.initMysql()
	})
	return instance
}

func (mp *MysqlPool) initMysql() {
	dsn := gl.V.GetString("mysql.username") + ":" +
		gl.V.GetString("mysql.password") + "@tcp(" +
		gl.V.GetString("mysql.path") + ":" +
		gl.V.GetString("mysql.port") + ")/" +
		gl.V.GetString("mysql.db-name") + "?" +
		gl.V.GetString("mysql.config")

	mysqlConfig := mysql.Config{
		DriverName:                "mysql",
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         255,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}

	//这一段才是真正去连接数据库
	db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig(gl.V.GetBool("mysql.log-mode")))
	if err != nil {
		fmt.Printf("数据库连接失败: %s\n\n\n", err.Error())
	} else {
		mp.db = db
		s, _ := mp.db.DB()
		s.SetMaxIdleConns(20)
		s.SetMaxOpenConns(50)
		fmt.Printf("============================数据库连接成功Ok\n\n\n")
	}

}

// gormConfig 根据配置决定是否开启日志
func gormConfig(mod bool) *gorm.Config {
	if mod {
		return &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Info),
			DisableForeignKeyConstraintWhenMigrating: true,
			//配置表面后面不要带s
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "",
				SingularTable: true,
				NameReplacer:  nil,
			},
		}
	} else {
		return &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Silent),
			DisableForeignKeyConstraintWhenMigrating: true,

			//配置表面后面不要带s
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "",
				SingularTable: true,
				NameReplacer:  nil,
			},
		}
	}

}
