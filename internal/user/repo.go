package user

import (
	"clomingo/pkg/db"
	"cloud.google.com/go/datastore"
	"context"
	"go.uber.org/zap"
)

type Repo struct {
	logger *zap.Logger
	ds     *datastore.Client
}

func (r Repo) GetByEmail(ctx context.Context, email string) (*User, error) {
	var users []User
	q := datastore.NewQuery("user").Filter("Email =", email).Limit(1)
	keys, err := r.ds.GetAll(ctx, q, &users)
	if err != nil {
		r.logger.Error("error querying users", zap.Error(err), zap.String("email", email))
		return nil, err
	}
	if len(users) < 1 {
		return nil, nil
	}
	userDb := users[0]
	userDb.Id = keys[0].ID
	return &userDb, nil
}

func (r Repo) GetByUserId(ctx context.Context, id int64) (*User, error) {
	var userDb User
	key := datastore.IDKey("user", id, nil)
	err := r.ds.Get(ctx, key, &userDb)
	if err != nil {
		r.logger.Error("error querying users", zap.Error(err), zap.Int64("id", id))
		return nil, err
	}
	return &userDb, nil
}

func (r Repo) Save(ctx context.Context, user *User) (*User, error) {
	put, err := r.ds.Put(ctx, datastore.IncompleteKey("user", nil), user)
	if err != nil {
		r.logger.Error("cannot save user", zap.Error(err), zap.Any("user", user))
		return user, err
	}
	user.Id = put.ID
	return user, err
}

func NewRepo(db *db.Datastore, logger *zap.Logger) *Repo {
	return &Repo{ds: db.DS, logger: logger}
}
