package terminal

type Terminal interface {
	Start() error
	Run(cmd string) (string, error)
	Close() error
}
