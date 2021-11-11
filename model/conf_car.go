package model

import (
	"path/filepath"
	"time"

	"github.com/filswan/go-swan-client/config"
)

type ConfCar struct {
	LotusClientApiUrl      string //required
	LotusClientAccessToken string //required
	OutputDir              string //required
	InputDir               string //required
	GocarFileSizeLimit     int64  //required only when creating gocar file(s)
	GenerateMd5            bool   //required
}

func GetConfCar(inputDir string, outputDir *string) *ConfCar {
	confCar := &ConfCar{
		LotusClientApiUrl:      config.GetConfig().Lotus.ClientApiUrl,
		LotusClientAccessToken: config.GetConfig().Lotus.ClientAccessToken,
		OutputDir:              filepath.Join(config.GetConfig().Sender.OutputDir, time.Now().Format("2006-01-02_15:04:05")),
		InputDir:               inputDir,
		GocarFileSizeLimit:     config.GetConfig().Sender.GocarFileSizeLimit,
		GenerateMd5:            config.GetConfig().Sender.GenerateMd5,
	}

	if outputDir != nil && len(*outputDir) != 0 {
		confCar.OutputDir = *outputDir
	}

	return confCar
}
