package common

import (
	
)

const (
	// game service
)


// 创建角色
type CreatePlayerMsgReq struct {
	Name string // 角色名
}

type CreatePlayerMsgAck struct {
	Ret	int // retcode
	PlayerId int64 // 角色id
}

// 开始战斗
type DoBattleMsgReq struct {
	PlayerId1 int64
	PlayerId2 int64
}

type DoBattleMsgAck struct {
	Ret int // retcode
	Winner int64
}