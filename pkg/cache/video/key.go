package videocache

import "strconv"

func getPopularVideoKey() string {
	return "popular_videos"
}

func getVideoKey(videoID int64) string {
	return "video:" + strconv.FormatInt(videoID, 10)
}

func getUserVideoListVersionKey(userID int64) string {
	return "user:" + strconv.FormatInt(userID, 10) + ":video_list:version"
}

func getUserVideoListKey(userID int64, version int64, pageSize int, pageNum int) string {
	return "user:" + strconv.FormatInt(userID, 10) +
		":video_list:" + strconv.FormatInt(version, 10) +
		":" + strconv.Itoa(pageSize) +
		":" + strconv.Itoa(pageNum)
}
