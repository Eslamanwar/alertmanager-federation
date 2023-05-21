package alertmanager

import (
	"log"
	"strings"
	"time"
)

var (
	silanceFoundCheck = false
)

func SyncAlertManagerSilences(f AlertManagerFederation, syncPeriod time.Duration) {
	for {

		log.Printf("Syncing alertManager silances for %s env", f.Env)

		parentAlertmanagerSilences, err := GetAlertmanagerSilances(f.ParentAlertManagerURL)
		if err != nil {
			log.Printf("error while getting silances: %s", err)
		}

		if len(parentAlertmanagerSilences) == 0 {
			log.Printf("No silances found for parent alertManager %s ", f.ParentAlertManagerURL)
		}

		for _, pSilence := range parentAlertmanagerSilences {
			if strings.Compare(*pSilence.Status.State, "active") == 0 {
				for _, matcher := range pSilence.Matchers {
					if strings.Compare(*matcher.Name, "cluster") == 0 {
						log.Printf("Found silance on parent AlertManager %s, with match name: %s and value: %s", f.ParentAlertManagerURL, *matcher.Name, *matcher.Value)

						if strings.Contains(*matcher.Value, "01") {
							childAlertManagerURL := f.ChildsAlertManagerURLs["01"]
							ChildSilences, err := GetAlertmanagerSilances(childAlertManagerURL)
							if err != nil {
								log.Printf("error: %s", err)
							}

							for _, cSilence := range ChildSilences {
								if AreSilencesEqual(*pSilence, *cSilence) {
									log.Printf("Skipping silance on the parent alertManager with ID %s, as it is already exist on child alertManager", *pSilence.ID)
									silanceFoundCheck = true
								}
							}

							if !silanceFoundCheck {
								err = PostAlertmanagerSilences(childAlertManagerURL, pSilence.Matchers, pSilence.StartsAt, pSilence.EndsAt, *pSilence.CreatedBy, *pSilence.Comment)
								if err != nil {
									log.Printf("error while pushing silance with error %s", err)
								}
							}
						} else if strings.Contains(*matcher.Value, "02") {
							childAlertManagerURL := f.ChildsAlertManagerURLs["02"]
							ChildSilences, err := GetAlertmanagerSilances(childAlertManagerURL)
							if err != nil {
								log.Printf("error: %s", err)
							}

							for _, cSilence := range ChildSilences {
								if AreSilencesEqual(*pSilence, *cSilence) {
									log.Printf("Skipping silance on the parent alertManager with ID %s, as it is already exist on child alertManager", *pSilence.ID)
									silanceFoundCheck = true
								}
							}

							if !silanceFoundCheck {
								err = PostAlertmanagerSilences(childAlertManagerURL, pSilence.Matchers, pSilence.StartsAt, pSilence.EndsAt, *pSilence.CreatedBy, *pSilence.Comment)
								if err != nil {
									log.Printf("error while pushing silance with error %s", err)
								}
							}
						}
					}
				}
			}
			// revert check to false for the second iteration
			silanceFoundCheck = false
		}

		time.Sleep(syncPeriod)
	}
}
