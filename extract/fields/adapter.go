package fields

import "github.com/editorpost/donq/pkg/vars"

func MapExtractor(m map[string]any) (*Field, error) {

	e := &Field{}
	if err := vars.FromJSON(m, e); err != nil {
		return nil, err
	}

	return e, nil
}

func ExtractorMap(ex *Field) map[string]any {
	return map[string]any{
		"Name":         ex.Name,
		"Cardinality":  ex.Cardinality,
		"Required":     ex.Required,
		"InputFormat":  ex.InputFormat,
		"OutputFormat": ex.OutputFormat,
		"Selector":     ex.Selector,
		"BetweenStart": ex.BetweenStart,
		"BetweenEnd":   ex.BetweenEnd,
		"FinalRegex":   ex.FinalRegex,
		"Multiline":    ex.Multiline,
		"Children":     ex.Children,
	}
}
