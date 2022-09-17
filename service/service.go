package service

import (
	"context"

	"github.com/joshsoftware/sparkode-core/isolateutil"
	"go.uber.org/zap"
)

type Service interface {
	Run(context.Context, string, string, string) (string, error)
}

type runService struct {
	logger *zap.SugaredLogger
}

func (s *runService) Run(ctx context.Context, code string, language string, input string) (string, error) {
	var output string = ""
	file, err := createFile(code, language)
	if err != nil {
		return output, err
	}

	output, err = isolateutil.Run(ctx, file.Name())
	if err != nil {
		return output, err
	}
	return output, nil
}

func NewService(l *zap.SugaredLogger) Service {
	return &runService{
		logger: l,
	}
}
