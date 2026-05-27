package media

type Store interface {
	Save(id string, contentType string, data []byte) (string, error)
	PublicURL(id string) string
}
