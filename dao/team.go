package dao

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
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

func (d *Dao) List(searchOp *model.TeamListSearch) (res []model.ResTeam, err error) {
	teams := make([]*model.Team, 0)
	err = d.DB.Table(TableTeam).Find(&teams).Error
	if err != nil {
		fmt.Println("--->err:", err)
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	resTeam := make([]*ResTeamUser, 0)
	err = d.DB.Table(TableTeamUser).Select("id, count(user_id) AS count, team_id").Group("team_id").Find(&resTeam).Error
	if err != nil {
		fmt.Println("--->err2:", err)
		if err.Error() == "record not found" {
			err = nil
		}
		return
	}

	mTeam := make(map[int64]*ResTeamUser, 0)
	for _, item := range resTeam {
		mTeam[item.TeamId] = item
	}

	for _, item := range teams {

		itemResTeam := model.ResTeam{
			Id:       item.Id,
			Name:     item.Name,
			Score:    item.Score,
			LeaderId: item.LeaderId,
		}
		if itemTeam, ok := mTeam[item.Id]; ok {
			itemResTeam.Count = itemTeam.Count

		}
		res = append(res, itemResTeam)
	}

	// todo test and add team img_url into return

	return
}

// 参加队伍前退出队伍
func (d *Dao) BeforeJoinTeamQuitTeam(uid int) error {
	return d.DB.Table(TableTeamUser).Where("user_id = ?", uid).Delete(&model.TeamUserMap{}).Error
}

func (d *Dao) JoinTeam(uid, teamId int) error {
	teamUser := model.TeamUserMap{
		TeamId: int64(teamId),
		UserId: int64(uid),
	}
	return d.DB.Table(TableTeamUser).Create(&teamUser).Error
}

func (d *Dao) QuitTeam(uid, teamId int) (err error) {
	team := new(model.Team)
	err = d.DB.Table(TableTeam).Where("id = ?", teamId).First(&team).Error
	if err != nil {
		return
	}

	if team.LeaderId == int64(uid) {
		err = d.DB.Table(TableTeam).Where("id = ?", teamId).Delete(&model.Team{}).Error
		if err != nil {
			return
		}
	}

	return d.DB.Table(TableTeamUser).Where("user_id = ? AND team_id = ?", uid, teamId).Delete(&model.TeamUserMap{}).Error
}

func (d *Dao) UpdateStatus(teamId int64, status int) error {
	return d.DB.Table(TableTeam).Where("id = ?", teamId).Update(map[string]interface{}{"status": status}).Error
}

func (d *Dao) TeamStartGame(teamId int64) (err error) {
	now := time.Now()
	nowTime := now.Unix() + 3600
	err = d.DB.Table(TableTeam).Where("id = ?", teamId).Update(map[string]interface{}{"status": 1, "end_time": nowTime}).Error
	return
}

func (d *Dao) GetTeamByLeaderId(leaderId int64) (team *model.Team, err error) {
	team = new(model.Team)
	err = d.DB.Table(TableTeam).Where("leader_id = ?", leaderId).First(&team).Error
	return
}

func (d *Dao) GetTeamInfo(uid int) (resTeam *model.ResTeamInfo, err error) {
	teamUser := new(model.TeamUserMap)
	err = d.DB.Table(TableTeamUser).Where("user_id = ?", uid).First(&teamUser).Error
	if err != nil {
		fmt.Println("--->", err)
		return
	}

	teamUserList := make([]*model.TeamUserMap, 0)
	uids := make([]int64, 0)

	err = d.DB.Table(TableTeamUser).Select("user_id").Where("team_id = ?", teamUser.TeamId).Find(&teamUserList).Error
	if err != nil {
		return
	}

	for _, item := range teamUserList {
		uids = append(uids, item.UserId)
	}

	fmt.Println("uids=>", uids)

	teamInfo := new(model.Team)

	err = d.DB.Table(TableTeam).Where("id = ?", teamUser.TeamId).First(&teamInfo).Error
	if err != nil {
		return
	}

	users := make([]model.WxUser, 0)

	err = d.DB.Table(TableWxUser).Where("id IN (?)", uids).Find(&users).Error
	if err != nil {
		return
	}

	resTeam = &model.ResTeamInfo{
		Users: make([]model.ResWxUser, 0),
	}

	for _, item := range users {
		resTeam.Users = append(resTeam.Users, model.ResWxUser{
			Id:        int64(item.ID),
			NickName:  item.NickName,
			AvatarUrl: item.AvatarUrl,
		})
	}

	resTeam.Id = teamInfo.Id
	resTeam.Name = teamInfo.Name
	resTeam.Score = teamInfo.Score
	resTeam.Status = teamInfo.Status
	resTeam.LeaderId = teamInfo.LeaderId
	return
}

func (d *Dao) GetTeamUserMapsByUid(uid int) (r []*model.TeamUserMap, err error) {
	teamUser := new(model.TeamUserMap)

	err = d.DB.Table(TableTeamUser).Where("user_id = ?", uid).First(&teamUser).Error
	if err != nil {
		return
	}

	r = make([]*model.TeamUserMap, 0)
	err = d.DB.Table(TableTeamUser).Where("team_id = ?", teamUser.TeamId).Find(&r).Error
	return
}
