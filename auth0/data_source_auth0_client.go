package auth0

import (
	"errors"
	"fmt"

	"github.com/auth0/go-auth0"
	"github.com/auth0/go-auth0/management"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func newDataClient() *schema.Resource {
	return &schema.Resource{
		Read:   readDataClient,
		Schema: newClientSchema(),
	}
}

func newClientSchema() map[string]*schema.Schema {
	clientSchema := datasourceSchemaFromResourceSchema(newClient().Schema)
	delete(clientSchema, "client_secret_rotation_trigger")
	delete(clientSchema, "client_secret")
	addOptionalFieldsToSchema(clientSchema, "name", "client_id")
	return clientSchema
}

func readDataClient(d *schema.ResourceData, m interface{}) error {
	clientId := auth0.StringValue(String(d, "client_id"))
	if clientId != "" {
		d.SetId(clientId)
		return readClient(d, m)
	}

	// If not provided ID, perform looking of client by name
	name := auth0.StringValue(String(d, "name"))
	if name == "" {
		return errors.New("no 'client_id' or 'name' was specified")
	}

	api := m.(*management.Management)
	clients, err := api.Client.List(management.IncludeFields("client_id", "name"))
	if err != nil {
		return err
	}
	for _, client := range clients.Clients {
		if auth0.StringValue(client.Name) == name {
			clientId = auth0.StringValue(client.ClientID)
			d.SetId(clientId)
			return readClient(d, m)
		}
	}
	return fmt.Errorf("no client found with 'name' = '%s'", name)
}
