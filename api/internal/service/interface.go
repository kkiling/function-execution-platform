package service

import "context"

type ITemplateService interface {
	InitBaseTemplate(ctx context.Context) error
}
