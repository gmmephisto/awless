package stats

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/config"
)

const AWLESS_BUCKET = "awless"

type DB struct {
	*bolt.DB
}

func OpenDB(name string) (*DB, error) {
	boltdb, err := bolt.Open(name, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	db := &DB{boltdb}

	if id, err := db.GetStringValue(AWLESS_ID_KEY); err != nil {
		return nil, err
	} else if id == "" {
		if err = db.NewDB(); err != nil {
			return nil, err
		}

	}

	return db, nil
}

func (db *DB) NewDB() error {
	userId, err := cloud.Current.GetUserId()
	if err != nil {
		return err
	}
	newId, err := generateAnonymousId(userId)
	if err != nil {
		return err
	}
	if err = db.SetStringValue(AWLESS_ID_KEY, newId); err != nil {
		return err
	}
	accountId, err := cloud.Current.GetAccountId()
	if err != nil {
		return err
	}
	aId, err := generateAnonymousId(accountId)
	if err != nil {
		return err
	}
	if err = db.SetStringValue(AWLESS_AID_KEY, aId); err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteBucket(name string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(name))
		if b == nil {
			return nil
		}
		e := tx.DeleteBucket([]byte(name))
		return e
	})
}

func (db *DB) GetValue(key string) ([]byte, error) {
	var value []byte
	err := db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket([]byte(AWLESS_BUCKET)); b != nil {
			value = b.Get([]byte(key))
		}
		return nil
	})
	if err != nil {
		return value, err
	}

	return value, nil
}

func (db *DB) GetStringValue(key string) (string, error) {
	str, err := db.GetValue(key)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

func (db *DB) GetTimeValue(key string) (time.Time, error) {
	var t time.Time
	bin, err := db.GetValue(key)
	if err != nil {
		return t, err
	}
	if len(bin) == 0 {
		return t, nil
	}
	err = t.UnmarshalBinary(bin)
	return t, err
}

func (db *DB) GetIntValue(key string) (int, error) {
	str, err := db.GetStringValue(key)
	if err != nil {
		return 0, err
	}
	if str == "" {
		return 0, nil
	}
	return strconv.Atoi(str)
}

func (db *DB) SetValue(key string, value []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(AWLESS_BUCKET))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), value)
	})
}

func (db *DB) SetStringValue(key, value string) error {
	return db.SetValue(key, []byte(value))
}

func (db *DB) SetTimeValue(key string, t time.Time) error {
	bin, err := t.MarshalBinary()
	if err != nil {
		return err
	}
	return db.SetValue(key, bin)
}

func (db *DB) SetIntValue(key string, value int) error {
	return db.SetStringValue(key, strconv.Itoa(value))
}

func generateAnonymousId(seed string) (string, error) {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(config.Salt+seed))), nil
}
