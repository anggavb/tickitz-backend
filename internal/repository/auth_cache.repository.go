package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
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

func (r *AuthCacheRepository) StoreTokenForgotPassword(
	ctx context.Context,
	token string,
	userID int,
	expiresAt time.Time,
) error {

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return errors.New("token already expired")
	}

	key := fmt.Sprintf("%s:auth:reset-password:%s", r.prefix, token)
	log.Println(key)
	err := r.rdb.Set(ctx, key, userID, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthCacheRepository) IsFogotPasswordKeyActive(ctx context.Context, key string) (bool, error) {
	exists, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (r *AuthCacheRepository) GetValueAndDelete(ctx context.Context, key string) (string, error) {
	val, err := r.rdb.GetDel(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
