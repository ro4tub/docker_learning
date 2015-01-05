package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"net/rpc"
	"encoding/json"
	common "github.com/ro4tub/docker_learning/common"
)

// Command-line flags.
var (
	bind   = flag.String("bind", ":8080", "Listen address")
	backend = flag.String("backend", "127.0.0.1:9527", "Backend address")
)

type Server struct {
	Backend	*rpc.Client
}

func NewServer() (*Server, error) {
	client, err := rpc.Dial("tcp", *backend)
	if err != nil {
		log.Fatal("rpc.Dial failed:", err)
		return nil, err
	}
	server := &Server{}
	server.Backend = client
	return server, nil
}

// ServeHTTP implements the HTTP user interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 把请求组成rpc转发给后面的logicsvr
	msgid, err := strconv.Atoi(r.FormValue("msg"))
	if err != nil {
		log.Printf("invalid argument msg: %s", r.FormValue("msg"))
		return
	}
	switch msgid {
	case 1:
		log.Printf("recv msgid: %d", msgid)
		// 创建角色：参数name
		req := &common.CreatePlayerMsgReq{r.FormValue("name")}
		ack := common.CreatePlayerMsgAck{}
		err = s.Backend.Call("GameService.CreatePlayer", req, &ack)
		if err != nil {
			log.Fatal("GameService.CreatePlayer error:", err)
			return
		}
		data, err := json.Marshal(ack)
		if err != nil {
			log.Fatal("json.Marshal error:", err)
			return
		}
		w.Write(data)
		break
	case 2:
		log.Printf("recv msgid: %d", msgid)
		// 战斗：参数pid1, 参数pid2
		pid1, err := strconv.ParseInt(r.FormValue("pid1"), 0, 64)
		pid2, err := strconv.ParseInt(r.FormValue("pid2"), 0, 64)
		req := &common.DoBattleMsgReq{pid1, pid2}
		ack := common.DoBattleMsgAck{}
		err = s.Backend.Call("GameService.DoBattle", req, &ack)
		if err != nil {
			log.Fatal("GameService.DoBattle error:", err)
		}
		data, err := json.Marshal(ack)
		if err != nil {
			log.Fatal("json.Marshal error:", err)
			return
		}
		w.Write(data)
		break
	default:
		log.Printf("Invalid msgid: %d", msgid)
	}
}

func main() {
	flag.Parse()
	server, err := NewServer()
	if err != nil {
		log.Fatal("NewServer failed: %v", err)
		return
	}
	http.Handle("/game/", server)
	log.Fatal(http.ListenAndServe(*bind, nil))
}
