package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Team struct {
	Id       int64     `json:"id" gorm:"AUTO_INCREMENT;primary_key;"`
	Name     string    `json:"name" gorm:"column:name;type:varchar(100);unique_index"`
	Score    int64     `json:"score"`
	LeaderId int64     `json:"leader_id"`
	EndTime  int64     `json:"end_time"`
	Status   int       `json:"status"`
	Created  time.Time `json:"created" gorm:"column:created"`
}

func (t Team) TeamName() string {
	return "team"
}

type TeamUserMap struct {
	Id     int64 `json:"id" gorm:"AUTO_INCREMENT;primary_key;"`
	TeamId int64 `json:"team_id"`
	UserId int64 `json:"user_id"`
}

// 试题信息
type Subject struct {
	Id      int64 `json:"id" gorm:"AUTO_INCREMENT;primary_key;"`
	Name    string
	Content string `gorm:"type:TEXT"`
	Hits    string `gorm:"type:TEXT"`
	Hits2   string
	Answer  string
	Pic     string
	AnsInfo string
}

//
type ResTeam struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Score    int64  `json:"score"`
	Count    int64  `json:"count"`
	Status   int    `json:"status"`
	LeaderId int64  `json:"leader_id"`
}

// 测试题信息
type TeamTest struct {
	Id           int64 `json:"id" gorm:"AUTO_INCREMENT;primary_key;"`
	TeamId       int64 `json:"team_id"`       // 队伍id
	SortNo       int64 `json:"sort_no"`       // 排序
	SubjectId    int64 `json:"subject_id"`    // 试题id
	AnswerStatus int   `json:"answer_status"` // 回答状态
	HitCount     int64 `json:"hit_count"`     // 提示次数

}

type ResTeamTest struct {
	Id           int64    `json:"id"`
	SortNo       int64    `json:"sort_no"`
	Name         string   `json:"name"`
	Content      string   `json:"content"`
	Hits         []string `json:"hits"`
	AnswerStatus int      `json:"answer_status"` // 回答状态
	HitCount     int64    `json:"hit_count"`     // 提示次数
	Pic          string   `json:"pic"`
	AnsInfo      string   `json:"ans_info"` // 答案详情
}

// 测试题信息
type TeamTestLog struct {
	Id           int64  `json:"id" gorm:"AUTO_INCREMENT;primary_key;"`
	TeamId       int64  `json:"team_id"` // 队伍id
	TestId       int64  `json:"test_id"`
	UserId       int64  `json:"user_id"`
	AnswerStatus int    `json:"answer_status"` // 回答状态
	Log          string `json:"log"`           // 回答log
}

type WxUser struct {
	gorm.Model
	NickName  string `gorm:"column:nickname"`
	AvatarUrl string `gorm:"column:avatar_url"`
	Gender    int    `gorm:"column:gender"`
	OpenId    string `gorm:"column:open_id;type:varchar(70);unique_index"`
}

type WxCode struct {
	gorm.Model
	Code       string
	SessionKey string
	OpenID     string
}

type ResWxUser struct {
	Id        int64  `json:"id"`
	NickName  string `json:"nickname" gorm:"column:nickname"`
	AvatarUrl string `json:"avatar_url" gorm:"column:avatar_url"`
}

type ResTeamInfo struct {
	ResTeam
	Users []ResWxUser `json:"users"`
}
