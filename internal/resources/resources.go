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

func GetResources(db *sql.DB) (*[]Resource, error) {
	query := fmt.Sprintf("SELECT type, resource FROM resources")
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for rows.Next() {
		rowsErr := rows.Err()
		if rowsErr != nil {
				return &resources, rowsErr
		}
		var resource Resource
		scanErr := rows.Scan(&resource.Type, &resource.Resource)
		if scanErr != nil {
				return &resources, scanErr
		}
		resources = append(resources, resource)
	}
	return &resources, nil
}

func UpdateResource(resource Resource, newResource string) error {
	query := fmt.Sprintf("UPDATE resources SET resource=%q WHERE " +
		"resource=%q", newResource, resource.Resource)
	_, err := resource.DB.Exec(query)
	return err
}

