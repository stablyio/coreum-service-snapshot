package utils

type Stage string

const (
	Prod  Stage = "prod"
	Beta  Stage = "beta"
	Local Stage = "local"
	Test  Stage = "test"
)

func (s Stage) String() string {
	return string(s)
}

func GetStage() Stage {
	stageString := GetEnvVar("STAGE", true)
	switch stageString {
	case Prod.String():
		return Prod
	case Beta.String():
		return Beta
	case Local.String():
		return Local
	case Test.String():
		return Test
	default:
		panic("Invalid STAGE environment variable: " + stageString)
	}
}

func GetAllStages() []Stage {
	return []Stage{Prod, Beta, Local, Test}
}
