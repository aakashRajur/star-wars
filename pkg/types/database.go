package types

type NotificationListener interface {
	OnNotification(notification Notification)
}

type Database interface {
	Close() error
	GetObject(query Query) (map[string]interface{}, error)
	GetArray(query Query) ([]map[string]interface{}, error)
	GetPaginatedArray(query Query, pagination Pagination, recordIdKey string) ([]map[string]interface{}, *Pagination, error)
	Set(queries []Query) error
	GenerateUpdateQuery(tableName string, args map[string]interface{}, constraints []Constraint) Query
	Notify(listener NotificationListener) error
}
