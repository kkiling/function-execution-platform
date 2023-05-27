package model

type TemplateFile struct {
	FilePath string
	Body     string
}

type ContainerParams struct {
	MemoryLimitMb       int
	MemoryReservationMb int
	DiskSizeMb          int
	CPULimit            float32
	CPUReservation      float32
	TimeoutSec          int
}

type Template struct {
	Name               string
	Language           string
	Version            string
	Description        string
	CopyFiles          []string
	ForbiddenFileNames []string
	ReadOnlyFile       []string
	ContainerParams    ContainerParams
	Files              []TemplateFile
}

type TemplateMetaData struct {
	MetaData
	Template
}
