package constants

import "time"

const (
	FeedCount                   = 20
	DefaultVideoPageSize        = 10
	MaxVideoPageSize            = 50
	VideoCacheExpiration        = 24 * time.Hour
	UserVideoCacheExpiration    = 5 * time.Minute
	PopularVideoCacheExpiration = time.Minute * 1
	PopularVideoCacheCount      = 1000
)
