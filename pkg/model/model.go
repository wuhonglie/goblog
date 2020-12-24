package model

import (
    "fmt"
    "goblog/pkg/config"
    "goblog/pkg/logger"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
    var err error
    // 初始化 MySQL 连接信息
    var(
        host = config.GetString("database.mysql.host")
        port = config.GetString("database.mysql.port")
        database = config.GetString("database.mysql.database")
        username = config.GetString("database.mysql.username")
        password = config.GetString("database.mysql.password")
        charset = config.GetString("database.mysql.charset")
    )
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",username,password,host,port,database,charset,true,"Local")
    gormConfig := mysql.New(mysql.Config{
        DSN:dsn,
    })
    var level gormlogger.LogLevel
    if config.GetBool("app.debug"){
        level = gormlogger.Warn
    }else{
        level = gormlogger.Error
    }
    // 准备数据库连接池
    DB, err = gorm.Open(gormConfig, &gorm.Config{
        Logger: gormlogger.Default.LogMode(level),
    })
    logger.LogError(err)
    return DB
}
