package mysql

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"notification-api/config"
)

type DBRepository struct {
	NotificationDB *gorm.DB
}

var AppDB *DBRepository

func NewDatabase() *DBRepository {
	var DBConnection DBRepository
	return &DBConnection
}

func (r *DBRepository) ConnectNotificationDB() *DBRepository {
	notificationDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Config.GetString("MYSQL_NOTIFICATION_USER"), config.Config.GetString("MYSQL_NOTIFICATION_PASS"), config.Config.GetString("MYSQL_NOTIFICATION_HOST"), config.Config.GetString("MYSQL_NOTIFICATION_PORT"), config.Config.GetString("MYSQL_NOTIFICATION_DBNAME"))
	notification, err := gorm.Open(mysql.Open(notificationDSN), &gorm.Config{Logger: logger.Default.LogMode(logger.Info), SkipDefaultTransaction: true})
	if err != nil {
		log.Println("Cannot Connect to Database")
	}
	r.NotificationDB = notification
	return r
}

func (r *DBRepository) Execute() *DBRepository {
	AppDB = r
	return r
}

func (r *DBRepository) FindAll(db *gorm.DB, ctx context.Context, i interface{}, sql string, params ...interface{}) error {
	sql = " /*FORCE_SLAVE*/ " + sql
	db_ := db.WithContext(ctx)
	err := db_.Raw(sql, params...).Find(i).Error
	fmt.Println(i)
	return err
}

func (r *DBRepository) FindOne(db *gorm.DB, ctx context.Context, i interface{}, sql string, params ...interface{}) error {
	sql = " /*FORCE_SLAVE*/ " + sql
	db_ := db.WithContext(ctx)
	err := db_.Raw(sql, params...).First(i).Error
	fmt.Println(i)
	return err
}

func (r *DBRepository) Insert(db *gorm.DB, ctx context.Context, i interface{}) error {
	db_ := db.WithContext(ctx)
	err := db_.Create(i).Error
	fmt.Println(i)
	return err
}

func (r *DBRepository) Update(db *gorm.DB, ctx context.Context, i interface{}) error {
	db_ := db.WithContext(ctx)
	err := db_.Updates(i).Error
	fmt.Println(i)
	return err
}
