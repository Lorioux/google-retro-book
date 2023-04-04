package collectors

import (
	crm "google.golang.org/api/cloudresourcemanager/v1"
	cam "google.golang.org/api/cloudasset/v1"
)

type (
	// Project struct {
	// 	// createTime string
	// 	// lifecycleState string
	// 	// name string
	// 	// parent string
	// 	// projectId string
	// 	// projectNumber int
	// }
)

var (
	CrmService crm.Service
	CamService cam.Service
)
