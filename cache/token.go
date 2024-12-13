package cache

import (
	"log"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
)

type CacheClient struct {
	mod *memcache.Client
}

func NewMemCachedStorage() *CacheClient {
	return &CacheClient{mod: memcache.New(os.Getenv("MEMCACHE_SERVER"))}
}

func (cc *CacheClient) SaveAccessToken(orderId, accessToken string) error {
	if err := cc.mod.Add(&memcache.Item{Key: orderId, Value: []byte(accessToken), Expiration: int32(7*60+50) * 60}); err != nil {
		if err == memcache.ErrNotStored {
			log.Printf("Token's already stored")
			return nil
		} else {
			return err
		}
	}
	return nil
}

func (cc *CacheClient) GetAccessToken(orderId string) (string, error) {
	item, err := cc.mod.Get(orderId)
	if err != nil {
		return "", err
	}
	return string(item.Value), nil
}
