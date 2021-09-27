package dao

import (
	"github.com/luxingwen/secret-game/model"

	"github.com/jinzhu/gorm"

	"math/rand"
	"time"
)

func getRandList(id int, max int) (list []int) {
	for i := id; i <= max; i++ {
		list = append(list, i)
	}

	for i := 1; i < id; i++ {
		list = append(list, i)
	}
	return
}

func (d *Dao) GenTeamTest(teamid int64) (err error) {

	err = d.DB.Table(TableTeamTest).Where("team_id = ?", teamid).Delete(&model.TeamTest{}).Error
	if err != nil {
		return
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))

	index := r.Intn(8) + 1

	ids := getRandList(index, 8)
	ids = append(ids, 9)

	tests := make([]*model.TeamTest, 0)

	var id int64 = 1
	for _, itemId := range ids {
		tests = append(tests, &model.TeamTest{
			TeamId:    teamid,
			SortNo:    id,
			SubjectId: int64(itemId),
		})
		id++
	}

	for _, item := range tests {
		err = d.DB.Table(TableTeamTest).Create(&item).Error
		if err != nil {
			return
		}
	}

	return
}

func (d *Dao) GetSubjectByTestId(id int) (res *model.Subject, err error) {

	test := new(model.TeamTest)
	err = d.DB.Table(TableTeamTest).Where("id = ?", id).First(&test).Error
	if err != nil {
		return
	}
	res = new(model.Subject)
	err = d.DB.Table(TableSubject).Where("id = ?", test.SubjectId).First(&res).Error
	return
}

func (d *Dao) TeamTestList(teamId int64) (res []model.ResTeamTest, err error) {

	resTeamTests := make([]*model.TeamTest, 0)
	err = d.DB.Table(TableTeamTest).Where("team_id = ?", teamId).Order("sort_no ASC").Find(&resTeamTests).Error
	if err != nil {
		return
	}

	subjects := make([]*model.Subject, 0)
	err = d.DB.Table(TableSubject).Find(&subjects).Error
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
			Pic:          subjectItem.Pic,
		}

		if item.AnswerStatus == 1 {
			resItem.AnsInfo = subjectItem.AnsInfo
		}

		if item.HitCount == 1 {
			resItem.Hits = append(resItem.Hits, subjectItem.Hits)
		}

		if item.HitCount >= 2 {
			resItem.Hits = append(resItem.Hits, subjectItem.Hits)
			resItem.Hits = append(resItem.Hits, subjectItem.Hits2)
		}

		res = append(res, resItem)
		if item.AnswerStatus == 0 {
			return
		}
	}
	return
}

//
func (d *Dao) TeatTestUpdateAnswerStatus(id int) error {
	return d.DB.Table(TableTeamTest).Where("id = ?", id).Update(map[string]interface{}{"answer_status": 1}).Error
}

func (d *Dao) GetTeamTestHitsById(id int) (r string, err error) {
	teamTest := new(model.TeamTest)
	err = d.DB.Table(TableTeamTest).Where("id = ?", id).First(&teamTest).Error
	if err != nil {
		return
	}
	subject := new(model.Subject)

	err = d.DB.Table(TableSubject).Where("id = ?", teamTest.SubjectId).First(&subject).Error
	if err != nil {
		return
	}

	err = d.AddTeamTestHitCount(id)
	if err != nil {
		return
	}

	if teamTest.HitCount < 1 {
		r = subject.Hits
		return
	}
	r = subject.Hits2
	return
}

func (d *Dao) AddTeamTestHitCount(id int) error {
	return d.DB.Table(TableTeamTest).Where("id = ?", id).Update("hit_count", gorm.Expr("hit_count + ?", 1)).Error
}

func (d *Dao) AddTeamTestLog(teamTestLog *model.TeamTestLog) error {
	return d.DB.Table(TableTeamTestLog).Create(teamTestLog).Error
}
