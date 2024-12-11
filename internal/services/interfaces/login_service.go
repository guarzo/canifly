package interfaces

type LoginService interface {
	ResolveAccountByState(state string) (string, bool)
	GenerateAnStoreState(value string) (string, error)
	ClearState(state string)
}
