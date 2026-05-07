package events

import (
	"log"

	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mrun"
)

func init() {
	var err error
	Instance.AppPublishedWorkflowWasUpdatedSig, err = mrun.NewSignal("app-published-workflow-was-updated", func(*models.App, *models.Workflow) {})
	if err != nil {
		log.Printf("[E]create signal failed:%v", err)
	}
	Instance.AppWasCreatedSig, err = mrun.NewSignal("app-was-created", func(*models.App, *models.Account) {})
	if err != nil {
		log.Printf("[E]create signal failed:%v", err)
	}
	Instance.AppModelConfigWasUpdatedSig, err = mrun.NewSignal("app-model-config-was-updated", func(*models.App, *models.AppModelConfig) {})
	if err != nil {
		log.Printf("[E]create signal failed:%v", err)
	}
	Instance.AppDraftWorkflowWasSyncedSig, err = mrun.NewSignal("app-draft-workflow-was-synced", func(*models.App, *models.Workflow) {})
	if err != nil {
		log.Printf("[E]create signal failed:%v", err)
	}
}
