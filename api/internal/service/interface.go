package service

import "context"

type ITemplateService interface {
	InitBaseTemplates(ctx context.Context) error
	LoadGitTemplates(ctx context.Context, gitUrl, branch string) error
}
