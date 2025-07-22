package datamanageres

type ManagerGroup struct {
	ApiKeyAuth *ApiKeyAuthManager
}

var ManagerGroupApp = ManagerGroup{
	ApiKeyAuth: &ApiKeyAuthManager{},
}
