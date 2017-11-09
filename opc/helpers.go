package opc

import (
	"fmt"
	"sort"

	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/go-oracle-terraform/java"
	"github.com/hashicorp/terraform/helper/schema"
)

// Helper function to get a string list from the schema, and alpha-sort it
func getStringList(d *schema.ResourceData, key string) []string {
	if _, ok := d.GetOk(key); !ok {
		return nil
	}
	l := d.Get(key).([]interface{})
	res := make([]string, len(l))
	for i, v := range l {
		res[i] = v.(string)
	}
	sort.Strings(res)
	return res
}

// Helper function to set a string list in the schema, in an alpha-sorted order.
func setStringList(d *schema.ResourceData, key string, value []string) error {
	sort.Strings(value)
	return d.Set(key, value)
}

// Helper function to get an int list from the schema, and numerically sort it
func getIntList(d *schema.ResourceData, key string) []int {
	if _, ok := d.GetOk(key); !ok {
		return nil
	}

	l := d.Get(key).([]interface{})
	res := make([]int, len(l))
	for i, v := range l {
		res[i] = v.(int)
	}
	sort.Ints(res)
	return res
}

func setIntList(d *schema.ResourceData, key string, value []int) error {
	sort.Ints(value)
	return d.Set(key, value)
}

// A user may inadvertently call the database service without passing in the required parameters (because it's optional)
// so we check to make sure that the database client has been initialized
func getDatabaseClient(meta interface{}) (*database.DatabaseClient, error) {
	client := meta.(*OPCClient).databaseClient
	if client == nil {
		return nil, fmt.Errorf("Database Client is not initialized. Make sure to use `database_endpoint` variable or `OPC_DATABASE_ENDPOINT` env variable")
	}
	return client, nil
}

// A user may inadvertently call the java without passing in the required parameters to use that service
// (because it's optional) so we check to make sure that the database client has been initialized
func getJavaClient(meta interface{}) (*java.JavaClient, error) {
	client := meta.(*OPCClient).javaClient
	if client == nil {
		return nil, fmt.Errorf("Java Client is not initialized. Make sure to use `java_endpoint` variable or `OPC_JAVA_ENDPOINT` env variable")
	}
	return client, nil
}
