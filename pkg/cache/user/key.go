package usercache

import "strconv"

func getFollowingKey(userID int64) string {
	return "user:" + strconv.FormatInt(userID, 10) + ":following"
}

func getLikedVideosKey(userID int64) string {
	return "user:" + strconv.FormatInt(userID, 10) + ":liked_videos"
}

func getJwchSessionKey(userID int64) string {
	return "user:" + strconv.FormatInt(userID, 10) + ":jwch"
}
