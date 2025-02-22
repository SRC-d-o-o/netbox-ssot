// Common structs and interfaces for all sources
package source

import (
	"context"
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/parser"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/source/dnac"
	"github.com/bl4ko/netbox-ssot/internal/source/fmc"
	"github.com/bl4ko/netbox-ssot/internal/source/fortigate"
	iosxe "github.com/bl4ko/netbox-ssot/internal/source/ios-xe"
	"github.com/bl4ko/netbox-ssot/internal/source/ovirt"
	"github.com/bl4ko/netbox-ssot/internal/source/paloalto"
	"github.com/bl4ko/netbox-ssot/internal/source/proxmox"
	"github.com/bl4ko/netbox-ssot/internal/source/vmware"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// NewSource creates a Source from the given configuration.
func NewSource(
	ctx context.Context,
	config *parser.SourceConfig,
	logger *logger.Logger,
	netboxInventory *inventory.NetboxInventory,
) (common.Source, error) {
	// First we create default tags for the source
	sourceNameTag, err := netboxInventory.AddTag(ctx, &objects.Tag{
		Name:  config.Tag,
		Slug:  utils.Slugify("source-" + config.Name),
		Color: constants.Color(config.TagColor),
		Description: fmt.Sprintf(
			"Automatically created tag by netbox-ssot for source %s",
			config.Name,
		),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating sourceTag: %s", err)
	}
	sourceTypeTag, err := netboxInventory.AddTag(ctx, &objects.Tag{
		Name:  string(config.Type),
		Slug:  utils.Slugify("type-" + string(config.Type)),
		Color: constants.Color(constants.SourceTypeTagColorMap[config.Type]),
		Description: fmt.Sprintf(
			"Automatically created tag by netbox-ssot for source type %s",
			config.Type,
		),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating sourceTypeTag: %s", err)
	}
	commonConfig := common.Config{
		Logger:        logger,
		SourceConfig:  config,
		SourceNameTag: sourceNameTag,
		SourceTypeTag: sourceTypeTag,
		Ctx:           ctx,
		CAFile:        config.CAFile,
	}

	switch config.Type {
	case constants.Ovirt:
		return &ovirt.OVirtSource{Config: commonConfig}, nil
	case constants.Vmware:
		return &vmware.VmwareSource{Config: commonConfig}, nil
	case constants.Dnac:
		return &dnac.DnacSource{Config: commonConfig}, nil
	case constants.Proxmox:
		return &proxmox.ProxmoxSource{Config: commonConfig}, nil
	case constants.PaloAlto:
		return &paloalto.PaloAltoSource{Config: commonConfig}, nil
	case constants.Fortigate:
		return &fortigate.FortigateSource{Config: commonConfig}, nil
	case constants.FMC:
		return &fmc.FMCSource{Config: commonConfig}, nil
	case constants.IOSXE:
		return &iosxe.IOSXESource{Config: commonConfig}, nil
	default:
		return nil, fmt.Errorf("unsupported source type: %s", config.Type)
	}
}
