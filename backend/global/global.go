package global

import (
	"errors"
	"fmt"
	"log"
	"maps"
	"slices"
	"sync"

	"golang.org/x/sync/singleflight"
	"mlib.com/gofy/server/core/extension"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
)

var (
	// CONFIG              config.Config
	Concurrency_Control = &singleflight.Group{}
	AllModelProviders   sync.Map
	IsInitValidated     = false

	LANGUAGE_TIMEZONE_MAPPING = map[string]string{
		"zh-Hans": "Asia/Shanghai",
		"en-US":   "America/New_York",
		"zh-Hant": "Asia/Taipei",
		"pt-BR":   "America/Sao_Paulo",
		"es-ES":   "Europe/Madrid",
		"fr-FR":   "Europe/Paris",
		"de-DE":   "Europe/Berlin",
		"ja-JP":   "Asia/Tokyo",
		"ko-KR":   "Asia/Seoul",
		"ru-RU":   "Europe/Moscow",
		"it-IT":   "Europe/Rome",
		"uk-UA":   "Europe/Kyiv",
		"vi-VN":   "Asia/Ho_Chi_Minh",
		"ro-RO":   "Europe/Bucharest",
		"pl-PL":   "Europe/Warsaw",
		"hi-IN":   "Asia/Kolkata",
		"tr-TR":   "Europe/Istanbul",
		"fa-IR":   "Asia/Tehran",
		"sl-SI":   "Europe/Ljubljana",
		"th-TH":   "Asia/Bangkok",
	}
	LANGUAGES          = slices.AppendSeq([]string{}, maps.Keys(LANGUAGE_TIMEZONE_MAPPING))
	CodeBasedExtension = &extension.Extension{}
)

func RegisgterModelProvider(mp modelruntimeentities.ModelProvider) error {
	if _, ok := AllModelProviders.Load(mp.ProviderName()); ok {
		log.Printf("[E]provider(%s) already exist\n", mp.ProviderName())
		return fmt.Errorf("provider(%s) already exist", mp.ProviderName())
	}
	AllModelProviders.Store(mp.ProviderName(), mp)
	return nil
}

func RegisgterExtension(extension_class extension.Extensiblor, label map[string]any, form_schema []map[string]any, builtin bool, position int) error {
	if extension_class == nil {
		log.Printf("[E]extension_class is nil\n")
		return errors.New("extension_class is nil")
	}
	if string(extension_class.Module()) == "" || extension_class.Name() == "" {
		log.Printf("[E]extension_class Module(%s) and Name(%s) can't be empty\n", extension_class.Module(), extension_class.Name())
		return fmt.Errorf("extension_class Module(%s) and Name(%s) can't be empty", extension_class.Module(), extension_class.Name())
	}
	if CodeBasedExtension.ModuleExtensions == nil {
		CodeBasedExtension.ModuleExtensions = make(map[extension.ExtensionModuleType]map[string]*extension.ModuleExtension)
	}
	if _, ok := CodeBasedExtension.ModuleExtensions[extension_class.Module()]; !ok {
		CodeBasedExtension.ModuleExtensions[extension_class.Module()] = make(map[string]*extension.ModuleExtension)
	}
	if _, ok := CodeBasedExtension.ModuleExtensions[extension_class.Module()][extension_class.Name()]; ok {
		log.Printf("[E]extension_class Module(%s) and Name(%s) already exist\n", extension_class.Module(), extension_class.Name())
		return fmt.Errorf("extension_class Module(%s) and Name(%s) already exist", extension_class.Module(), extension_class.Name())
	}
	CodeBasedExtension.ModuleExtensions[extension_class.Module()][extension_class.Name()] = &extension.ModuleExtension{
		ExtensionClass: extension_class,
		Name:           extension_class.Name(),
		Label:          label,
		FormSchema:     form_schema,
		Builtin:        builtin,
		Position:       position,
	}
	return nil
}
