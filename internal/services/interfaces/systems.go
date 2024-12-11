package interfaces

type SystemRepository interface {
	GetSystemName(systemID int64) string
	LoadSystems() error
}
