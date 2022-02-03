package backstore

import (
	"fmt"
	"github.com/arcnadiven/elaina/models"
	"github.com/arcnadiven/elaina/tracelog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

const (
	// e.g. : "root:123456@tcp(172.16.215.215:3306)/persistent_volume?charset=utf8mb4&parseTime=True&loc=Local"
	dsnTempl = "%s:%s@tcp(%s:%s)/%s?%s"

	// TODO: add these config to start flag later
	defaultMaxOpenConns    = 10
	defaultMaxIdleConns    = 100
	defaultConnMaxLifetime = time.Hour
)

var (
	defaultConnectionArguments = map[string]string{
		"charset":   "utf8mb4",
		"parseTime": "True",
		"loc":       "Local",
	}
)

type ConnConfig struct {
	Username     string
	Password     string
	Host         string
	Port         string
	DataBaseName string
	ConnArgs     map[string]string
}

type SQLClient struct {
	client *gorm.DB
	log    tracelog.BaseLogger
}

func NewSQLClient(bl tracelog.BaseLogger, conf *ConnConfig) (*SQLClient, error) {
	if conf.ConnArgs == nil {
		conf.ConnArgs = defaultConnectionArguments
	}
	argList := []string{}
	for k, v := range conf.ConnArgs {
		argList = append(argList, fmt.Sprintf("%s=%s", k, v))
	}
	dsn := fmt.Sprintf(dsnTempl, conf.Username, conf.Password, conf.Host, conf.Port, conf.DataBaseName, strings.Join(argList, "&"))
	bl.Infoln(dsn)

	dbCli, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{DisableAutomaticPing: false})
	if err != nil {
		bl.Errorln(err)
		return nil, err
	}
	dbConn, err := dbCli.DB()
	if err != nil {
		bl.Errorln(err)
		return nil, err
	}
	dbConn.SetMaxOpenConns(defaultMaxOpenConns)
	dbConn.SetMaxIdleConns(defaultMaxIdleConns)
	dbConn.SetConnMaxLifetime(defaultConnMaxLifetime)
	return &SQLClient{
		client: dbCli,
		log:    bl,
	}, dbCli.AutoMigrate(&models.CSIPersiVol{})
}
