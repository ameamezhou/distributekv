package dbpkg

import (
	"github.com/ameamezhou/distributekv/xlog"
	bolt "go.etcd.io/bbolt"
)

var defaultBucket = []byte("default")

// 定义一下我们要用的变量type 以及相关的调用函数

// DataBase is an open bolt database
type DataBase struct {
	db 	*bolt.DB
}

// NewDataBase Init an instance of a DataBase, that we can work with this
func NewDataBase(dbLocation *string) (*DataBase, error) {
	db, err := bolt.Open(*dbLocation, 0600, nil)
	if err != nil{
		xlog.Fatalf("NewDataBase(%s): %v", *dbLocation, err)
	}
	if err = db.Update(func(tx *bolt.Tx) error {
		_, err1 := tx.CreateBucketIfNotExists(defaultBucket)
		return err1
	}); err != nil {
		db.Close()
		xlog.Fatal("Create default bucket error: %v", err)
	}

	return &DataBase{
		db: db,
	}, err
}

func (db *DataBase)Close (){
	db.db.Close()
}

// SetKey set the key to the requested value or returns error.
func (db *DataBase) SetKey (key string, value []byte) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		return b.Put([]byte(key), value)
	})
}

// GetValue get the value of the requested from a default database
func (db *DataBase) GetValue (key string) ([]byte, error) {
	var result []byte
	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		result = b.Get([]byte(key))
		return nil
	})
	return result, err
}