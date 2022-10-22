package email

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"time"
)

type Mail struct {
	ID uuid.UUID `json:"id" db:"id"`
}
type DelayedQueue interface {
	Enqueue(ctx context.Context, mail Mail, runAt int64) error
	GetReadyChannel() <-chan []Mail
	Run()
}

type queue struct {
	ready chan []Mail
	rds   *redis.Client
}

func NewQueue(ctx context.Context, addr, password string, db int) (DelayedQueue, error) {
	rds := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	res := rds.Ping(ctx)
	if res.Err() != nil {
		return nil, fmt.Errorf("can't ping redis: %w", res.Err())
	}

	return &queue{
		ready: make(chan []Mail),
		rds:   rds,
	}, nil
}

func (q *queue) Enqueue(ctx context.Context, mail Mail, runAt int64) error {
	jsonMail, err := json.Marshal(mail)
	if err != nil {
		return fmt.Errorf("can't marshal mail: %w", err)
	}

	_, err = q.rds.ZAdd(ctx, "mails", redis.Z{
		Score:  float64(runAt),
		Member: jsonMail,
	}).Result()
	if err != nil {
		return fmt.Errorf("can't enqueue mail: %w", err)
	}
	return nil
}

func (q *queue) GetReadyChannel() <-chan []Mail {
	return q.ready
}

func (q *queue) Run() {
	for {
		pipe := q.rds.TxPipeline()

		now := time.Now().Unix()
		res := pipe.ZRangeByScoreWithScores(context.Background(), "mails", &redis.ZRangeBy{
			Min: "0",
			Max: fmt.Sprint(now),
		})

		pipe.ZRemRangeByScore(context.Background(), "mails", "0", fmt.Sprint(now))

		_, err := pipe.Exec(context.Background())
		if err != nil {
			fmt.Printf("can't run queue: %v", err)
			time.Sleep(time.Second)
			continue
		}

		result, err := res.Result()
		if err != nil {
			fmt.Printf("can't run queue: %v", err)
			time.Sleep(time.Second)
			continue
		}

		var mails []Mail
		for _, v := range result {
			var mail Mail
			err := json.Unmarshal([]byte(v.Member.(string)), &mail)
			if err != nil {
				fmt.Printf("can't run queue: %v", err)
				time.Sleep(time.Second)
				continue
			}
			mails = append(mails, mail)
		}

		if len(mails) > 0 {
			q.ready <- mails
		}

		time.Sleep(time.Second)
	}
}
