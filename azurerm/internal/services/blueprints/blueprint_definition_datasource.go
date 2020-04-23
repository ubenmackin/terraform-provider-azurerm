package blueprints

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"time"
)

func dataSourceArmBlueprintDefinition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmBlueprintDefinitionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"scope": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"subscription",
					"managementGroup",
				}, false),
			},

			// Computed
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"target_scope": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"time_created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func dataSourceArmBlueprintDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Blueprints.BlueprintsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	scope := d.Get("scope").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, scope, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Blueprint Definition %q not found in Scope (%q): %+v", name, scope, err)
		}

		return fmt.Errorf("Read failed for Blueprint Definition (%q) in Sccope (%q): %+v", name, scope, err)
	}

	if resp.ID != nil {
		d.SetId(*resp.ID)
	} else {
		return fmt.Errorf("Failed to retrieve ID for Blueprint %q", name)
	}

	if resp.Description != nil {
		d.Set("description", resp.Description)
	}

	if resp.DisplayName != nil {
		d.Set("display_name", resp.DisplayName)
	}

	d.Set("last_modified", resp.Status.LastModified.String())

	d.Set("time_created", resp.Status.TimeCreated.String())

	if resp.Type != nil {
		d.Set("type", resp.Type)
	}

	if resp.Versions != nil {
		d.Set("versions", resp.Versions)
	}

	return nil
}
