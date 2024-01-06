package database

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/tymbaca/stackexchange-analyzer/puller"
)

type DB struct {
	clickhouse.Conn
}

func New() *DB {
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

	return &DB{conn}
}

func (db *DB) PushQuestions(questions []puller.Question) error {
	ctx := context.Background()
	batch, err := db.PrepareBatch(ctx, "INSERT INTO questions")
	if err != nil {
		return err
	}

	for _, q := range questions {
		err := batch.Append(
			q.Tags,
			q.Owner.UserId,
			q.IsAnswered,
			q.ViewCount,
			q.AnswerCount,
			q.Score,
			q.LastActivityDate,
			q.CreationDate,
			q.QuestionId,
			q.ContentLicense,
			q.Link,
			q.Title,
		)
		if err != nil {
			return err
		}
	}

	err = batch.Send()
	if err != nil {
		return err
	}
	return nil

}
