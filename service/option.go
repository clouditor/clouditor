package service

// Option is a functional option type to configure services.
type Option[T any] func(*T)
