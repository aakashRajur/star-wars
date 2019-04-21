package types

type FatalHandler interface {
	HandleFatal(err error)
}
