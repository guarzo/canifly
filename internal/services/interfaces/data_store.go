// interfaces/data_store.go
package interfaces

// DataStore defines methods persistence
type DataStore interface {
	FileSystemRepository
	DeletedCharactersRepository
	ConfigRepository
	UserSelectionsRepository
	AccountRepository
	StateRepository
}
