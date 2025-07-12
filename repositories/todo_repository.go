package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"
	"todolist/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TodoRepository interface {
	CreateTodo(ctx context.Context, todo *models.Todo) (*models.Todo, error)
	GetAll(ctx context.Context, filter bson.M, limit, skip int64) ([]models.Todo, error)
	GetTodoByID(ctx context.Context, id primitive.ObjectID) (*models.Todo, error)
	UpdateTodo(ctx context.Context, id primitive.ObjectID, todo* models.Todo) (*models.Todo, error)
	DeleteTodo(ctx context.Context, id primitive.ObjectID) error
	CountTodos(ctx context.Context, filter bson.M) (int64, error)
}

type todoRepository struct {
	collection *mongo.Collection
}

func NewTodoRepository(db *mongo.Database) TodoRepository { 
	return &todoRepository{
		collection: db.Collection("todos"),
	}
}

func (r *todoRepository) CreateTodo(ctx context.Context, todo *models.Todo) (*models.Todo, error)  {
	todo.ID = primitive.NewObjectID()
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, todo);

	if err != nil {
		return nil, fmt.Errorf("Error the to create task: %w", err)
	}

	return todo, nil;
}

func (r *todoRepository) GetAll(ctx context.Context, filter bson.M, limit, skip int64) ([]models.Todo, error) {
	var todos []models.Todo
	findOptions := options.Find();

	if limit > 0 {
		findOptions.SetLimit(limit)
	}

	if skip >= 0 {
		findOptions.SetSkip(skip)
	}

	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, findOptions);

	defer cursor.Close(ctx);

	if err = cursor.All(ctx, &todos); err != nil {
		return nil, fmt.Errorf("failed to decode tasks: %w", err)
	}

	return todos, nil;
}

func (r *todoRepository) GetTodoByID(ctx context.Context, id primitive.ObjectID) (*models.Todo, error) {
	var todo models.Todo;
	filter := bson.M{"_id": id};

	err := r.collection.FindOne(ctx, filter).Decode(&todo);
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil;
		}
		return nil, fmt.Errorf("Fail the to search the task by ID! error: %w", err)
	}

	return &todo, nil;
}

func (r *todoRepository) UpdateTodo(ctx context.Context, id primitive.ObjectID, todo* models.Todo) (*models.Todo, error) {
	todo.UpdatedAt = time.Now(); 
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "title", Value: todo.Title},
			{Key: "description", Value: todo.Description},
			{Key: "complete", Value: todo.Complete},
			{Key: "updatedAt", Value: todo.UpdatedAt},
		}},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After);

	var updatedTodo models.Todo
	err := r.collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&updatedTodo);

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("Fail the update task: %w", err)
	}
	return &updatedTodo, nil
}

func (r *todoRepository) DeleteTodo(ctx context.Context, id primitive.ObjectID) error  {
	filter := bson.M{"_id": id};

	result, err := r.collection.DeleteOne(ctx, filter);

	if err != nil {
		return fmt.Errorf("Fail the to delete task: %w", err);
	}

	if result.DeletedCount == 0 {
		return errors.New("No tasks found with id to delete")
	}

	return nil;
}

func (r *todoRepository) CountTodos(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter);

	if err != nil {
		return 0, fmt.Errorf("Fail the to count tasks");
	}

	return count, nil;
}