package kaptcha

type kaptcha interface {
	Gen() string
	Get(key string) string
	Save(key,code string) error
	Auth(key,code string) error
	Fire(msg string) error
}
