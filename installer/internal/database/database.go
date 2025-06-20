package database

import (
	"database/sql"
	"fmt"
	"time"

	"context"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSL      bool   `json:"ssl"`
}

type Manager interface {
	TestConnection(config Config) error
	Setup(config Config) error
	CreateDatabase(config Config) error
	CreateTables(config Config) error
}

type manager struct{}

func New() Manager {
	return &manager{}
}

func (m *manager) TestConnection(config Config) error {
	switch config.Type {
	case "mysql":
		return m.testMySQLConnection(config)
	case "mongodb":
		return m.testMongoDBConnection(config)
	default:
		return fmt.Errorf("不支持的数据库类型: %s", config.Type)
	}
}

func (m *manager) Setup(config Config) error {
	// 测试连接
	if err := m.TestConnection(config); err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	// 创建数据库
	if err := m.CreateDatabase(config); err != nil {
		return fmt.Errorf("创建数据库失败: %w", err)
	}

	// 创建表/集合
	if err := m.CreateTables(config); err != nil {
		return fmt.Errorf("创建表失败: %w", err)
	}

	return nil
}

func (m *manager) CreateDatabase(config Config) error {
	switch config.Type {
	case "mysql":
		return m.createMySQLDatabase(config)
	case "mongodb":
		return m.createMongoDatabase(config)
	default:
		return fmt.Errorf("不支持的数据库类型: %s", config.Type)
	}
}

func (m *manager) CreateTables(config Config) error {
	switch config.Type {
	case "mysql":
		return m.createMySQLTables(config)
	case "mongodb":
		return m.createMongoCollections(config)
	default:
		return fmt.Errorf("不支持的数据库类型: %s", config.Type)
	}
}

// MySQL相关方法
func (m *manager) testMySQLConnection(config Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", config.Username, config.Password, config.Host, config.Port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}

func (m *manager) createMySQLDatabase(config Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", config.Username, config.Password, config.Host, config.Port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", config.Database)
	_, err = db.Exec(query)
	return err
}

func (m *manager) createMySQLTables(config Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// 创建用户表
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	// 创建插件表
	pluginTable := `
	CREATE TABLE IF NOT EXISTS plugins (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL,
		version VARCHAR(50) NOT NULL,
		description TEXT,
		author VARCHAR(255),
		enabled BOOLEAN DEFAULT TRUE,
		installed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	// 创建外链表
	externalLinksTable := `
	CREATE TABLE IF NOT EXISTS external_links (
		id INT AUTO_INCREMENT PRIMARY KEY,
		url VARCHAR(500) UNIQUE NOT NULL,
		title VARCHAR(255),
		description TEXT,
		click_count INT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	tables := []string{userTable, pluginTable, externalLinksTable}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return err
		}
	}

	return nil
}

// MongoDB相关方法
func (m *manager) testMongoDBConnection(config Config) error {
	uri := m.buildMongoURI(config)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Ping(ctx, nil)
}

func (m *manager) createMongoDatabase(config Config) error {
	uri := m.buildMongoURI(config)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	// MongoDB会在首次写入时自动创建数据库
	db := client.Database(config.Database)

	// 创建一个临时集合来确保数据库被创建
	_, err = db.Collection("_temp").InsertOne(context.Background(), map[string]interface{}{"init": true})
	if err != nil {
		return err
	}

	// 删除临时集合
	return db.Collection("_temp").Drop(context.Background())
}

func (m *manager) createMongoCollections(config Config) error {
	uri := m.buildMongoURI(config)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	db := client.Database(config.Database)

	// 创建集合和索引
	collections := []struct {
		name    string
		indexes []mongo.IndexModel
	}{
		{
			name: "users",
			indexes: []mongo.IndexModel{
				{Keys: map[string]interface{}{"username": 1}, Options: options.Index().SetUnique(true)},
				{Keys: map[string]interface{}{"email": 1}, Options: options.Index().SetUnique(true)},
			},
		},
		{
			name: "plugins",
			indexes: []mongo.IndexModel{
				{Keys: map[string]interface{}{"name": 1}, Options: options.Index().SetUnique(true)},
				{Keys: map[string]interface{}{"enabled": 1}},
			},
		},
		{
			name: "external_links",
			indexes: []mongo.IndexModel{
				{Keys: map[string]interface{}{"url": 1}, Options: options.Index().SetUnique(true)},
				{Keys: map[string]interface{}{"created_at": -1}},
			},
		},
	}

	for _, coll := range collections {
		// 创建集合
		err := db.CreateCollection(context.Background(), coll.name)
		if err != nil && !mongo.IsDuplicateKeyError(err) {
			return err
		}

		// 创建索引
		if len(coll.indexes) > 0 {
			collection := db.Collection(coll.name)
			_, err := collection.Indexes().CreateMany(context.Background(), coll.indexes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *manager) buildMongoURI(config Config) string {
	auth := ""
	if config.Username != "" && config.Password != "" {
		auth = fmt.Sprintf("%s:%s@", config.Username, config.Password)
	}

	return fmt.Sprintf("mongodb://%s%s:%d/%s", auth, config.Host, config.Port, config.Database)
}
