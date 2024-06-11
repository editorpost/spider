package fields

type Builder interface {
	Extractor() (ExtractFn, error)
}
