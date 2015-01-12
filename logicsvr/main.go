package main

import (
	"flag"
	"log"
	"net/rpc"
	"net"
	"github.com/garyburd/redigo/redis"
	"fmt"
	"time"
	"math/rand"
	common "github.com/ro4tub/docker_learning/common"
)

// Command-line flags.
var (
	bind   = flag.String("bind", ":9527", "Listen Address")
)

// gatesvr转发来消息
type GameService struct {
}

func (r *GameService) createPlayerRedis(name string) (int64, error) {
    conn, err := redis.Dial("tcp", ":6379")
    if err != nil {
        return -1, err
    }
    defer conn.Close()
	id, err := redis.Int64(conn.Do("INCR", "playerid"))
	if err != nil {
		return -1, err
	}

	_, err = conn.Do("HMSET", fmt.Sprintf("player:%d", id), "ID", id, "Name", name, "CreateTime", time.Now().Unix(), "FightValue", rand.Intn(10000))
	if err != nil {
		return -1, err
	}
	return id, nil
}

// 创建角色
func (r *GameService) CreatePlayer(req *common.CreatePlayerMsgReq, ack *common.CreatePlayerMsgAck) error {
	if req == nil || req.Name == "" {
		return common.ErrParam
	}
	log.Printf("CreatePlayer: %s\n", req.Name)
	id, err := r.createPlayerRedis(req.Name)
	if err != nil {
		ack.Ret = common.InternalErr
		ack.PlayerId = -1
		return common.ErrNotFoundPlayer
	}
	ack.Ret = common.OK
	ack.PlayerId = id
	return nil
}


func (r *GameService) getFightValueByIdRedis(id int64) (int, error) {
    conn, err := redis.Dial("tcp", ":6379")
    if err != nil {
        return -1, err
    }
    defer conn.Close()
	value, err := redis.Int(conn.Do("HGET", fmt.Sprintf("player:%d", id), "FightValue"))
	if err != nil {
		return -1, err
	}
	return value, nil
}

// 创建角色
func (r *GameService) DoBattle(req *common.DoBattleMsgReq, ack *common.DoBattleMsgAck) error {
	if req == nil || req.PlayerId1 == 0 || req.PlayerId2 == 0 {
		return common.ErrParam
	}
	log.Printf("DoBattle: %d,%d\n", req.PlayerId1, req.PlayerId2)
	
	value1, err := r.getFightValueByIdRedis(req.PlayerId1)
	if err != nil {
		return err
	}
	value2, err := r.getFightValueByIdRedis(req.PlayerId2)
	if err != nil {
		return err
	}
	log.Printf("FightValue: %d, %d\n", value1, value2)
	ack.Ret = common.OK
	if value1 >= value2 {
		ack.Winner = req.PlayerId1
	} else {
		ack.Winner = req.PlayerId2
	}
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
	rand.Seed(time.Now().UnixNano())
	if err := InitRpc(); err != nil {
		panic(err)
	}
	sig := common.InitSignal()
	common.HandleSignal(sig)
}