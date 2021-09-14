package rule

type ProcessType string

const (
	ProcessTypePut        = "put"
	ProcessTypePutArray   = "put_array"
	ProcessTypeCache      = "cache"
	ProcessTypeCacheArray = "cache_array"
	ProcessTypeCompare    = "compare"
)
