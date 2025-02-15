package src

type Root struct {
	Results map[string]QueryResult `json:"results"`
}

// QueryResult represents the result for a specific query (e.g., "C")
type QueryResult struct {
	Status int     `json:"status"`
	Frames []Frame `json:"frames"`
}

// Frame represents a single frame in the response
type Frame struct {
	Schema Schema `json:"schema"`
	Data   Data   `json:"data"`
}

// Schema represents the schema of the response
type Schema struct {
	RefID  string  `json:"refId"`
	Meta   Meta    `json:"meta"`
	Fields []Field `json:"fields"`
}

// Meta contains metadata for the query
type Meta struct {
	Type                string     `json:"type"`
	TypeVersion         []int      `json:"typeVersion"`
	Custom              CustomMeta `json:"custom"`
	ExecutedQueryString string     `json:"executedQueryString"`
}

// CustomMeta represents additional custom metadata
type CustomMeta struct {
	ResultType string `json:"resultType"`
}

// Field represents a single field in the schema
type Field struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	TypeInfo TypeInfo          `json:"typeInfo"`
	Labels   map[string]string `json:"labels,omitempty"`
	Config   FieldConfig       `json:"config"`
}

// TypeInfo provides type-specific information
type TypeInfo struct {
	Frame string `json:"frame"`
}

// FieldConfig provides configuration settings for a field
type FieldConfig struct {
	Interval          int    `json:"interval,omitempty"`
	DisplayNameFromDS string `json:"displayNameFromDS,omitempty"`
}

// Data represents the actual data in the response
type Data struct {
	Values [][]interface{} `json:"values"` // Can contain int64 timestamps and float64 values
}
