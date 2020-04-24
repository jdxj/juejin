package module

import "encoding/json"

// RespData 是 HTTP 响应的通用形式
type RespData struct {
	State   int             `json:"s"`
	Message string          `json:"m"`
	Data    json.RawMessage `json:"d"`
}

type UserInfo struct {
	ObjectId              string  `json:"objectId"`
	Username              string  `json:"username"`
	FolloweesCount        int     `json:"followeesCount"`
	FollowersCount        int     `json:"followersCount"`
	RankIndex             float64 `json:"rankIndex"`
	TotalViewsCount       int     `json:"totalViewsCount"` // 文章被阅读数量
	Level                 int     `json:"level"`
	JuejinPower           int     `json:"juejinPower"`
	TotalCollectionsCount int     `json:"totalCollectionsCount"` // 获得点赞
}

type FollowInfo struct {
	Followee *struct {
		FolloweeId string `json:"objectId"`
	}
	FollowDatetime string `json:"createdAtString"`
}
