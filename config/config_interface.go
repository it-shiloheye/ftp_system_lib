package configuration

type ProcessType string

const (
	ProcessTypeClient ProcessType = "client"
	ProcessTypeServer ProcessType = "server"
)

type Config interface {
	Id() string               // unique id of this specific process
	Version() string          // git release version
	TempDir() string          // temporary directory to store files in transit
	DataPath() string         // permanent directory to stor files
	ProcessType() ProcessType // process type to differentiate client and server
}

type ServerConfig interface {
	Clients() []string // returns array of clients
}
