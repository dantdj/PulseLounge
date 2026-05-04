package media

type Store interface {
	Save(fileName string, data []byte) (string, error)
	Exists(fileName string) (bool, error)
	PublicURL(fileName string) string
}
