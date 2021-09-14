package tools

import (
	"container/list"
	"encoding/json"
	"errors"
)

const (
	CACHELEN int = 20
)

var teamChatCache map[int]*list.List = make(map[int]*list.List)

type Chat struct {
	Uid     int    `json:"uid,omitempty"`
	Content string `json:"content,omitempty"`
}

// 聊天信息
func TeamChat(uid int, teamId int, content string) {
	msg := &Chat{Uid: uid, Content: content}

	if teamChatList, ok := teamChatCache[teamId]; ok {
		if teamChatList.Len() >= CACHELEN {
			// 删除第一个
			teamChatList.Remove(teamChatList.Front())
		}
		teamChatList.PushBack(msg)
		return
	}
	teamChatList := list.New()
	teamChatList.PushBack(msg)
	teamChatCache[teamId] = teamChatList
}

// 队伍解散
func TeamDisband(teamId int) {
	if _, ok := teamChatCache[teamId]; ok {
		delete(teamChatCache, teamId)
	}
}

// 获得消息列表
func TeamChatList(teamId int) ([]byte, error) {
	chatList := []Chat{}
	if teamChatList, ok := teamChatCache[teamId]; ok {
		for i := teamChatList.Front(); i != nil; i = i.Next() {
			value := i.Value
			val, _ := value.(Chat)
			chatList = append(chatList, val)
		}

		return json.Marshal(chatList)
	}
	return nil, errors.New("nil chat list")
}
