package global

import (
	"fmt"
	"log"
	"maps"
	"slices"
	"sync"

	"golang.org/x/sync/singleflight"
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
	LANGUAGES = slices.AppendSeq([]string{}, maps.Keys(LANGUAGE_TIMEZONE_MAPPING))
)

func RegisgterModelProvider(mp modelruntimeentities.ModelProvider) error {
	if _, ok := AllModelProviders.Load(mp.ProviderName()); ok {
		log.Printf("[E]provider(%s) already exist\n", mp.ProviderName())
		return fmt.Errorf("provider(%s) already exist", mp.ProviderName())
	}
	AllModelProviders.Store(mp.ProviderName(), mp)
	return nil
}
