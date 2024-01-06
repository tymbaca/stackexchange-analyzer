package db

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func ConnectAndInit() clickhouse.Conn {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "tymbaca",
			Password: "qwerty",
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "an-example-go-client", Version: "0.1"},
			},
		},

		Debugf: func(format string, v ...interface{}) {
			fmt.Printf(format, v)
		},
	})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	err = conn.Exec(
		ctx,
		`CREATE TABLE IF NOT EXISTS questions
                (
                    tags Array(String),
                    owner UInt32,
                    is_answered Bool,
                    view_count UInt32,
                    answer_count UInt32,
                    score UInt32,
                    last_activity_date UInt64,
                    creation_date UInt64,
                    question_id UInt32,
                    content_license String,
                    link String,
                    title String
                )
                ENGINE = MergeTree
                PRIMARY KEY question_id`,
	)
	if err != nil {
		panic(err)
	}

	return conn
}
