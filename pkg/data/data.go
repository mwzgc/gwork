package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
)

type any interface{}

// Data .
type Data struct {
	rdb *redis.Client
	db  *sql.DB
}

// NewData .
func NewData() (*Data, func(), error) {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "test",
		AllowNativePasswords: true,
	}

	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "test",
		// DB:           int(conf.Redis.Db),
		// DialTimeout:  conf.Redis.DialTimeout.AsDuration(),
		// WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		// ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
	})
	rdb.AddHook(redisotel.TracingHook{})
	d := &Data{
		rdb: rdb,
		db:  db,
	}
	return d, func() {
		fmt.Println("data init error")
		if err := d.rdb.Close(); err != nil {

		}
	}, nil
}

func (d *Data) GetDb() sql.DB {
	return *d.db
}

func (d *Data) Set(ctx context.Context, key string, data any) {
	if data == nil {
		return
	}

	byteData, err0 := json.Marshal(data)
	if err0 != nil {
		panic(err0)
	}

	err := d.rdb.Set(ctx, key, string(byteData), 0).Err()
	if err != nil {
		panic(err)
	}
}

func (d *Data) SetWithExpir(ctx context.Context, key string, data any, expiration time.Duration) {
	if data == nil {
		return
	}

	byteData, err0 := json.Marshal(data)
	if err0 != nil {
		panic(err0)
	}

	err := d.rdb.Set(ctx, key, string(byteData), expiration).Err()
	if err != nil {
		panic(err)
	}
}

func (d *Data) Get(ctx context.Context, key string) string {
	val, err := d.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return ""
	} else if err != nil {
		panic(err)
	} else {
		// var newtype []TestObj
		// json.Unmarshal(val, &newtype)
		return val
	}
}
