// persist/nameresolver.go
package persist

import (
	"strconv"
)

type NameResolver interface {
	GetSystemName(systemID int64) string
	GetSkillName(skillID int32) string
}

// NameResolverImpl is a separate struct that uses DataStore internally.
type NameResolverImpl struct {
	ds *DataStore
}

func NewNameResolver(ds *DataStore) NameResolver {
	return &NameResolverImpl{ds: ds}
}

func (r *NameResolverImpl) GetSystemName(systemID int64) string {
	name, ok := r.ds.SysIdToName[strconv.FormatInt(systemID, 10)]
	if !ok {
		r.ds.logger.Warnf("System ID %d not found", systemID)
		return ""
	}
	return name
}

func (r *NameResolverImpl) GetSkillName(skillID int32) string {
	skill, ok := r.ds.GetSkillTypeByID(strconv.FormatInt(int64(skillID), 10))
	if !ok {
		r.ds.logger.Warnf("Skill ID %d not found", skillID)
		return ""
	}
	return skill.TypeName
}
