// persist/system_resolver.go
package persist

import (
	"strconv"
)

type SystemNameResolver interface {
	GetSystemName(systemID int64) string
}

// SystemNameResolverImpl is a separate struct that uses DataStore internally.
type SystemNameResolverImpl struct {
	ds *DataStore
}

func NewSystemNameResolver(ds *DataStore) SystemNameResolver {
	return &SystemNameResolverImpl{ds: ds}
}

func (r *SystemNameResolverImpl) GetSystemName(systemID int64) string {
	name, ok := r.ds.SysIdToName[strconv.FormatInt(systemID, 10)]
	if !ok {
		r.ds.logger.Warnf("System ID %d not found", systemID)
		return ""
	}
	return name
}
