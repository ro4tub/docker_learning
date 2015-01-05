package main

import (
	"flag"
	"log"
	"net/rpc"
	"net"
	common "github.com/ro4tub/docker_learning/common"
)

// Command-line flags.
var (
	bind   = flag.String("bind", ":9527", "Listen Address")
)

// gatesvr转发来消息
type GameService struct {
}

// 创建角色
func (r *GameService) CreatePlayer(req *common.CreatePlayerMsgReq, ack *common.CreatePlayerMsgAck) error {
	if req == nil || req.Name == "" {
		return common.ErrParam
	}
	log.Printf("CreatePlayer: %s", req.Name)
	ack.Ret = common.OK
	ack.PlayerId = 1234567
	return nil
}

// 创建角色
func (r *GameService) DoBattle(req *common.DoBattleMsgReq, ack *common.DoBattleMsgAck) error {
	if req == nil || req.PlayerId1 == 0 || req.PlayerId2 == 0 {
		return common.ErrParam
	}
	log.Printf("DoBattle: %d,%d", req.PlayerId1, req.PlayerId2)
	ack.Ret = common.OK
	ack.Winner = req.PlayerId1
	return nil
}

func InitRpc() error {
	game := &GameService{}
	rpc.Register(game)
	go rpcListen(*bind)
	return nil
}

func rpcListen(remoteip string) {
	log.Printf("rpc listen: %s\n", remoteip)
	l, err := net.Listen("tcp", remoteip)
	if err != nil {
		log.Printf("net.Listen failed: %s, %v\n", remoteip, err)
		panic(err)
	}
	defer func() {
		if err := l.Close(); err != nil {
			log.Printf("Close failed %v\n", err)
		}
	}()
	rpc.Accept(l)
}



func main() {
	flag.Parse()
	if err := InitRpc(); err != nil {
		panic(err)
	}
	sig := common.InitSignal()
	common.HandleSignal(sig)
}