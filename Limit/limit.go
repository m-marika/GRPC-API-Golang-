package Limit

import (
	"calcLab2/grpc_api"
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	grpc_api.UnimplementedCalculateLabParamServer
}

func NewServer() *Server {
	return &Server{}
}

var a = `{  
	"DNA_N":[  
	  {  
		"S": "1"
		"value": {
				"Lower":"50",
				"Upper":"800",
			 	},
		"S": "2"
		"value":{
				"Lower":"800",
				"Upper":"1100",
				},
	  }
	],
	"DNA_Q3":[  
	  {  
		"S": "1"
		"value":{
				"Lower":"40",
				"Upper":"60",
			 	},
		"S": "1"
		"value":{
				"Lower":"8000",
				"Upper":"10000",
				},
	  }
	]
}`

var b = `{
	"applications": [
		{
			"name": "app1",
			"db": {
			   "host": "db2",
			   "user": "root",
			   "pass": "",
			   "dbname": "test"
			}
		},
		{
			"name": "app2",
			"db": {
			   "host": "db2",
			   "user": "root",
			   "pass": "",
			   "dbname": "test"
			}
		}
	 ]
 }`

type Param struct {
	Text []*S `json:"DNA_N"`
}

type S struct {
	S   int32  `json:"S"`
	Res *Value `json:"Value"`
}
type Value struct {
	Lower float32 `json:"Lower"`
	Upper float32 `json:"Upper"`
}

type Range struct {
	Lower        float32
	Upper        float32
	LowerUnbound bool
	UpperUnbound bool
}

///////////////////////////////////////////
type Config struct {
	Applications []Application
}

type Application struct {
	Name string
	Db   Db
}

type Db struct {
	Host   string
	User   string
	Pass   string
	Dbname string
}

///////////////////////////////////////////
func (s *Server) Limit(ctx context.Context, DNA *grpc_api.LimitRequest, S int32) (*grpc_api.LimitReply, error) {

	log.WithFields(log.Fields{
		"package": "main",
		"func":    "Limit",
		"Param":   DNA.GetParam(),
	}).Info("Method limit start, Param =", DNA.GetParam())

	var j Param
	var v Config

	err1 := json.Unmarshal([]byte(a), &v)
	if err1 != nil {
		log.Fatalf("error parsing JSON: %s\n", err1.Error())
	}
	err := json.Unmarshal([]byte(a), &j)
	if err != nil {
		log.Fatalf("error parsing JSON: %s\n", err.Error())
	}
	fmt.Printf("j.Text[0].Res.Lower", j.Text[0].Res.Lower)

	fmt.Printf("j.Text[1].Res.Upper", j.Text[1].Res.Upper)

	//fmt.Printf("%+v\n", j["DNA_Q3"].Text[0].Res1[1].Lower)

	//return &grpc_api.LimitRequest{LimitRange:{Ranges: {Lower:j.Text[0].Res1[0].Lower,Upper:j.Text[0].Res1[0].Upper,LowerUnbound:0,UpperUnbound:0}}}, nil
	//return &grpc_api.LimitReply{Ranges:R}, nil

	return nil, nil
}
