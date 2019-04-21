package types

type Cache interface {
	GetObject(key string) (map[string]interface{}, error)
	GetArray(key string) ([]map[string]interface{}, error)
	GetPaginatedArray(key string, pagination Pagination) ([]map[string]interface{}, *Pagination, error)
	SetObject(key string, data map[string]interface{}) error
	SetArray(key string, data []map[string]interface{}) error
	SetPaginatedArray(key string, old Pagination, data []map[string]interface{}, new Pagination) error
	Delete(key ...string) error
	DeletePagination(key ...string) error
	GeneratePaginationCacheKeysForId(key string, id int64) ([]string, error)
}
