package goid

type GoId interface {
	Gen() string
	GenWithLength(length int) string
	NewUUID() string
}
