package resources

import (
	"database/sql"
	"fmt"
)

type Resource struct {
	DB       *sql.DB
	Type     string
	Resource string
}

func AddResource(resource Resource) error {
	query := fmt.Sprintf("INSERT INTO resources VALUES(%q, %q)",
		resource.Type, resource.Resource)
	_, err := resource.DB.Exec(query)
	return err
}

func RemoveResource(resource Resource) error {
	query := fmt.Sprintf("DELETE FROM resources WHERE resource=%q",
		resource.Resource)
	_, err := resource.DB.Exec(query)
	return err
}
