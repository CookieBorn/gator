package config

func Working(x string) string {
	return x + "Works"
}

type Config struct {
	DB_URL   string `json:"db_url"`
	Username string `json:"current_user_name"`
}
