package constants

import "time"

const (
	FeedCount                   = 20
	DefaultVideoPageSize        = 10
	MaxVideoPageSize            = 50
	PopularVideoCacheExpiration = time.Minute * 10
	PopularVideoCacheCount      = 100
)
