package kvDB

import (
	"github.com/dgraph-io/badger/v3"
	"mirror-api/util/logger"
)

var bDB *badger.DB

func Init() {
	logger.L.Info("Badger key-value DB initializing...")
	db, err := badger.Open(badger.DefaultOptions("badger"))
	if err != nil {
		logger.L.Fatalw(err.Error(), "func", "Init()", "extra", `badger.Open(badger.DefaultOptions("badger"))`)
		return
	}

	bDB = db
	logger.L.Info("Badger key-value DB has been successfully initialized!")
}

func SetValue(key, value string) error {
	txn := bDB.NewTransaction(true)
	if err := txn.Set([]byte(key), []byte(value)); err != nil {
		return err
	}

	return txn.Commit()
}

func DeleteValue(key string) error {
	txn := bDB.NewTransaction(true)
	if err := txn.Delete([]byte(key)); err != nil {
		return err
	}

	return txn.Commit()
}

func GetValue(key string) (string, error) {
	txn := bDB.NewTransaction(false)
	item, err := txn.Get([]byte(key))
	if err != nil {
		return "", err
	}

	val, err := item.ValueCopy(nil)
	if err != nil {
		return "", err
	}

	return string(val), nil
}
