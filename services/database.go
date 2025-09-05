package services

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/f-alotaibi/go-starter/models"
	"github.com/glebarez/sqlite"
	"github.com/go-gorm/caches/v4"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/joho/godotenv/autoload"
)

type DBType string

const (
	DBTypeSQLite     DBType = "sqlite"
	DBTypeMySQL      DBType = "mysql"
	DBTypePostgreSQL DBType = "postgresql"
)

func NewDB() (*gorm.DB, error) {
	dbType := os.Getenv("DATABASE_TYPE")
	var db *gorm.DB
	var err error
	switch DBType(dbType) {
	case (DBTypeSQLite):
		{
			dsn := fmt.Sprintf("%s.db", os.Getenv("DATABASE_DB"))
			db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
			break
		}
	case (DBTypeMySQL):
		{
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_DB"))
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
			break
		}
	case (DBTypePostgreSQL):
		{
			dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_DB"), os.Getenv("DATABASE_PORT"))
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			break
		}
	}
	if db == nil {
		return nil, fmt.Errorf("database: couldn't find database for type: %s", dbType)
	}

	if err != nil {
		return nil, err
	}

	cachePlugin := &caches.Caches{
		Conf: &caches.Config{
			Easer:  true,
			Cacher: &memoryCacher{},
		},
	}
	err = db.Use(cachePlugin)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.User{})

	return db, err
}

// Caching
type memoryCacher struct {
	store *sync.Map
}

func (c *memoryCacher) init() {
	if c.store == nil {
		c.store = &sync.Map{}
	}
}

func (c *memoryCacher) Get(ctx context.Context, key string, q *caches.Query[any]) (*caches.Query[any], error) {
	c.init()
	val, ok := c.store.Load(key)
	if !ok {
		return nil, nil
	}

	if err := q.Unmarshal(val.([]byte)); err != nil {
		return nil, err
	}

	return q, nil
}

func (c *memoryCacher) Store(ctx context.Context, key string, val *caches.Query[any]) error {
	c.init()
	res, err := val.Marshal()
	if err != nil {
		return err
	}

	c.store.Store(key, res)
	return nil
}

func (c *memoryCacher) Invalidate(ctx context.Context) error {
	c.store = &sync.Map{}
	return nil
}
