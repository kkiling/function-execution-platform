package template_impl

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/go-version"
	"github.com/kkiling/function-execution-platform/api/internal/common"
	"github.com/kkiling/function-execution-platform/api/internal/factory"
	"github.com/kkiling/function-execution-platform/api/internal/service"
	"github.com/kkiling/function-execution-platform/api/internal/service/model"
	"github.com/kkiling/function-execution-platform/api/internal/service/tostorage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	rootPath1   = "../template/"
	packageFile = "package.json"
)

type TemplateService struct {
	fact            factory.IScopeFactory
	templateStorage tostorage.ITemplateStorage
}

func NewService(fact factory.IScopeFactory) *TemplateService {
	return &TemplateService{
		fact:            fact,
		templateStorage: tostorage.NewTemplateStorage(fact),
	}
}

func (t TemplateService) loadDefaultContainerParams() model.ContainerParams {
	cfg := t.fact.GetConfig().DefaultContainerParams
	return model.ContainerParams{
		MemoryLimitMb:       cfg.MemoryLimitMb,
		MemoryReservationMb: cfg.MemoryReservationMb,
		DiskSizeMb:          cfg.DiskSizeMb,
		CPULimit:            cfg.CPULimit,
		CPUReservation:      cfg.CPUReservation,
		TimeoutSec:          cfg.TimeoutSec,
	}
}

func (t TemplateService) loadPackagePaths(dir string) ([]string, error) {
	packagePaths := make([]string, 0)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Base(path) == packageFile {
			packagePaths = append(packagePaths, path)
		}
		return nil
	})

	if err != nil {
		return packagePaths, errors.Wrap(err, "fail read package.json")
	}

	return packagePaths, nil
}

func (t TemplateService) loadTemplateFiles(dir string) ([]model.TemplateFile, error) {
	var files []model.TemplateFile
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		bodyBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files = append(files, model.TemplateFile{
			FilePath: relPath,
			Body:     string(bodyBytes),
		})
		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "fail load files")
	}

	return files, nil
}

func (t TemplateService) loadTemplate(packagePath string) (*model.Template, error) {
	jsonFile, err := os.Open(packagePath)
	if err != nil {
		return nil, errors.Wrap(err, "fail open package.json")
	}
	defer common.Close(jsonFile)

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, errors.Wrap(err, "fail read json file")
	}

	var template model.Template
	err = json.Unmarshal(byteValue, &template)
	if err != nil {
		return nil, errors.Wrap(err, "fail unmarshal package.json")
	}

	files, err := t.loadTemplateFiles(filepath.Dir(packagePath))
	if err != nil {
		return nil, errors.Wrap(err, "fail load template files")
	}

	template.Files = files
	template.ContainerParams = t.loadDefaultContainerParams()
	return &template, nil
}

func newVersionGreaterThanOldVersion(oldVersion, newVersion string) (bool, error) {
	oldV, err := version.NewVersion(oldVersion)
	if err != nil {
		return false, errors.Wrapf(err, "fail compare %s version", oldVersion)
	}

	newV, err := version.NewVersion(newVersion)
	if err != nil {
		return false, errors.Wrapf(err, "fail compare %s version", newVersion)
	}

	return newV.GreaterThan(oldV), nil
}

func (t TemplateService) loadTemplatesFromDir(ctx context.Context, dir string) error {
	packagePaths, err := t.loadPackagePaths(dir)
	if err != nil {
		return service.MakeRuntimeWrapErr(err, "fail load packages")
	}

	templates := make([]*model.Template, 0, len(packagePaths))
	for _, packagePath := range packagePaths {
		template, err := t.loadTemplate(packagePath)
		if err != nil {
			return service.MakeRuntimeWrapErr(err, "fail load template")
		}
		templates = append(templates, template)
	}

	for _, template := range templates {
		find, err := t.templateStorage.FindTemplate(ctx, template.Name, template.Language)
		if err != nil {
			return service.MakeRuntimeWrapErr(err, "fail find template")
		}

		if find != nil {
			if greater, err := newVersionGreaterThanOldVersion(find.Version, template.Version); err != nil {
				return service.MakeRuntimeWrapErr(err, "fail compare version")
			} else if greater {
				if err := t.templateStorage.UpdateTemplate(ctx, find.Id, template); err != nil {
					return service.MakeRuntimeWrapErr(err, "fail save template")
				}
				log.Info().Msgf("update template %s/%s", template.Name, template.Language)
				break
			} else {
				log.Info().Msgf("template %s/%s already saved", template.Name, template.Language)
				break
			}
		}

		if _, err := t.templateStorage.SaveTemplate(ctx, template); err != nil {
			return service.MakeRuntimeWrapErr(err, "fail save template")
		}

		log.Info().Msgf("save template %s/%s", template.Name, template.Language)
	}

	return nil
}

func (t TemplateService) InitBaseTemplates(ctx context.Context) error {
	log.Info().Msg("load default template")
	return t.loadTemplatesFromDir(ctx, rootPath1)
}

func (t TemplateService) LoadGitTemplates(ctx context.Context, gitUrl, branch string) error {
	dir, err := os.MkdirTemp("", "clone_template_functions")
	if err != nil {
		return service.MakeRuntimeWrapErr(err, "fail create temp dir")
	}
	// defer os.RemoveAll(dir)

	keyPath := dir + "/id_rsa"
	pk := t.fact.GetConfig().Git.PrivateKey
	err = os.WriteFile(keyPath, []byte(pk), 0600)
	if err != nil {
		return service.MakeRuntimeWrapErr(err, "fail save private key")
	}

	dir += "/repo"
	cmd := exec.Command("git", "clone", gitUrl, dir)
	cmd.Env = append(os.Environ(), "GIT_SSH_COMMAND=ssh -i "+keyPath+" -o IdentitiesOnly=yes -F /dev/null")
	err = cmd.Run()
	if err != nil {
		return service.MakeRuntimeWrapErr(err, "fail clone git url")
	}

	cmd = exec.Command("git", "-C", dir, "checkout", branch)
	err = cmd.Run()
	if err != nil {
		return service.MakeRuntimeWrapErr(err, "fail checkout branch")
	}

	return t.loadTemplatesFromDir(ctx, dir)
}
