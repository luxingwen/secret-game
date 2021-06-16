package dao

import (
	"github.com/luxingwen/secret-game/model"
)

func (d *Dao) AddTeam(team *model.Team) error {
	err := d.DB.Table("team").Create(team).Error
	if err != nil {
		return err
	}
	tuser := model.TeamUserMap{
		TeamId: team.Id,
		UserId: team.LeaderId,
	}

	err = d.DB.Table("team_user").Create(tuser).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao) DeleteTeam(id int64) error {
	return d.DB.Table("team").Where("id = ?", id).Delete(model.Team{}).Error
}

type ResTeamUser struct {
	Id     int64
	Count  int64
	TeamId int64
}

func (d *Dao) List() (res []model.ResTeam, err error) {
	teams := make([]*model.Team, 0)
	err = d.DB.Table("team").Find(&teams).Error
	if err != nil {
		return
	}

	resTeam := make([]*ResTeamUser, 0)
	err = d.DB.Table("team_user").Select("count(user_id) AS count, team_id").Group("team_id").Find(&resTeam).Error
	if err != nil {
		return
	}

	mTeam := make(map[int64]*ResTeamUser, 0)
	for _, item := range resTeam {
		mTeam[item.TeamId] = item
	}

	for _, item := range teams {
		itemTeam := mTeam[item.Id]
		res = append(res, model.ResTeam{
			Id:    item.Id,
			Name:  item.Name,
			Score: item.Score,
			Count: itemTeam.Count,
		})
	}
	return
}

func (d *Dao) JoinTeam(uid, teamId int64) error {
	teamUser := model.TeamUserMap{
		TeamId: teamId,
		UserId: uid,
	}
	return d.DB.Table("team_user").Create(&teamUser).Error
}

func (d *Dao) QuitTeam(uid, teamId int64) error {
	return d.DB.Table("team_user").Where("user_id = ? AND team_id = ?", uid, teamId).Delete(&model.TeamUserMap{}).Error
}
