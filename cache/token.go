package cache

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
)

var order_encryption_key string

type CacheClient struct {
	mod *memcache.Client
}

func NewMemCachedStorage() *CacheClient {
	order_encryption_key = os.Getenv("ORDER_ENCRYPTION_KEY")
	return &CacheClient{mod: memcache.New(os.Getenv("MEMCACHE_SERVER"))}
}

func (cc *CacheClient) SaveAccessToken(orderId, accessToken string) (string, error) {
	fmt.Println("KEY:", order_encryption_key)
	if err := cc.mod.Add(&memcache.Item{
		Key:        orderId,
		Value:      []byte(accessToken),
		Expiration: int32(7*60+50) * 60,
	}); err != nil {
		if err == memcache.ErrNotStored {
			return "", fmt.Errorf("token is already stored")
		} else {
			return "", err
		}
	}
	encryptedOrderId, err := encryptOrderID(orderId)
	if err != nil {
		return "", err
	}
	return encryptedOrderId, nil
}

func (cc *CacheClient) GetAccessToken(encryptedOrderId string) (string, string, error) {
	orderId, err := decryptOrderID(encryptedOrderId)
	if err != nil {
		return "", "", err
	}
	item, err := cc.mod.Get(orderId)
	if err != nil {
		return "", "", err
	}
	return orderId, string(item.Value), nil
}

func encryptOrderID(orderID string) (string, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(order_encryption_key)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ciphertext := aead.Seal(nonce, nonce, []byte(orderID), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptOrderID(encryptedOrderID string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedOrderID)
	if err != nil {
		return "", err
	}
	decodedKey, err := base64.StdEncoding.DecodeString(order_encryption_key)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return "", err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce, ciphertext := ciphertext[:12], ciphertext[12:]
	decryptedOrderID, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(decryptedOrderID), nil
}
