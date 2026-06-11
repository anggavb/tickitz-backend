package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthCacheRepository struct {
	rdb    *redis.Client
	prefix string
}

func NewAuthCacheRepository(rdb *redis.Client) *AuthCacheRepository {
	return &AuthCacheRepository{
		rdb:    rdb,
		prefix: strings.TrimRight(os.Getenv("RDB_PREFIX"), ":"),
	}
}

func (r *AuthCacheRepository) StoreToken(ctx context.Context, tokenHash string, userID int, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return errors.New("token already expired")
	}

	return r.rdb.Set(ctx, r.tokenKey(userID, tokenHash), userID, ttl).Err()
}

func (r *AuthCacheRepository) IsTokenActive(ctx context.Context, tokenHash string, userID int) (bool, error) {
	exists, err := r.rdb.Exists(ctx, r.tokenKey(userID, tokenHash)).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func (r *AuthCacheRepository) DeleteToken(ctx context.Context, tokenHash string, userID int) error {
	return r.rdb.Del(ctx, r.tokenKey(userID, tokenHash)).Err()
}

func (r *AuthCacheRepository) tokenKey(userID int, tokenHash string) string {
	key := fmt.Sprintf("auth:token~%d:%s", userID, tokenHash)
	if r.prefix == "" {
		return key
	}

	return r.prefix + ":" + key
}
