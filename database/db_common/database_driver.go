package db_common

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBTypeDriver string

const (
	PostgresSQL DBTypeDriver = "postgres"
	MySQL       DBTypeDriver = "mysql"
)

type IDBRelDriver interface {
	GetConfigStr() string
	GetDB() (*gorm.DB, error)
	GetType() DBTypeDriver
}

type PostgresConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
	IsSSL    bool
	Timezone DBTimeZone
}

func (c *PostgresConfig) GetConfigStr() string {
	timeZone := ShangHai
	if len(c.Timezone) > 0 {
		timeZone = c.Timezone
	}
	authStr := ""
	if len(c.User) > 0 && len(c.Password) > 0 {
		authStr = fmt.Sprintf("user=%s password=%s", c.User, c.Password)
	}
	return fmt.Sprintf(`host=%s %s dbname=%s 
port=%d TimeZone=%s`, c.Host, authStr, c.DBName, c.Port, timeZone)
}

func (c *PostgresConfig) GetDB() (*gorm.DB, error) {
	dsn := c.GetConfigStr()
	return gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
}

func (c *PostgresConfig) GetType() DBTypeDriver {
	return PostgresSQL
}

func NewPostgres(config PostgresConfig) IDBRelDriver {
	return &config
}

type MySQLConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Timezone DBTimeZone
	Port     int
}

func (c *MySQLConfig) GetConfigStr() string {
	timeZone := LocalTime
	if len(c.Timezone) > 0 {
		timeZone = c.Timezone
	}
	return fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?
charset=utf8mb4&parseTime=True&loc=%s`, c.User, c.Password, c.Host, c.Port, c.DBName, timeZone)
}

func (c *MySQLConfig) GetDB() (*gorm.DB, error) {
	dsn := c.GetConfigStr()
	return gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{})
}

func (c *MySQLConfig) GetType() DBTypeDriver {
	return MySQL
}

func NewMySQL(config MySQLConfig) IDBRelDriver {
	return &config
}

type MongoDBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type IDBRelationProvider interface {
	GetDB() *gorm.DB
}

type dbRelationProvider struct {
	db     *gorm.DB
	dBType DBTypeDriver
}

func (p *dbRelationProvider) GetDB() *gorm.DB {
	return p.db
}

func NewDBRelationProvider(dbDriver IDBRelDriver) (IDBRelationProvider, func()) {
	db, err := dbDriver.GetDB()
	if err != nil {
		panic(err)
	}

	cleanup := func() {
		dbInstance, _ := db.DB()
		if dbInstance != nil {
			dbInstance.Close()
		}
	}

	return &dbRelationProvider{
		db:     db,
		dBType: dbDriver.GetType(),
	}, cleanup
}

type IDBMongoProvider interface {
	GetDBName() string
	GetClient() *mongo.Client
}

type dbMongoProvider struct {
	client *mongo.Client
	DBName string
}

func (d *dbMongoProvider) GetClient() *mongo.Client {
	return d.client
}

func (d *dbMongoProvider) GetDBName() string {
	return d.DBName
}

func NewDBMongoProvider(config MongoDBConfig) (IDBMongoProvider, func()) {
	var auth string
	if len(config.User) > 0 && len(config.Password) > 0 {
		auth = fmt.Sprintf("%s:%s@", config.User, config.Password)
	}
	connectString := fmt.Sprintf("mongodb://%s%s:%d/%s", auth, config.Host, config.Port, config.DBName)
	client, err := mongo.NewClient(options.Client().ApplyURI(connectString))
	if err != nil {
		panic(err)
	}
	return &dbMongoProvider{
			client: nil,
			DBName: "",
		}, func() {
			client.Disconnect(context.Background())
		}
}
