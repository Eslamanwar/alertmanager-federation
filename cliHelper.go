package main

import (
	"alertmanager-federation/pkg/alertmanager"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
)

func syncSilences(c *cli.Context) error {
	var wg1 sync.WaitGroup

	dryRun := c.String("dry-run")
	env := c.String("env")
	syncPeriod := c.Int64("sync-period")

	if strings.Compare(dryRun, "true") == 0 {
		log.Println("Running in Dry Run mode")
	} else {
		var AlertManagerFederationList = []alertmanager.AlertManagerFederation{
			{
				Env:                   "dev",
				ParentAlertManagerURL: "https://alertmanager.dev.com/api/v2/silences",
				ChildsAlertManagerURLs: map[string]string{
					"01": "https://alertmanager-mesh-eks-dev01-usw2.dev.com/api/v2/silences",
					"02": "https://alertmanager-mesh-eks-dev02-usw2.dev.com/api/v2/silences",
				},
			},
			{
				Env:                   "staging",
				ParentAlertManagerURL: "https://alertmanager.staging.com/api/v2/silences",
				ChildsAlertManagerURLs: map[string]string{
					"01": "https://alertmanager-mesh-eks-staging01-usw2.staging.com/api/v2/silences",
					"02": "https://alertmanager-mesh-eks-staging02-usw2.staging.com/api/v2/silences",
				},
			},
			{
				Env:                   "prod",
				ParentAlertManagerURL: "https://alertmanager.prod.com/api/v2/silences",
				ChildsAlertManagerURLs: map[string]string{
					"01": "https://alertmanager-mesh-eks-prod01-usw2.prod.com/api/v2/silences",
					"02": "https://alertmanager-mesh-eks-prod02-usw2.prod.com/api/v2/silences",
				},
			},
		}

		for _, f := range AlertManagerFederationList {
			if strings.Compare(f.Env, env) == 0 {
				go alertmanager.SyncAlertManagerSilences(f, time.Duration(syncPeriod)*time.Second)
				wg1.Add(1)
			}
		}
	}
	wg1.Wait()

	return nil
}
