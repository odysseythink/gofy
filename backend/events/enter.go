package events

import "github.com/odysseythink/mrun"

type EventGroup struct {
	MessageWasCreatedSig              *mrun.Signal
	AppPublishedWorkflowWasUpdatedSig *mrun.Signal
	AppWasCreatedSig                  *mrun.Signal
	AppModelConfigWasUpdatedSig       *mrun.Signal
	AppDraftWorkflowWasSyncedSig      *mrun.Signal
	TenantWasCreatedSig               *mrun.Signal
	TenantWasUpdatedSig               *mrun.Signal
}

var Instance = new(EventGroup)
