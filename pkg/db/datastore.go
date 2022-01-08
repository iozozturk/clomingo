package db

import (
	"clomingo/pkg/config"
	"cloud.google.com/go/datastore"
	"context"
	"go.uber.org/zap"
)

type Datastore struct {
	DS *datastore.Client
}

func NewDatastore(conf *config.ApplicationConf, logger *zap.Logger) *Datastore {
	logger.Info("creating datastore client", zap.String("projectId", conf.GoogleProjectId))
	client, err := datastore.NewClient(context.Background(), conf.GoogleProjectId)
	if err != nil {
		logger.Fatal("cannot create datastore client", zap.Error(err))
	}
	return &Datastore{client}
}
