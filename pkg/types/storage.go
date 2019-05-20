package types

type UpdateListener func(storage Storage, notification Notification)

type Storage interface {
	Close() error
	GetDatabase() Database
	GetCache() Cache
	GetObject(key string, query Query) (map[string]interface{}, error)
	GetArray(key string, query Query) ([]map[string]interface{}, error)
	GetPaginatedArray(key string, query Query, pagination Pagination, recordIdKey string) ([]map[string]interface{}, *Pagination, error)
	Set(queries []Query) error
	Listen(listener UpdateListener, channels ...string) error
	GenerateCacheKey(key string, args ...interface{}) (string, error)
	GetLogger() Logger
}
