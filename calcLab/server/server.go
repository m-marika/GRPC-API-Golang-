package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"sync"

	//	"calcLab2/Limit"
	"calcLab2/grpc_api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	grpc_api.UnimplementedCalculateLabParamServer
}

func NewServer() *Server {
	return &Server{}
}

func StartServer(host string, port int, args []string) error {

	log.WithFields(log.Fields{
		"package": "main",
		"func":    "calcLab_server",
		"host":    host,
		"port":    port,
	}).Info("Server start")

	grpcServer := grpc.NewServer()

	c := NewServer()

	reflection.Register(grpcServer)
	grpc_api.RegisterCalculateLabParamServer(grpcServer, c)

	hostS := fmt.Sprintf("%s:%v", host, port)

	lis, err := net.Listen("tcp", hostS)
	if err != nil {
		log.WithFields(log.Fields{
			"package": "main",
			"func":    "calcLab_server",
		}).Error("failed to listen: %v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.WithFields(log.Fields{
			"package": "main",
			"func":    "calcLab_server",
		}).Error("failed to serve: %v", err)
	}

	return nil
}

func (s *Server) CleanStage(ctx context.Context, DNA *grpc_api.DNAcon) (*grpc_api.ResultCleanStage, error) {

	log.WithFields(log.Fields{
		"package": "server",
		"func":    "CleanStage",
	}).Info("CleanStage start, C = %d", DNA.GetC())

	if DNA.GetC() == 0 {

		log.WithFields(log.Fields{
			"package": "server",
			"func":    "CleanStage",
		}).Warning("CleanStage stop because, C = 0")

		return &grpc_api.ResultCleanStage{S: -1}, nil
	}

	if DNA.GetC() > 1000/24 {
		return &grpc_api.ResultCleanStage{S: 2}, nil
	} else {
		return &grpc_api.ResultCleanStage{S: 1}, nil
	}
}

func (s *Server) Volume(ctx context.Context, DNA *grpc_api.DNA) (*grpc_api.ResultVolume, error) {

	log.WithFields(log.Fields{
		"package": "server",
		"func":    "Volume",
		"params":  DNA.GetC(),
	}).Info("Volume start, N = %d", DNA.GetN())

	flagV := false
	if DNA.GetN() == 0 {

		log.WithFields(log.Fields{
			"package": "server",
			"func":    "Volume",
		}).Warning("Volume stop because, N = 0")
		return &grpc_api.ResultVolume{Vsample: -1, Vsalt: -1, Vwater: -1, FlagV: false}, nil
	}
	var Vsample float32
	var Vsalt float32
	var Vwater float32

	Vsample = float32(DNA.GetN()) / float32(DNA.GetC())

	if Vsample > 34 {
		flagV = true
	}

	Vsalt = (35.0 - 0.5*Vsample) / 35.0
	Vwater = 35 - (Vsample + Vsalt)

	return &grpc_api.ResultVolume{Vsample: Vsample, Vsalt: Vsalt, Vwater: Vwater, FlagV: flagV}, nil
}

func (s *Server) VolumeQ(ctx context.Context, DNA *grpc_api.DNAconQ1) (*grpc_api.ResultVolumeQ, error) {

	log.WithFields(log.Fields{
		"package": "server",
		"func":    "VolumeQ",
	}).Info("VolumeQ start, Q1 = %d", DNA.GetQ1())

	if DNA.GetQ1() == 0 {
		log.WithFields(log.Fields{
			"package": "server",
			"func":    "VolumeQ",
		}).Warning("VolumeQ stop because, Q1 = %d", DNA.GetQ1())

		return &grpc_api.ResultVolumeQ{Vq1: -1, Vnete: -1, Nq1: -1}, nil
	}

	var Vq1 float32
	var Vnete float32
	var Nq1 float32

	if 200/DNA.GetQ1() >= 40 {
		Vq1 = 40
	} else {
		Vq1 = 200 / float32(DNA.GetQ1())
	}

	Vnete = 40 - Vq1
	Nq1 = float32(DNA.GetQ1()) - Vq1

	return &grpc_api.ResultVolumeQ{Vq1: Vq1, Vnete: Vnete, Nq1: Nq1}, nil
}

type Params struct {
	Params []*Param
}

type Param struct {
	Param  string `json:"Param"`
	Values []*S
}

type S struct {
	S   string `json:"S"`
	Val Val
}

type Val struct {
	Lower string `json:"Lower"`
	Upper string `json:"Upper"`
}

func (s *Server) Limit(ctx context.Context, DNA *grpc_api.LimitRequest) (*grpc_api.LimitReply, error) {

	log.WithFields(log.Fields{
		"package":     "main",
		"func":        "limit",
		"ParamName":   DNA.GetParamName(),
		"ParamResult": DNA.GetParamResult(),
		"S":           DNA.GetS(),
	}).Info("Method Limit start, Param =", DNA.GetParamName())

	var j Params
	var R grpc_api.LimitRange
	var (
		_lock = &sync.Mutex{}
	)

	count := 0

	S := DNA.GetS()
	if S != 0 {
		S = S - 1
	}

	file, err := os.Open("limits.json")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	_lock.Lock()
	defer _lock.Unlock()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	err1 := json.Unmarshal(data, &j)
	if err1 != nil {
		log.Fatalf("error parsing JSON: ", err.Error())
	}

	for i := 0; i < len(j.Params); i++ {
		if j.Params[i].Param == DNA.GetParamName() {
			count = i
		}
	}

	LowerToFloat, err := strconv.ParseFloat(j.Params[count].Values[S].Val.Lower, 64)
	fmt.Println(LowerToFloat)
	if err != nil {
		log.Fatal(err)
	}

	UpperToFloat, err := strconv.ParseFloat(j.Params[count].Values[S].Val.Upper, 64)
	if err != nil {
		log.Fatal(err)
	}

	R.Lower = float32(LowerToFloat)
	R.Upper = float32(UpperToFloat)

	if DNA.GetParamResult() != 0 {

		if R.Lower > DNA.GetParamResult() {
			R.LowerUnbound = true
		} else {
			R.LowerUnbound = false
		}

		if R.Upper < DNA.GetParamResult() {
			R.UpperUnbound = true
		} else {
			R.UpperUnbound = false
		}
	} else {
		R.LowerUnbound = false
		R.UpperUnbound = false
	}

	return &grpc_api.LimitReply{
		Ranges: []*grpc_api.LimitRange{&R},
	}, nil

}
