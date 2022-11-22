package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/trungnghia250/malo-api/rabitmq/model"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Queue interface {
	Push(routingKey string, data []byte) error
}

type useCase struct {
	coll  *mongo.Collection
	queue Queue
}

func NewUseCase(coll *mongo.Collection, queue Queue) *useCase {
	return &useCase{
		coll:  coll,
		queue: queue,
	}
}

// Process watches changed document from a collection of mongodb. It saves resume
// token to an external database and pushes messages to a queue. If it cannot get
// the resume token from the external database, it will continue watching changed
// document. If it cannot save the resume token to the external database, it will
// continue watching the next changed document. If it cannot push the message of
// a resume token to the queue, it will push the message of current resume token
// to the queue again until success.
func (u *useCase) Process() error {
	ctx := context.Background()
	opts := options.ChangeStream()
	opts.SetFullDocument(options.UpdateLookup) // return the last updated document. Visit docs for more information
	for {
		cs, err := u.coll.Watch(ctx, mongo.Pipeline{}, opts)
		if err != nil {
			fmt.Printf("watch err %v \n", err)
			// resume token stored in redis was removed in the oplog due to no longer update.
			if strings.Contains(err.Error(), "resume point may no longer be in the oplog") {
				// delete key in redis for refreshing new resume token in next loop
				if err != nil {
					fmt.Printf("failed to delete key in redis, err %v ", err)
				}
				fmt.Println("deleted redis key")
				continue
			}
		}

		for cs.Next(ctx) {
			changeStream := new(model.Mgostream)
			err = cs.Decode(changeStream)
			if err != nil {
				fmt.Printf("err %v \n", err)
				return err
			}

			var bytes []byte
			bytes, err = json.Marshal(changeStream)
			if err != nil {
				fmt.Printf("%v \n", err)
				return err
			}

			// push message to queue
			err = u.queue.Push(changeStream.OPType, bytes)
			if err != nil {
				fmt.Printf("failed to push message to queue of token %s, err %v \n", cs.ResumeToken().Lookup("_data").String(), err)
				err = cs.Close(ctx)
				if err != nil {
					fmt.Printf("close collection err %v \n", err)
				}
				fmt.Println("watch retrying...")
				time.Sleep(time.Second * 5)
				break
			}
		}
	}
}
