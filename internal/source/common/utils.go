package common

import (
	"context"
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// Function that matches cluster to tenant using regexRelationsMap.
//
// In case there is no match or regexRelations is nil, it will return nil.
func MatchClusterToTenant(ctx context.Context, nbi *inventory.NetboxInventory, clusterName string, clusterTenantRelations map[string]string) (*objects.Tenant, error) {
	if clusterTenantRelations == nil {
		return nil, nil
	}
	tenantName, err := utils.MatchStringToValue(clusterName, clusterTenantRelations)
	if err != nil {
		return nil, fmt.Errorf("matching cluster to tenant: %s", err)
	}
	if tenantName != "" {
		tenant, ok := nbi.GetTenant(tenantName)
		if !ok {
			tenant, err := nbi.AddTenant(ctx, &objects.Tenant{
				Name: tenantName,
				Slug: utils.Slugify(tenantName),
			})
			if err != nil {
				return nil, fmt.Errorf("add new tenant: %s", err)
			}
			return tenant, nil
		}
		return tenant, nil
	}
	return nil, nil
}

// Function that matches cluster to tenant using regexRelationsMap.
//
// In case there is no match or regexRelations is nil, it will return nil.
func MatchClusterToSite(ctx context.Context, nbi *inventory.NetboxInventory, clusterName string, clusterSiteRelations map[string]string) (*objects.Site, error) {
	if clusterSiteRelations == nil {
		return nil, nil
	}
	siteName, err := utils.MatchStringToValue(clusterName, clusterSiteRelations)
	if err != nil {
		return nil, fmt.Errorf("matching cluster to tenant: %s", err)
	}
	if siteName != "" {
		site, ok := nbi.GetSite(siteName)
		if !ok {
			newSite, err := nbi.AddSite(ctx, &objects.Site{
				Name: siteName,
				Slug: utils.Slugify(siteName),
			})
			if err != nil {
				return nil, fmt.Errorf("add new site: %s", err)
			}
			return newSite, nil
		}
		return site, nil
	}
	return nil, nil
}

// Function that matches vlanName to vlanGroupName using regexRelationsMap.
//
// In case there is no match or regexRelations is nil, it will return default VlanGroup.
func MatchVlanToGroup(ctx context.Context, nbi *inventory.NetboxInventory, vlanName string, vlanGroupRelations map[string]string, vlanGroupSiteRelations map[string]string) (*objects.VlanGroup, error) {
	if vlanGroupRelations == nil {
		vlanGroup, _ := nbi.GetVlanGroup(constants.DefaultVlanGroupName)
		return vlanGroup, nil
	}
	vlanGroupName, err := utils.MatchStringToValue(vlanName, vlanGroupRelations)
	if err != nil {
		return nil, fmt.Errorf("matching vlan to group: %s", err)
	}
	var vlanGroupSite *objects.Site
	if vlanGroupSiteRelations != nil {
		siteName, err := utils.MatchStringToValue(vlanName, vlanGroupSiteRelations)
		if err != nil {
			return nil, fmt.Errorf("matching vlan to site: %s", err)
		}
		if siteName != "" {
			vlanGroupSite, err = nbi.AddSite(ctx, &objects.Site{
				Name: siteName,
				Slug: utils.Slugify(siteName),
			})
			if err != nil {
				return nil, fmt.Errorf("add site: %s", err)
			}
		}
	}
	var vlanGroup *objects.VlanGroup
	if vlanGroupName != "" {
		vlanGroup := &objects.VlanGroup{
			Name:   vlanGroupName,
			Slug:   utils.Slugify(vlanGroupName),
			MinVid: constants.DefaultVID,
			MaxVid: constants.MaxVID,
		}
		if vlanGroupSite != nil {
			vlanGroup.ScopeType = constants.ContentTypeDcimSite
			vlanGroup.ScopeID = vlanGroupSite.ID
		}
		vlanGroup, err := nbi.AddVlanGroup(ctx, vlanGroup)
		if err != nil {
			return nil, fmt.Errorf("add vlan group %+v: %s", vlanGroup, err)
		}
		return vlanGroup, nil
	}
	return vlanGroup, nil
}

// Function that matches vlanName to tenant using vlanTenantRelations regex relations map.
//
// In case there is no match or vlanTenantRelations is nil, it will return nil.
func MatchVlanToTenant(ctx context.Context, nbi *inventory.NetboxInventory, vlanName string, vlanTenantRelations map[string]string) (*objects.Tenant, error) {
	if vlanTenantRelations == nil {
		return nil, nil
	}
	tenantName, err := utils.MatchStringToValue(vlanName, vlanTenantRelations)
	if err != nil {
		return nil, fmt.Errorf("matching vlan to tenant: %s", err)
	}
	if tenantName != "" {
		tenant, ok := nbi.GetTenant(tenantName)
		if !ok {
			tenant, err := nbi.AddTenant(ctx, &objects.Tenant{
				Name: tenantName,
				Slug: utils.Slugify(tenantName),
			})
			if err != nil {
				return nil, fmt.Errorf("add new tenant: %s", err)
			}
			return tenant, nil
		}
		return tenant, nil
	}

	return nil, nil
}

// Function that matches Host from hostName to Site using hostSiteRelations.
//
// In case that there is not match or hostSiteRelations is nil, it will return default site.
func MatchHostToSite(ctx context.Context, nbi *inventory.NetboxInventory, hostName string, hostSiteRelations map[string]string) (*objects.Site, error) {
	if hostSiteRelations == nil {
		return nil, nil
	}
	siteName, err := utils.MatchStringToValue(hostName, hostSiteRelations)
	if err != nil {
		return nil, fmt.Errorf("matching host to site: %s", err)
	}
	if siteName != "" {
		site, ok := nbi.GetSite(siteName)
		if !ok {
			newSite, err := nbi.AddSite(ctx, &objects.Site{
				Name: siteName,
				Slug: utils.Slugify(siteName),
			})
			if err != nil {
				return nil, fmt.Errorf("add new site: %s", err)
			}
			return newSite, nil
		}
		return site, nil
	}
	site, _ := nbi.GetSite(constants.DefaultSite)
	return site, nil
}

// Function that matches Host from hostName to Tenant using hostTenantRelations.
//
// In case that there is not match or hostTenantRelations is nil, it will return nil.
func MatchHostToTenant(ctx context.Context, nbi *inventory.NetboxInventory, hostName string, hostTenantRelations map[string]string) (*objects.Tenant, error) {
	if hostTenantRelations == nil {
		return nil, nil
	}
	tenantName, err := utils.MatchStringToValue(hostName, hostTenantRelations)
	if err != nil {
		return nil, fmt.Errorf("matching host to tenant: %s", err)
	}
	if tenantName != "" {
		site, ok := nbi.GetTenant(tenantName)
		if !ok {
			tenant, err := nbi.AddTenant(ctx, &objects.Tenant{
				Name: tenantName,
				Slug: utils.Slugify(tenantName),
			})
			if err != nil {
				return nil, fmt.Errorf("add new tenant: %s", err)
			}
			return tenant, nil
		}
		return site, nil
	}
	return nil, nil
}

// MatchHostToRole matches Host from hostName to DeviceRole using hostRoleRelations.
//
// In case that there is not match or hostRoleRelations is nil, it will return nil.
func MatchHostToRole(ctx context.Context, nbi *inventory.NetboxInventory, hostName string, hostRoleRelations map[string]string) (*objects.DeviceRole, error) {
	if hostRoleRelations == nil {
		return nil, nil
	}
	roleName, err := utils.MatchStringToValue(hostName, hostRoleRelations)
	if err != nil {
		return nil, fmt.Errorf("matching host to role: %s", err)
	}
	if roleName != "" {
		role, err := nbi.AddDeviceRole(ctx, &objects.DeviceRole{
			Name: roleName,
			Slug: utils.Slugify(roleName),
		})
		if err != nil {
			return nil, fmt.Errorf("add new host role: %s", err)
		}
		return role, nil
	}
	return nil, nil
}

// Function that matches Vm from vmName to Tenant using vmTenantRelations.
//
// In case that there is not match or hostTenantRelations is nil, it will return nil.
func MatchVMToTenant(ctx context.Context, nbi *inventory.NetboxInventory, vmName string, vmTenantRelations map[string]string) (*objects.Tenant, error) {
	if vmTenantRelations == nil {
		return nil, nil
	}
	tenantName, err := utils.MatchStringToValue(vmName, vmTenantRelations)
	if err != nil {
		return nil, fmt.Errorf("matching vm to tenant: %s", err)
	}
	if tenantName != "" {
		site, ok := nbi.GetTenant(tenantName)
		if !ok {
			tenant, err := nbi.AddTenant(ctx, &objects.Tenant{
				Name: tenantName,
				Slug: utils.Slugify(tenantName),
			})
			if err != nil {
				return nil, fmt.Errorf("add new tenant: %s", err)
			}
			return tenant, nil
		}
		return site, nil
	}
	return nil, nil
}

// MatchVMToRole matches VM from vmName to DeviceRole using vmRoleRelations.
//
// In case that there is not match or hostRoleRelations is nil, it will return nil.
func MatchVMToRole(ctx context.Context, nbi *inventory.NetboxInventory, vmName string, vmRoleRelations map[string]string) (*objects.DeviceRole, error) {
	if vmRoleRelations == nil {
		return nil, nil
	}
	roleName, err := utils.MatchStringToValue(vmName, vmRoleRelations)
	if err != nil {
		return nil, fmt.Errorf("matching vm to role: %s", err)
	}
	if roleName != "" {
		role, err := nbi.AddDeviceRole(ctx, &objects.DeviceRole{
			Name: roleName,
			Slug: utils.Slugify(roleName),
		})
		if err != nil {
			return nil, fmt.Errorf("add new vm role: %s", err)
		}
		return role, nil
	}
	return nil, nil
}
