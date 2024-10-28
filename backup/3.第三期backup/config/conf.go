package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/ameamezhou/distributekv/xlog"
	"hash/fnv"
)

// Shard describes a shard that holds the appropriate set of keys.
// shard represents an easier-to-use representation of
// the sharding config:  the shards count, current index and
// the addresses of all other shards too.
// Each shard has unique set of keys.

type Shard struct {
	Name string 	`toml:"name"`
	Index int		`toml:"idx"`
	Address string  `toml:"address"`
	// Addrs map[int]string `toml:"addrs"`
}

type Shards struct {
	Count int
	CurIndex int
	Addrs map[int]string
}

// Config is a describes the sharding config
type Config struct {
	Shard []Shard
}

func ParseFile(fileName string) (Config, error){
	var c Config
	if _, err := toml.DecodeFile(fileName, &c); err != nil {
		xlog.Errorf("toml.DecodeFile(%q): %v", fileName, err)
		return Config{}, err
	}
	return c, nil
}

// ParseShards converts and verifies the list of shards
// specified in the config into a form that can be used
// for routing
func ParseShards(shards []Shard, curShardName string) (*Shards, error) {
	shardCount := len(shards)
	shardIndex := -1
	addrs := make(map[int]string)

	for _, s := range shards {
		if _, ok := addrs[s.Index]; ok {
			return nil, fmt.Errorf("duplicate shard index: %d", s.Index)
		}

		addrs[s.Index] = s.Address
		if s.Name == curShardName {
			shardIndex = s.Index
		}
	}
	for i := 0; i < shardCount; i++ {
		if _, ok := addrs[i]; !ok {
			return nil, fmt.Errorf("shard %d is not found", i)
		}
	}
	if shardIndex < 0 {
		return nil, fmt.Errorf("shard %q was not found", curShardName)
	}

	return &Shards{
		Addrs: addrs,
		Count: shardCount,
		CurIndex: shardIndex,
	}, nil
}

func (s *Shards) Index (key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.Count))
}