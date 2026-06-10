CREATE TABLE IF NOT EXISTS resources(
	type TEXT CHECK (typeof(type)='TEXT'),
	resource TEXT PRIMARY KEY CHECK (typeof(resource)='TEXT')
) ;

