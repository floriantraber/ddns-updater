package server

import (
	"fmt"
	"ddns-updater/pkg/models"
	"net"
)

func healthcheckHandler(recordsConfigs []models.RecordConfigType) error {
	for i := range recordsConfigs {
		recordsConfigs[i].M.RLock()
		defer recordsConfigs[i].M.RUnlock()
		if recordsConfigs[i].Status.Code == models.FAIL {
			return fmt.Errorf("%s", recordsConfigs[i].String())
		}
		ips, err := net.LookupIP(recordsConfigs[i].Settings.BuildDomainName())
		if err != nil {
			return fmt.Errorf("%s", err)
		}
		if len(recordsConfigs[i].History.IPs) == 0 {
			return fmt.Errorf("no set IP address found")
		}
		for i := range ips {
			if ips[i].String() != recordsConfigs[i].History.IPs[0] {
				return fmt.Errorf(
					"lookup IP address of %s is not %s",
					recordsConfigs[i].Settings.BuildDomainName(),
					recordsConfigs[i].History.IPs[0],
				)
			}
		}
	}
	return nil
}