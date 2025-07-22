package services

import "mlib.com/gofy/server/proto/pbapi"

type ExtensionService struct {
}

func (service *ExtensionService) GetCodeBasedExtension(module string) []*pbapi.ModuleExtension {
	return nil
}
