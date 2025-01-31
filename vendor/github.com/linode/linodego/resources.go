package linodego

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/go-resty/resty/v2"
)

const (
	accountName                  = "account"
	accountSettingsName          = "accountsettings"
	databasesName                = "databases"
	domainRecordsName            = "records"
	domainsName                  = "domains"
	eventsName                   = "events"
	firewallsName                = "firewalls"
	firewallDevicesName          = "firewalldevices"
	firewallRulesName            = "firewallrules"
	imagesName                   = "images"
	instanceConfigsName          = "configs"
	instanceDisksName            = "disks"
	instanceIPsName              = "ips"
	instanceSnapshotsName        = "snapshots"
	instanceStatsName            = "instancestats"
	instanceVolumesName          = "instancevolumes"
	instancesName                = "instances"
	invoiceItemsName             = "invoiceitems"
	invoicesName                 = "invoices"
	ipaddressesName              = "ipaddresses"
	ipv6poolsName                = "ipv6pools"
	ipv6rangesName               = "ipv6ranges"
	kernelsName                  = "kernels"
	lkeClusterAPIEndpointsName   = "lkeclusterapiendpoints"
	lkeClustersName              = "lkeclusters"
	lkeClusterPoolsName          = "lkeclusterpools"
	lkeNodePoolsName             = "lkenodepools"
	lkeVersionsName              = "lkeversions"
	longviewName                 = "longview"
	longviewclientsName          = "longviewclients"
	longviewsubscriptionsName    = "longviewsubscriptions"
	managedName                  = "managed"
	mysqlName                    = "mysql"
	nodebalancerconfigsName      = "nodebalancerconfigs"
	nodebalancernodesName        = "nodebalancernodes"
	nodebalancerStatsName        = "nodebalancerstats"
	nodebalancersName            = "nodebalancers"
	notificationsName            = "notifications"
	oauthClientsName             = "oauthClients"
	objectStorageBucketsName     = "objectstoragebuckets"
	objectStorageBucketCertsName = "objectstoragebucketcerts"
	objectStorageClustersName    = "objectstorageclusters"
	objectStorageKeysName        = "objectstoragekeys"
	paymentsName                 = "payments"
	profileName                  = "profile"
	regionsName                  = "regions"
	sshkeysName                  = "sshkeys"
	stackscriptsName             = "stackscripts"
	tagsName                     = "tags"
	ticketsName                  = "tickets"
	tokensName                   = "tokens"
	typesName                    = "types"
	userGrantsName               = "usergrants"
	usersName                    = "users"
	volumesName                  = "volumes"
	vlansName                    = "vlans"

	accountEndpoint                = "account"
	accountSettingsEndpoint        = "account/settings"
	databasesEndpoint              = "databases"
	domainRecordsEndpoint          = "domains/{{ .ID }}/records"
	domainsEndpoint                = "domains"
	eventsEndpoint                 = "account/events"
	firewallsEndpoint              = "networking/firewalls"
	firewallDevicesEndpoint        = "networking/firewalls/{{ .ID }}/devices"
	firewallRulesEndpoint          = "networking/firewalls/{{ .ID }}/rules"
	imagesEndpoint                 = "images"
	instanceConfigsEndpoint        = "linode/instances/{{ .ID }}/configs"
	instanceDisksEndpoint          = "linode/instances/{{ .ID }}/disks"
	instanceIPsEndpoint            = "linode/instances/{{ .ID }}/ips"
	instanceSnapshotsEndpoint      = "linode/instances/{{ .ID }}/backups"
	instanceStatsEndpoint          = "linode/instances/{{ .ID }}/stats"
	instanceVolumesEndpoint        = "linode/instances/{{ .ID }}/volumes"
	instancesEndpoint              = "linode/instances"
	invoiceItemsEndpoint           = "account/invoices/{{ .ID }}/items"
	invoicesEndpoint               = "account/invoices"
	ipaddressesEndpoint            = "networking/ips"
	ipv6poolsEndpoint              = "networking/ipv6/pools"
	ipv6rangesEndpoint             = "networking/ipv6/ranges"
	kernelsEndpoint                = "linode/kernels"
	lkeClustersEndpoint            = "lke/clusters"
	lkeClusterAPIEndpointsEndpoint = "lke/clusters/{{ .ID }}/api-endpoints"
	lkeClusterPoolsEndpoint        = "lke/clusters/{{ .ID }}/pools"
	lkeNodePoolsEndpoint           = "lke/clusters/{{ .ID }}/pools"
	lkeVersionsEndpoint            = "lke/versions"
	longviewEndpoint               = "longview"
	longviewclientsEndpoint        = "longview/clients"
	longviewsubscriptionsEndpoint  = "longview/subscriptions"
	managedEndpoint                = "managed"
	mysqlEndpoint                  = "databases/mysql/instances"
	// @TODO we can't use these nodebalancer endpoints unless we include these templated fields
	// The API seems inconsistent about including parent IDs in objects, (compare instance configs to nb configs)
	// Parent IDs would be immutable for updates and are ignored in create requests ..
	// Should we include these fields in CreateOpts and UpdateOpts?
	nodebalancerconfigsEndpoint      = "nodebalancers/{{ .ID }}/configs"
	nodebalancernodesEndpoint        = "nodebalancers/{{ .ID }}/configs/{{ .SecondID }}/nodes"
	nodebalancerStatsEndpoint        = "nodebalancers/{{ .ID }}/stats"
	nodebalancersEndpoint            = "nodebalancers"
	notificationsEndpoint            = "account/notifications"
	oauthClientsEndpoint             = "account/oauth-clients"
	objectStorageBucketsEndpoint     = "object-storage/buckets"
	objectStorageBucketCertsEndpoint = "object-storage/buckets/{{ .ID }}/{{ .SecondID }}/ssl"
	objectStorageClustersEndpoint    = "object-storage/clusters"
	objectStorageKeysEndpoint        = "object-storage/keys"
	paymentsEndpoint                 = "account/payments"
	profileEndpoint                  = "profile"
	regionsEndpoint                  = "regions"
	sshkeysEndpoint                  = "profile/sshkeys"
	stackscriptsEndpoint             = "linode/stackscripts"
	tagsEndpoint                     = "tags"
	ticketsEndpoint                  = "support/tickets"
	tokensEndpoint                   = "profile/tokens"
	typesEndpoint                    = "linode/types"
	userGrantsEndpoint               = "account/users/{{ .ID }}/grants"
	usersEndpoint                    = "account/users"
	vlansEndpoint                    = "networking/vlans"
	volumesEndpoint                  = "volumes"
)

// Resource represents a linode API resource
type Resource struct {
	name             string
	endpoint         string
	isTemplate       bool
	endpointTemplate *template.Template
	R                func(ctx context.Context) *resty.Request
	PR               func(ctx context.Context) *resty.Request
}

// NewResource is the factory to create a new Resource struct. If it has a template string the useTemplate bool must be set.
func NewResource(client *Client, name string, endpoint string, useTemplate bool, singleType interface{}, pagedType interface{}) *Resource {
	var tmpl *template.Template

	if useTemplate {
		tmpl = template.Must(template.New(name).Parse(endpoint))
	}

	r := func(ctx context.Context) *resty.Request {
		return client.R(ctx).SetResult(singleType)
	}

	pr := func(ctx context.Context) *resty.Request {
		return client.R(ctx).SetResult(pagedType)
	}

	return &Resource{name, endpoint, useTemplate, tmpl, r, pr}
}

func (r Resource) render(data ...interface{}) (string, error) {
	if data == nil {
		return "", NewError("Cannot template endpoint with <nil> data")
	}
	out := ""
	buf := bytes.NewBufferString(out)

	var substitutions interface{}

	switch len(data) {
	case 1:
		substitutions = struct{ ID interface{} }{data[0]}
	case 2:
		substitutions = struct {
			ID       interface{}
			SecondID interface{}
		}{data[0], data[1]}
	default:
		return "", NewError("Too many arguments to render template (expected 1 or 2)")
	}

	if err := r.endpointTemplate.Execute(buf, substitutions); err != nil {
		return "", NewError(err)
	}
	return buf.String(), nil
}

// endpointWithParams will return the rendered endpoint string for the resource with provided parameters
func (r Resource) endpointWithParams(params ...interface{}) (string, error) {
	if !r.isTemplate {
		return r.endpoint, nil
	}
	return r.render(params...)
}

// Endpoint will return the non-templated endpoint string for resource
func (r Resource) Endpoint() (string, error) {
	if r.isTemplate {
		return "", NewError(fmt.Sprintf("Tried to get endpoint for %s without providing data for template", r.name))
	}
	return r.endpoint, nil
}
