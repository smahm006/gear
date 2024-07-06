package connection

type Connection interface {
	Connect() error
	Close() error
	WhoAmI() (string, error)
	Execute() error
	CopyFile(string, string) error
	WriteData(string, string) error
}
