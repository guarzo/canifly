package interfaces

type LoginService interface {
	ResolveAccountAndStatusByState(state string) (string, bool, bool)
	GenerateAndStoreInitialState(value string) (string, error)
	UpdateStateStatusAfterCallBack(state string) error
	ClearState(state string)
}
