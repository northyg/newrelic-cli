package execution

import "github.com/newrelic/newrelic-client-go/pkg/installevents"

type InstalleventsClient interface {
	CreateInstallEvent(installevents.InstallStatus) (*installevents.InstallEvent, error)
	CreateInstallMetadata(installevents.InputInstallMetadata) (*installevents.InstallMetadata, error)
}
