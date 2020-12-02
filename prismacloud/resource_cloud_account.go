package prismacloud

import (
	"encoding/json"
	"log"
	"time"

	pc "github.com/paloaltonetworks/prisma-cloud-go"
	"github.com/paloaltonetworks/prisma-cloud-go/cloud/account"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCloudAccount() *schema.Resource {
	return &schema.Resource{
		Create: createCloudAccount,
		Read:   readCloudAccount,
		Update: updateCloudAccount,
		Delete: deleteCloudAccountAction,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			// AWS type.
			account.TypeAws: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AWS account type",
				MaxItems:    1,
				ConflictsWith: []string{
					account.TypeAzure,
					account.TypeGcp,
					account.TypeAlibaba,
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "AWS account ID",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not the account is enabled",
							Default:     true,
						},
						"external_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "AWS account external ID",
							Sensitive:   true,
						},
						"group_ids": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of account IDs to which you are assigning this account",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name to be used for the account on the Prisma Cloud platform (must be unique)",
						},
						"role_arn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique identifier for an AWS resource (ARN)",
						},
						"account_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "account",
							Description: "Account type - organization or account",
						},
						"protection_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "MONITOR",
							Description: "Monitor or Monitor and Protect",
						},
					},
				},
			},

			// Azure type.
			account.TypeAzure: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Azure account type",
				MaxItems:    1,
				ConflictsWith: []string{
					account.TypeAws,
					account.TypeGcp,
					account.TypeAlibaba,
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Azure account ID",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not the account is enabled",
							Default:     true,
						},
						"group_ids": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of account IDs to which you are assigning this account",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name to be used for the account on the Prisma Cloud platform (must be unique)",
						},
						"client_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Application ID registered with Active Directory",
						},
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Application ID key",
							Sensitive:   true,
						},
						"monitor_flow_logs": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Automatically ingest flow logs",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Active Directory ID associated with Azure",
						},
						"service_principal_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique ID of the service principal object associated with the Prisma Cloud application that you create",
						},
						"account_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "account",
							Description: "Account type - organization or account",
						},
						"protection_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "MONITOR",
							Description: "Monitor or Monitor and Protect",
							ForceNew:    true,
						},
					},
				},
			},

			// GCP type.
			account.TypeGcp: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "GCP account type",
				MaxItems:    1,
				ConflictsWith: []string{
					account.TypeAws,
					account.TypeAzure,
					account.TypeAlibaba,
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "GCP project ID",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not the account is enabled",
							Default:     true,
						},
						"group_ids": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of account IDs to which you are assigning this account",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name to be used for the account on the Prisma Cloud platform (must be unique)",
						},
						"compression_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enable flow log compression",
						},
						"dataflow_enabled_project": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "GCP project for flow log compression",
						},
						"flow_log_storage_bucket": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "GCP flow logs storage bucket",
						},
						// Use a json string until this feature is added:
						// https://github.com/hashicorp/terraform-plugin-sdk/issues/248
						"credentials_json": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Content of the JSON credentials file",
							Sensitive:        true,
							DiffSuppressFunc: gcpCredentialsMatch,
						},
						"account_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "account",
							Description: "Account type - organization or account",
						},
						"protection_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "MONITOR",
							Description: "Monitor or Monitor and Protect",
						},
					},
				},
			},

			// Alibaba type.
			account.TypeAlibaba: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Alibaba account type",
				MaxItems:    1,
				ConflictsWith: []string{
					account.TypeAws,
					account.TypeAzure,
					account.TypeGcp,
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Alibaba account ID",
						},
						"group_ids": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of account IDs to which you are assigning this account",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name to be used for the account on the Prisma Cloud platform (must be unique)",
						},
						"ram_arn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique identifier for an Alibaba RAM role resource",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not the account is enabled",
							Default:     true,
						},
					},
				},
			},

			"disable_on_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not the account will be disabled on terraform destroy rather than deleted. True means the account will be disabled",
				Default:     false,
			},

			"update_on_create": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true and the account already exists, the account will be updated rather than failing on the initial creation of this resource",
				Default:     false,
			},
		},
	}
}

func gcpCredentialsMatch(k, old, new string, d *schema.ResourceData) bool {
	var (
		err       error
		prev, cur account.GcpCredentials
	)

	if err = json.Unmarshal([]byte(old), &prev); err != nil {
		return false
	}

	if err = json.Unmarshal([]byte(new), &cur); err != nil {
		return false
	}

	return (prev.Type == cur.Type &&
		prev.ProjectId == cur.ProjectId &&
		prev.PrivateKeyId == cur.PrivateKeyId &&
		prev.PrivateKey == cur.PrivateKey &&
		prev.ClientEmail == cur.ClientEmail &&
		prev.ClientId == cur.ClientId &&
		prev.AuthUri == cur.AuthUri &&
		prev.TokenUri == cur.TokenUri &&
		prev.ProviderCertUrl == cur.ProviderCertUrl &&
		prev.ClientCertUrl == cur.ClientCertUrl)
}

func parseCloudAccount(d *schema.ResourceData) (string, string, interface{}) {
	if x := ResourceDataInterfaceMap(d, account.TypeAws); len(x) != 0 {
		return account.TypeAws, x["name"].(string), account.Aws{
			AccountId:      x["account_id"].(string),
			Enabled:        x["enabled"].(bool),
			ExternalId:     x["external_id"].(string),
			GroupIds:       ListToStringSlice(x["group_ids"].([]interface{})),
			Name:           x["name"].(string),
			RoleArn:        x["role_arn"].(string),
			ProtectionMode: x["protection_mode"].(string),
			AccountType:    x["account_type"].(string),
		}
	} else if x := ResourceDataInterfaceMap(d, account.TypeAzure); len(x) != 0 {
		return account.TypeAzure, x["name"].(string), account.Azure{
			Account: account.CloudAccount{
				AccountId:      x["account_id"].(string),
				Enabled:        x["enabled"].(bool),
				GroupIds:       ListToStringSlice(x["group_ids"].([]interface{})),
				Name:           x["name"].(string),
				ProtectionMode: x["protection_mode"].(string),
				AccountType:    x["account_type"].(string),
			},
			ClientId:           x["client_id"].(string),
			Key:                x["key"].(string),
			MonitorFlowLogs:    x["monitor_flow_logs"].(bool),
			TenantId:           x["tenant_id"].(string),
			ServicePrincipalId: x["service_principal_id"].(string),
		}
	} else if x := ResourceDataInterfaceMap(d, account.TypeGcp); len(x) != 0 {
		var creds account.GcpCredentials
		_ = json.Unmarshal([]byte(x["credentials_json"].(string)), &creds)

		return account.TypeGcp, x["name"].(string), account.Gcp{
			Account: account.CloudAccount{
				AccountId:      x["account_id"].(string),
				Enabled:        x["enabled"].(bool),
				GroupIds:       ListToStringSlice(x["group_ids"].([]interface{})),
				Name:           x["name"].(string),
				ProtectionMode: x["protection_mode"].(string),
				AccountType:    x["account_type"].(string),
			},
			CompressionEnabled:     x["compression_enabled"].(bool),
			DataflowEnabledProject: x["dataflow_enabled_project"].(string),
			FlowLogStorageBucket:   x["flow_log_storage_bucket"].(string),
			Credentials:            creds,
		}
	} else if x := ResourceDataInterfaceMap(d, account.TypeAlibaba); len(x) != 0 {
		return account.TypeAlibaba, x["name"].(string), account.Alibaba{
			AccountId: x["account_id"].(string),
			GroupIds:  ListToStringSlice(x["group_ids"].([]interface{})),
			Name:      x["name"].(string),
			RamArn:    x["ram_arn"].(string),
			Enabled:   x["enabled"].(bool),
		}
	}

	return "", "", nil
}

func saveCloudAccount(d *schema.ResourceData, dest string, obj interface{}) {
	var val map[string]interface{}

	switch v := obj.(type) {
	case account.Aws:
		val = map[string]interface{}{
			"account_id":      v.AccountId,
			"enabled":         v.Enabled,
			"external_id":     v.ExternalId,
			"group_ids":       v.GroupIds,
			"name":            v.Name,
			"role_arn":        v.RoleArn,
			"protection_mode": v.ProtectionMode,
			"account_type":    v.AccountType,
		}
	case account.Azure:
		val = map[string]interface{}{
			"account_id":           v.Account.AccountId,
			"enabled":              v.Account.Enabled,
			"group_ids":            v.Account.GroupIds,
			"name":                 v.Account.Name,
			"client_id":            v.ClientId,
			"key":                  v.Key,
			"monitor_flow_logs":    v.MonitorFlowLogs,
			"tenant_id":            v.TenantId,
			"service_principal_id": v.ServicePrincipalId,
			"protection_mode":      v.Account.ProtectionMode,
			"account_type":         v.Account.AccountType,
		}
	case account.Gcp:
		b, _ := json.Marshal(v.Credentials)
		val = map[string]interface{}{
			"account_id":               v.Account.AccountId,
			"enabled":                  v.Account.Enabled,
			"group_ids":                v.Account.GroupIds,
			"name":                     v.Account.Name,
			"compression_enabled":      v.CompressionEnabled,
			"dataflow_enabled_project": v.DataflowEnabledProject,
			"flow_log_storage_bucket":  v.FlowLogStorageBucket,
			"credentials_json":         string(b),
			"protection_mode":          v.Account.ProtectionMode,
			"account_type":             v.Account.AccountType,
		}
	case account.Alibaba:
		val = map[string]interface{}{
			"account_id": v.AccountId,
			"group_ids":  v.GroupIds,
			"name":       v.Name,
			"ram_arn":    v.RamArn,
			"enabled":    v.Enabled,
		}
	}

	for _, key := range []string{account.TypeAws, account.TypeAzure, account.TypeGcp, account.TypeAlibaba} {
		if key != dest {
			d.Set(key, nil)
			continue
		}

		if err := d.Set(key, []interface{}{val}); err != nil {
			log.Printf("[WARN] Error setting %q field for %q: %s", key, d.Id(), err)
		}
	}
}

func createCloudAccount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pc.Client)
	cloudType, name, obj := parseCloudAccount(d)
	updateIfExists := d.Get("update_on_create").(bool)
	successfulUpdate := false

	if err := account.Create(client, obj); err != nil {
		if err == pc.DuplicateCloudAccountError && updateIfExists {
			log.Print("[WARN] Duplicate cloud account detected. Attempting to update")
			// We don't call updateCloudAccount here because it makes a call
			// to readCloudAccount which we also do below. Don't want to duplicate
			// that call. Also, we couldn't readCloudAccount yet since d.SetId hasn't
			// been called yet.
			if err := account.Update(client, obj); err != nil {
				return err
			}
			successfulUpdate = true
		} else {
			return err
		}
	}

	id, err := account.Identify(client, cloudType, name)
	if err != nil {
		if err == pc.ObjectNotFoundError && updateIfExists && successfulUpdate {
			// I've observed when changing the name of the account, for example,
			// that the name is correctly updated in the results of: https://api.docs.prismacloud.io/reference#get-cloud-accounts
			// but is not correctly updated in the results of: https://api.docs.prismacloud.io/reference#get-cloud-account-names.
			// This causes Identify to fail, since its comparison includes name, when
			// update_on_create is true and we update an existing cloud account.
			// In this case, build the value of id rather than getting it from Identify. Subsequent
			// updates do not fail if multiple state files constantly change the one account's name.
			log.Printf("[WARN] Failed to identify cloud account when updating existing account. type: %s, name: %s. Constructing id", cloudType, name)
			parsedID, parsedIDErr := account.GetID(obj)
			if parsedIDErr != nil {
				return parsedIDErr
			}
			id = parsedID
			log.Printf("[DEBUG] Account id is %s. type: %s, name: %s", id, cloudType, name)
		} else {
			return err
		}
	}

	d.SetId(TwoStringsToId(cloudType, id))
	return readCloudAccount(d, meta)
}

func readCloudAccount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pc.Client)
	cloudType, id := IdToTwoStrings(d.Id())

	obj, err := account.Get(client, cloudType, id)
	if err != nil {
		if err == pc.ObjectNotFoundError {
			d.SetId("")
			return nil
		}
		return err
	}

	saveCloudAccount(d, cloudType, obj)

	return nil
}

func updateCloudAccount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pc.Client)

	_, _, obj := parseCloudAccount(d)

	if err := account.Update(client, obj); err != nil {
		return err
	}

	return readCloudAccount(d, meta)
}

func deleteCloudAccountAction(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pc.Client)
	cloudType, id := IdToTwoStrings(d.Id())
	shouldDisable := d.Get("disable_on_destroy").(bool)

	var err error

	if shouldDisable {
		err = account.Disable(client, cloudType, id)
	} else {
		err = account.Delete(client, cloudType, id)
	}

	if err != nil {
		if err != pc.ObjectNotFoundError {
			return err
		}
	}

	d.SetId("")
	return nil
}
