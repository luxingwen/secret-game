package model

type User struct {
	Id   int
	Name string
	Pic  string
}

type Team struct {
	Id       int64  `json:"id" gorm:"AUTO_INCREMENT;primary_key;"`
	Name     string `json:"name"`
	Score    int64  `json:"score"`
	LeaderId int64  `json:"leader_id"`
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
	Content string
	Hits    []string
	Answer  string
}

//
type ResTeam struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Score int64  `json:"score"`
	Count int64  `json:"count"`
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
