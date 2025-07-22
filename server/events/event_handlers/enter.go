package eventhandlers

import (
	"mlib.com/gofy/server/events"
	"mlib.com/mrun"
)

func Init() {
	if events.Instance.MessageWasCreatedSig != nil {
		mrun.Connect(events.Instance.MessageWasCreatedSig, DeductQuotaWhenMessageCreatedHandle)
		mrun.Connect(events.Instance.MessageWasCreatedSig, UpdateProviderLastUsedAtWhenMessageCreatedHandle)
	}
	if events.Instance.AppWasCreatedSig != nil {
		mrun.Connect(events.Instance.AppWasCreatedSig, create_installed_app_when_app_created)
		mrun.Connect(events.Instance.AppWasCreatedSig, create_site_record_when_app_created)
	}
}
