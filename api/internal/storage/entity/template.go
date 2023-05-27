package entity

type TemplateFile struct {
	FilePath string `bson:"filePath"`
	Body     string `bson:"body"`
}

type ContainerParams struct {
	MemoryLimitMb       int     `bson:"memory_limit_mb"`
	MemoryReservationMb int     `bson:"memory_reservation_mb"`
	DiskSizeMb          int     `bson:"disk_size_mb"`
	CPULimit            float32 `bson:"cpu_limit"`
	CPUReservation      float32 `bson:"cpu_reservation"`
	TimeoutSec          int     `bson:"timeout_sec"`
}

type Template struct {
	Name               string          `bson:"name"`
	Language           string          `bson:"language"`
	Version            string          `bson:"version"`
	Description        string          `bson:"description"`
	CopyFiles          []string        `bson:"copyFiles"`
	ForbiddenFileNames []string        `bson:"forbiddenFileNames"`
	ReadOnlyFile       []string        `bson:"readOnlyFile"`
	ContainerParams    ContainerParams `bson:"containerParams"`
	Files              []TemplateFile  `bson:"files"`
}

type TemplateMetaData struct {
	MetaData `bson:",inline"`
	Template `bson:",inline"`
}
