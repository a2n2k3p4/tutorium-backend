package dbserver

func NewConfig() *Config {
	return &Config{
		// Input your connection details here...
		DBUser:     "postgres",
		DBPassword: "2546",
		DBHost:     "localhost",
		DBPort:     "5432",
		DBName:     "tutorium",
		//
	}
}
