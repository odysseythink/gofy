package events

import (
	"log"

	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mrun"
)

func init() {
	var err error
	Instance.TenantWasCreatedSig, err = mrun.NewSignal("tenant-was-created", func(*models.Tenant) {})
	if err != nil {
		log.Printf("[E]create signal failed:%v", err)
	}
	Instance.TenantWasUpdatedSig, err = mrun.NewSignal("tenant-was-updated", func(*models.Tenant) {})
	if err != nil {
		log.Printf("[E]create signal failed:%v", err)
	}
}
