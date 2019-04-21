package types

type Validator func(key string, value interface{}, exists bool) error

type Normalizor func(key string, value interface{}) (interface{}, error)
