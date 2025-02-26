// Package usercache contains user related CRUD functionality with caching.
package usercache

import (
	"context"
	"net/mail"
	"time"

	"github.com/ardanlabs/service/business/domain/publicuesrbus"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/viccon/sturdyc"
)

// Store manages the set of APIs for user data and caching.
type Store struct {
	log    *logger.Logger
	storer publicuesrbus.Storer
	cache  *sturdyc.Client[publicuesrbus.PublicUser]
}

// NewStore constructs the api for data and caching access.
func NewStore(log *logger.Logger, storer publicuesrbus.Storer, ttl time.Duration) *Store {
	const capacity = 10000
	const numShards = 10
	const evictionPercentage = 10

	return &Store{
		log:    log,
		storer: storer,
		cache:  sturdyc.New[publicuesrbus.PublicUser](capacity, numShards, ttl, evictionPercentage),
	}
}

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, usr publicuesrbus.PublicUser) error {
	if err := s.storer.Create(ctx, usr); err != nil {
		return err
	}

	s.writeCache(usr)

	return nil
}

// Update replaces a user document in the database.
func (s *Store) Update(ctx context.Context, usr publicuesrbus.PublicUser) error {
	if err := s.storer.Update(ctx, usr); err != nil {
		return err
	}

	s.writeCache(usr)

	return nil
}

// Delete removes a user from the database.
func (s *Store) Delete(ctx context.Context, usr publicuesrbus.PublicUser) error {
	if err := s.storer.Delete(ctx, usr); err != nil {
		return err
	}

	s.deleteCache(usr)

	return nil
}

// QueryByEmail gets the specified user from the database by email.
func (s *Store) QueryByEmail(ctx context.Context, email mail.Address) (publicuesrbus.PublicUser, error) {
	cachedUsr, ok := s.readCache(email.Address)
	if ok {
		return cachedUsr, nil
	}

	usr, err := s.storer.QueryByEmail(ctx, email)
	if err != nil {
		return publicuesrbus.PublicUser{}, err
	}

	s.writeCache(usr)

	return usr, nil
}

// readCache performs a safe search in the cache for the specified key.
func (s *Store) readCache(key string) (publicuesrbus.PublicUser, bool) {
	usr, exists := s.cache.Get(key)
	if !exists {
		return publicuesrbus.PublicUser{}, false
	}

	return usr, true
}

// writeCache performs a safe write to the cache for the specified userbus.
func (s *Store) writeCache(bus publicuesrbus.PublicUser) {
	s.cache.Set(bus.ID.String(), bus)
	s.cache.Set(bus.Email.Address, bus)
}

// deleteCache performs a safe removal from the cache for the specified userbus.
func (s *Store) deleteCache(bus publicuesrbus.PublicUser) {
	s.cache.Delete(bus.ID.String())
	s.cache.Delete(bus.Email.Address)
}
