package handler

import (
	"context"
	"github.com/y-ttkt/baker/gen/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type BakerHandler struct {
	api.UnimplementedPancakeBakerServiceServer
	report *report
}

type report struct {
	sync.Mutex
	data map[api.Pancake_Menu]int
}

func NewBakerHandler() *BakerHandler {
	return &BakerHandler{
		report: &report{
			data: make(map[api.Pancake_Menu]int),
		},
	}
}

func (h *BakerHandler) Bake(
	ctx context.Context,
	req *api.BakeRequest,
) (*api.BakeResponse, error) {
	//バリデーション
	if req.Menu == api.Pancake_UNKNOWN || req.Menu > api.Pancake_SPICY_CURRY {
		return nil, status.Errorf(codes.InvalidArgument, "パンケーキを選んでください！")
	}
	//パンケーキを焼いて、数を記録します
	now := time.Now()
	h.report.Lock()
	h.report.data[req.Menu] = h.report.data[req.Menu] + 1
	h.report.Unlock()
	//レスポンスを作って返します
	return &api.BakeResponse{
		Pancake: &api.Pancake{
			Menu:           req.Menu,
			ChefName:       "gami", //ワンオペ
			TechnicalScore: rand.Float32(),
			CreateTime: &timestamppb.Timestamp{
				Seconds: now.Unix(),
				Nanos:   int32(now.Nanosecond()),
			},
		},
	}, nil
}

// Report は、焼けたパンケーキの数を報告します。
func (h *BakerHandler) Report(
	ctx context.Context,
	req *api.ReportRequest,
) (*api.ReportResponse, error) {
	counts := make([]*api.Report_BakeCount, 0)
	//レポートを作ります
	h.report.Lock()
	for k, v := range h.report.data {
		counts = append(counts, &api.Report_BakeCount{
			Menu:  k,
			Count: int32(v),
		})
	}
	h.report.Unlock()
	//レスポンスを作って返します
	return &api.ReportResponse{Report: &api.Report{
		BakeCounts: counts,
	},
	}, nil
}
