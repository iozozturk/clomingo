package session

import (
	"clomingo/pkg/db"
	"cloud.google.com/go/datastore"
	"context"
	"go.uber.org/zap"
)

type SessionRepo struct {
	logger *zap.Logger
	ds     *datastore.Client
}

func NewSessionRepo(db *db.Datastore, logger *zap.Logger) *SessionRepo {
	return &SessionRepo{ds: db.DS, logger: logger}
}

func (r *SessionRepo) save(ctx context.Context, sess *Session) (*Session, error) {
	put, err := r.ds.Put(ctx, datastore.IncompleteKey("session", nil), sess)
	if err != nil {
		r.logger.Error("cannot save session", zap.Error(err), zap.Any("session", sess))
		return sess, err
	}
	sess.Id = put.ID
	return sess, err
}

func (r *SessionRepo) findKeysByDeviceId(ctx context.Context, deviceId string) ([]int64, error) {
	query := datastore.NewQuery("session").Filter("DeviceId =", deviceId).KeysOnly()
	keys, err := r.ds.GetAll(ctx, query, nil)
	if err != nil {
		r.logger.Error("database query failed", zap.Error(err))
		return nil, err
	}
	var ids []int64
	for _, key := range keys {
		ids = append(ids, key.ID)
	}
	return ids, nil
}

func (r *SessionRepo) deleteMany(ctx context.Context, ids []int64) error {
	keys := idsToKeys(ids)
	err := r.ds.DeleteMulti(ctx, keys)
	if err != nil {
		r.logger.Error("database delete failed", zap.Error(err))
		return err
	}
	return nil
}

func (r *SessionRepo) deleteById(ctx context.Context, id int64) {
	err := r.ds.Delete(ctx, idToKey(id))
	if err != nil {
		r.logger.Error("database delete session failed", zap.Error(err))
	}
	return
}

func (r SessionRepo) GetBySessionToken(ctx context.Context, sessionToken string) (*Session, error) {
	query := datastore.NewQuery("session").Filter("SessionToken =", sessionToken)
	var sessions []Session
	keys, err := r.ds.GetAll(ctx, query, &sessions)
	if err != nil {
		r.logger.Error("database query failed", zap.Error(err))
		return nil, err
	}
	if len(sessions) < 1 {
		return nil, nil
	}
	sessionDb := sessions[0]
	sessionDb.Id = keys[0].ID
	return &sessionDb, nil
}

func idsToKeys(ids []int64) []*datastore.Key {
	var keys []*datastore.Key
	for _, id := range ids {
		keys = append(keys, idToKey(id))
	}
	return keys
}

func idToKey(id int64) *datastore.Key {
	return datastore.IDKey("session", id, nil)
}
