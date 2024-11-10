package grpc

import (
	"context"
	"net"
	"strings"

	"rttmas-backend/pkg/grpc/recognition_report"
	rttmas_binding "rttmas-backend/pkg/binding"
	"rttmas-backend/pkg/utils/logger"

	"google.golang.org/grpc"
)

type RecognitionReportServer struct {
	recognition_report.UnimplementedRecognitionReportServer
}

func (s RecognitionReportServer) Create(ctx context.Context, r *recognition_report.RecognitionReportRequest) (*recognition_report.RecognitionReportResponse, error) {

	logger.Info("Recognition results received via gRPC.\n")

	reportTime := r.ReportTime.GetSeconds()
	reporterUID := r.ReporterUID
	lat := float64(r.Lat)
	lon := float64(r.Lon)
	plateNumbers := strings.Split(r.PlateNumbers, ",")

	for _, plate := range plateNumbers {
		logger.Info(reportTime, " ", reporterUID, " ", lat, " ", lon, " ", plate)
		rttmas_binding.RTTMAS_OnPlateReport(reportTime, lat, lon, plate, reporterUID)
	}

	return &recognition_report.RecognitionReportResponse{
		ResponseStatus: 100,
	}, nil
}



func SetupGrpc() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatal(err)
	}

	serverRegistrar := grpc.NewServer()

	recognitionReportServer := &RecognitionReportServer{}

	recognition_report.RegisterRecognitionReportServer(serverRegistrar, recognitionReportServer)

	err = serverRegistrar.Serve(lis)
	if err != nil {
		logger.Fatal(err)
	}
}
