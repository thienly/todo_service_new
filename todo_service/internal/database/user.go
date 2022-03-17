package database

import (
	"context"
	"new_todo_project/internal/domain"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type TodoDb interface {
	Add(ctx context.Context, user *domain.User) (*domain.User, error)
	Verify(ctx context.Context, email, password string) (*domain.User, error)
	Deactive(ctx context.Context, user *domain.User) error
	AddTodo(ctx context.Context,user *domain.User, todo *domain.Todo) error
	MarkTodoDone(ctx context.Context, user *domain.User, todo *domain.Todo) error

}
type userDbImpl struct {
	DB         *mongo.Database
	UserCollection *mongo.Collection
	TodoCollection *mongo.Collection
}

func NewTodoDb(DB *mongo.Database) TodoDb {
	userCollection:= DB.Collection("users")
	todosCollection:= DB.Collection("todos")
	return &userDbImpl{DB: DB,
		UserCollection: userCollection,
		TodoCollection: todosCollection,
	}
}
func (u *userDbImpl) Add(ctx context.Context, user *domain.User) (*domain.User, error) {
	_, err := u.UserCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userDbImpl) Verify(ctx context.Context, email, password string) (*domain.User, error) {
	singleResult := u.UserCollection.FindOne(ctx, bson.D{{"_id",email}, {"password",password}})
	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}
	result:= &domain.User{}
	err := singleResult.Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *userDbImpl) Deactive(ctx context.Context, user *domain.User) error {
	_, err := u.UserCollection.UpdateOne(ctx, bson.M{"_id": user.Email}, bson.D{{"$set", bson.D{{"is_active", false}}}})
	return err
}

func (u *userDbImpl) AddTodo(ctx context.Context,user *domain.User, todo *domain.Todo) error {
	session, err:= u.DB.Client().StartSession()
	if err != nil{
		return err
	}
	defer session.EndSession(ctx)
	callback:= func(sessionCtx mongo.SessionContext) (interface{}, error){
		_, err = u.TodoCollection.InsertOne(ctx, todo)
		if err != nil {
			return nil, err
		}
		user.Todos = append(user.Todos, *todo)
		_, err = u.UserCollection.ReplaceOne(ctx, bson.M{"_id": user.Email}, user)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return err
	}
	return nil
}

func (u *userDbImpl) MarkTodoDone(ctx context.Context, user *domain.User, todo *domain.Todo) error {
	session, err := u.DB.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	callback:= func(sessionCtx mongo.SessionContext) (interface{}, error) {
		// move it out of user collections.
		for i, todo := range user.Todos {
			if todo.Id == todo.Id {
				user.Todos = append(user.Todos[:i], user.Todos[i+1:]...)
			}
		}
		_, err = u.UserCollection.ReplaceOne(ctx, bson.M{"_id": user.Email}, user)
		if err != nil {
			return nil, err
		}
		_, err = u.TodoCollection.ReplaceOne(ctx, bson.M{"_id": todo.Id}, todo)
		if err != nil {
			return nil, err
		}
		return nil,nil
	}
	_, err = session.WithTransaction(ctx, callback)
	if err != nil{
		return err
	}
	return nil
}
