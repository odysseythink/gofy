package services

type ServiceGroup struct {
	Python *PythonService
}

var ServiceGroupApp = ServiceGroup{
	Python: &PythonService{},
}
