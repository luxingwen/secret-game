package dao

import (
	"github.com/luxingwen/secret-game/model"
)

func (d *Dao) AddTeam(team *model.Team) error {
	err := d.DB.Table(TableTeam).Create(team).Error
	if err != nil {
		return err
	}
	tuser := model.TeamUserMap{
		TeamId: team.Id,
		UserId: team.LeaderId,
	}

	err = d.DB.Table(TableTeamUser).Create(&tuser).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao) DeleteTeam(id int64) error {
	return d.DB.Table(TableTeam).Where("id = ?", id).Delete(model.Team{}).Error
}

type ResTeamUser struct {
	Id     int64
	Count  int64
	TeamId int64
}

func (d *Dao) List() (res []model.ResTeam, err error) {
	teams := make([]*model.Team, 0)
	err = d.DB.Table(TableTeam).Find(&teams).Error
	if err != nil {
		return
	}

	resTeam := make([]*ResTeamUser, 0)
	err = d.DB.Table(TableTeamUser).Select("id, count(user_id) AS count, team_id").Group("team_id").Find(&resTeam).Error
	if err != nil {
		return
	}

	mTeam := make(map[int64]*ResTeamUser, 0)
	for _, item := range resTeam {
		mTeam[item.TeamId] = item
	}

	for _, item := range teams {

		itemResTeam := model.ResTeam{
			Id:    item.Id,
			Name:  item.Name,
			Score: item.Score,
		}
		if itemTeam, ok := mTeam[item.Id]; ok {
			itemResTeam.Count = itemTeam.Count
		}
		res = append(res, itemResTeam)

	}
	return
}

func (d *Dao) JoinTeam(uid, teamId int64) error {
	teamUser := model.TeamUserMap{
		TeamId: teamId,
		UserId: uid,
	}
	return d.DB.Table(TableTeamUser).Create(&teamUser).Error
}

func (d *Dao) QuitTeam(uid, teamId int64) error {
	return d.DB.Table(TableTeamUser).Where("user_id = ? AND team_id = ?", uid, teamId).Delete(&model.TeamUserMap{}).Error
}
