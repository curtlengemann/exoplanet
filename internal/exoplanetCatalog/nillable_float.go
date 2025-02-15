package exoplanetCatalog

import "encoding/json"

// Custom float type that allows the value to be nil
type nillableFloat struct {
	Value *float64
}

// UnmarshalJSON parses the JSON encoded data. If the data is a number,
// then it is always parsed as a float64. Otherwise, the value is set to nil.
func (nf *nillableFloat) UnmarshalJSON(data []byte) error {
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	if floatValue, ok := value.(float64); ok {
		nf.Value = &floatValue
		return nil
	}
	nf.Value = nil
	return nil
}
