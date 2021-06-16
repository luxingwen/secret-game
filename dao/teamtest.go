package dao

import (
	"github.com/luxingwen/secret-game/model"

	"github.com/jinzhu/gorm"
)

func (d *Dao) GenTeamTest(teamid int64) (err error) {
	resSubject := make([]*model.Subject, 0)

	err = d.DB.Table("subject").Find(&resSubject).Error

	if err != nil {
		return
	}

	err = d.DB.Table("team_test").Where("team_id = ?", teamid).Delete(&model.TeamTest{}).Error
	if err != nil {
		return
	}

	mSubject := make(map[int64]*model.Subject, 0)
	for _, item := range resSubject {
		mSubject[item.Id] = item
	}
	tests := make([]model.TeamTest, 0)

	var id int64 = 1
	for _, item := range mSubject {
		tests = append(tests, model.TeamTest{
			TeamId:    teamid,
			SortNo:    id,
			SubjectId: item.Id,
		})
		id++
	}

	err = d.DB.Table("team_test").Create(&tests).Error

	return
}

func (d *Dao) TeamTestList(teamId int64) (res []model.ResTeamTest, err error) {

	var count int64
	err = d.DB.Table("team_test").Where("team_id = ?", teamId).Count(&count).Error
	if err != nil {
		return
	}
	if count == 0 {
		err = d.GenTeamTest(teamId)
		if err != nil {
			return
		}
	}
	resTeamTests := make([]model.TeamTest, 0)
	err = d.DB.Table("team_test").Where("team_id = ?", teamId).Order("sort_no ASC").Find(&resTeamTests).Error
	if err != nil {
		return
	}

	subjects := make([]*model.Subject, 0)
	err = d.DB.Table("subject").Find(&subjects).Error
	if err != nil {
		return
	}

	mSubject := make(map[int64]*model.Subject, 0)
	for _, item := range subjects {
		mSubject[item.Id] = item
	}

	for _, item := range resTeamTests {
		subjectItem := mSubject[item.SubjectId]

		resItem := model.ResTeamTest{
			Id:           item.Id,
			SortNo:       item.SortNo,
			Name:         subjectItem.Name,
			Content:      subjectItem.Content,
			AnswerStatus: item.AnswerStatus,
		}

		if item.HitCount > 0 {
			resItem.Hits = subjectItem.Hits[:item.HitCount]
		}

		res = append(res, resItem)
		if item.AnswerStatus == 0 {
			return
		}
	}
	return

}

//
func (d *Dao) TeatTestUpdateAnswerStatus(id int64) error {
	return d.DB.Table("team_test").Where("id = ?", id).Update("answer_status = ?", 1).Error
}

func (d *Dao) AddTeatTestHitCount(id int64) error {
	return d.DB.Table("team_test").Where("id = ?", id).Update("hit_count", gorm.Expr("hit_count + ?", 1)).Error
}

func (d *Dao) AddTeamTestLog(teamTestLog *model.TeamTestLog) error {
	return d.DB.Table("team_test_log").Create(teamTestLog).Error
}