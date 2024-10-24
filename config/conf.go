package config


// Shard describes a shard that holds the appropriate set of keys.
// Each shard has unique set of keys.
type Shard struct {
	Name string 	`toml:"name"`
	Index int		`toml:"idx"`
	Address string  `toml:"address"`
}

// Config is a describes the sharding config
type Config struct {
	Shard []Shard
}