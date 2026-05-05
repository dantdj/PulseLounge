package media

type Store interface {
	Save(id string, contentType string, data []byte) (string, error)
	Exists(id string) (bool, error)
	PublicURL(id string) string
}
