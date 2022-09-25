package orm

import "database/sql"

type DBOption func(*DB)

// DB 是sql.DB 的装饰器
type DB struct {
	db *sql.DB
	r *registry
}

// 如果用户指定了 registry，就用用户指定的，否则用默认的

// db := Open()

// r1 := NewRegistry()
// db1 := Open(r1)
// db2 := Open(r1)

func Open(driver string, dsn string, opts...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)

	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

// OpenDB
// 我可以利用 OpenDB 来传入一个 mock 的DB
// sqlmock.Open 的 DB
func OpenDB(db *sql.DB, opts...DBOption) (*DB, error) {
	res := &DB{
		r: &registry{},
		db: db,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}


// func MustNewDB(opts...DBOption) *DB{
// 	res, err := Open(opts...)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return res
// }

func DBWithRegistry(r *registry) DBOption {
	return func(db *DB) {
		db.r = r
	}
}

//
// func (db *DB) NewSelector[T any]() *Selector[T] {
// 	return &Selector[T]{
// 		db: db,
// 	}
// }