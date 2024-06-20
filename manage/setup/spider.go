package setup

import (
	"encoding/json"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract"
)

type Spider struct {
	*config.Args
	*extract.Config
}

// UnmarshalJSON is the custom unmarshalling for Spider
func (s *Spider) UnmarshalJSON(data []byte) error {
	type Alias Spider
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}
