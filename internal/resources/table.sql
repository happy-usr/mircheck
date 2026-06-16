CREATE TABLE IF NOT EXISTS resources(
	type TEXT CHECK(length(type) > 0),
	resource TEXT PRIMARY KEY CHECK(length(resource) > 0) 
) STRICT;

