package utils

const (
	Prefix = "blog:" //项目key前缀
)

func GetRedisKey(key string) string {
	return Prefix + key
}
