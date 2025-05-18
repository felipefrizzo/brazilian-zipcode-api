package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/felipefrizzo/brazilian-zipcode-api/internal/address"
	"github.com/felipefrizzo/brazilian-zipcode-api/internal/correios"
	"github.com/redis/go-redis/v9"
	"github.com/richardwilkes/toolbox/errs"
)

const (
	prefix = "brazilian-zipcode-api:address:"
)

type client struct {
	redis *redis.Client
	ttl   time.Duration

	correiosService correios.Correios
}

// NewClient creates a new address redis client.
func NewClient(addr, username, password string, db int, correiosService correios.Correios, cacheTTL string) (address.AddressRepository, error) {
	ttl, err := strconv.Atoi(cacheTTL)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return &client{
		redis: redis.NewClient(&redis.Options{
			Addr:     addr,
			Username: username,
			Password: password,
			DB:       db,
		}),
		ttl: time.Duration(ttl) * time.Second,

		correiosService: correiosService,
	}, nil
}

func (c *client) getAndCreate(ctx context.Context, zipcode string) (*address.Address, error) {
	addr, err := c.correiosService.GetAddressByZipcode(ctx, zipcode)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	if err := c.Save(ctx, addr, zipcode); err != nil {
		return nil, errs.Wrap(err)
	}

	return addr, nil
}

// Get retrieves an address from redis.
func (c *client) Get(ctx context.Context, zipcode string) (*address.Address, error) {
	key := fmt.Sprintf("%s%s", prefix, zipcode)

	exists, err := c.redis.Exists(ctx, key).Result()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	if exists == 0 {
		slog.Info("address not found in redis, searching in correios api and saving to redis", "zipcode", zipcode)
		addr, err := c.getAndCreate(ctx, zipcode)
		if err != nil {
			return nil, errs.Wrap(err)
		}

		return addr, nil

	}

	ttl, err := c.redis.TTL(ctx, key).Result()
	if err != nil {
		return nil, errs.Wrap(err)
	}

	if ttl <= 0 {
		slog.Info("address expired, searching in correios api and saving to redis", "zipcode", zipcode)
		if err := c.Delete(ctx, zipcode); err != nil {
			return nil, errs.Wrap(err)
		}

		addr, err := c.getAndCreate(ctx, zipcode)
		if err != nil {
			return nil, errs.Wrap(err)
		}

		return addr, nil
	}

	content, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errs.Newf("address not found")
		}
		return nil, errs.Wrap(err)
	}

	var addr address.Address
	if err := json.Unmarshal([]byte(content), &addr); err != nil {
		return nil, errs.Wrap(err)
	}

	return &addr, nil
}

// Save saves an address to redis.
func (c *client) Save(ctx context.Context, address *address.Address, zipcode string) error {
	key := fmt.Sprintf("%s%s", prefix, zipcode)

	content, err := json.Marshal(address)
	if err != nil {
		return err
	}

	return c.redis.Set(ctx, key, content, c.ttl).Err()
}

// Delete deletes an address from redis.
func (c *client) Delete(ctx context.Context, zipcode string) error {
	key := fmt.Sprintf("%s%s", prefix, zipcode)
	return c.redis.Del(ctx, key).Err()
}
