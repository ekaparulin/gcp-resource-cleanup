// Package cleanup contains a Pub/Sub Cloud Function.
package cleanup

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// Cleanup consumes a Pub/Sub message.
func Cleanup(ctx context.Context, m PubSubMessage) error {
	log.Println(string(m.Data))

	if string(m.Data) == "instance templates" {
		return cleanupInstanceTemplates(ctx)
	}

	return nil
}

func cleanupInstanceTemplates(ctx context.Context) error {
	project := getGcpProject()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)

	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}

	instanceTemplates := make(map[string]bool)

	for _, tpl := range getInstanceTemplatesInUse(ctx, computeService, project) {
		instanceTemplates[tpl] = true
	}

	offset := getDeleteOlderDays()
	timeStampMonthAgo := time.Now().Add(time.Duration(-1*offset*24) * time.Hour)

	req := computeService.InstanceTemplates.List(project)
	if err := req.Pages(ctx, func(page *compute.InstanceTemplateList) error {
		for _, instanceTemplate := range page.Items {
			if _, inUse := instanceTemplates[instanceTemplate.SelfLink]; inUse {
				continue
			}

			timeStampCreated, _ := time.Parse(time.RFC3339, instanceTemplate.CreationTimestamp)
			isOlder := timeStampCreated.Before(timeStampMonthAgo)
			if !isOlder {
				continue
			}

			fmt.Printf("Going to delete instance template: %v (Created: %v) \n", instanceTemplate.Name, instanceTemplate.CreationTimestamp)
			deleteInstanceTemplate(ctx, computeService, project, instanceTemplate.Name)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	return nil
}

func getGcpProject() string {
	project, present := os.LookupEnv("GCP_PROJECT")
	if present == false {
		log.Fatal("Environment variable GCP_PROJECT is not set!")
	}
	return project
}

func getDeleteOlderDays() int {
	value, present := os.LookupEnv("DELETE_OLDER_DAYS")
	if present == false {
		log.Fatal("Environment variable DELETE_OLDER_DAYS is not set!")
	}
	ret, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal("Environment variable DELETE_OLDER_DAYS is not a number!")
	}

	return ret
}

func getZonesInRegion(ctx context.Context, computeService *compute.Service, project string, region string) []string {
	zones := []string{}

	req := computeService.Zones.List(project)
	if err := req.Pages(ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			if !strings.Contains(zone.Description, region) {
				continue
			}
			zones = append(zones, zone.Description)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	return zones
}

func getInstanceTemplatesInUse(ctx context.Context, computeService *compute.Service, project string) []string {
	temlatesInUse := []string{}
	req := computeService.InstanceGroupManagers.AggregatedList(project)
	if err := req.Pages(ctx, func(page *compute.InstanceGroupManagerAggregatedList) error {
		for _, instanceGroupManagersScopedList := range page.Items {
			for _, mgr := range instanceGroupManagersScopedList.InstanceGroupManagers {
				temlatesInUse = append(temlatesInUse, mgr.InstanceTemplate)
			}

		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	return temlatesInUse
}

func deleteInstanceTemplate(ctx context.Context, computeService *compute.Service, project string, instanceTemplate string) {
	resp, err := computeService.InstanceTemplates.Delete(project, instanceTemplate).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Delete operation status: %#v, server response code %v\n", resp.Status, resp.ServerResponse.HTTPStatusCode)
}
