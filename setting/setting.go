package setting

const (
	mb                                   = 1024 * 1024
	DEFAULT_PART_SIZE              int64 = 5 * mb // 5 MB
	DEFAULT_CONCURRENCY                  = 10
	DEFAULT_MAX_TRIES                    = 3
	DEFAULT_SIMMULATANOUS_DOWNLOAD       = 1
	DEFAULT_SAVE_LOCATION                = "./tmp"
)

type Setting struct {
	Partsize        int64
	Concurrency     int
	Maxtries        int
	SimmultanousNum int
	SaveLocation    string
}

func New() Setting {
	return Setting{
		Partsize:        DEFAULT_PART_SIZE,
		Concurrency:     DEFAULT_CONCURRENCY,
		Maxtries:        DEFAULT_MAX_TRIES,
		SimmultanousNum: DEFAULT_SIMMULATANOUS_DOWNLOAD,
		SaveLocation:    DEFAULT_SAVE_LOCATION,
	}
}
