package store

import (
	"context"
	"fmt"
	"time"

	"github.com/codenito/example-go-todo-list-api/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ms MongoStore) GetTasks(ctx context.Context) ([]types.Task, error) {
	c, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	// Select tasks collection
	collection := ms.DataBase.Collection("tasks")

	// Get tasks
	cursor, err := collection.Find(c, bson.D{})
	if err != nil {
		return nil, err
	}

	// Convert result to task object
	var tasks []types.Task
	err = cursor.All(c, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (ms MongoStore) CreateTask(ctx context.Context, in types.Task) (*types.Task, error) {
	c, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	// convert task to bson
	bsonTask, err := bson.Marshal(in)
	if err != nil {
		return nil, err
	}

	// Select tasks collection
	collection := ms.DataBase.Collection("tasks")

	res, err := collection.InsertOne(c, bsonTask)
	if err != nil {
		return nil, err
	}

	// populate task with mongo id
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		in.Id = oid
		return &in, nil
	}

	return nil, fmt.Errorf("mongo store: error decoding ObjectID for %s (%s)", in.Name, res.InsertedID)
}

func (ms MongoStore) DeleteTask(ctx context.Context, in types.Task) error {
	c, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	collection := ms.DataBase.Collection("tasks")

	// Select the task with it's unique id
	filter := bson.M{"_id": in.Id}

	res := collection.FindOneAndDelete(c, filter)

	return res.Err()
}
