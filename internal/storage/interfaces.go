package storage

type Config interface {
	TypeOfStorage() string
	PathToFileStorage() string
	String() string
}
