package main

import (
	"flag"
	"github.com/ameamezhou/distributekv/config"
	"github.com/ameamezhou/distributekv/dbpkg"
	"github.com/ameamezhou/distributekv/handlers"
	"github.com/ameamezhou/distributekv/xlog"
	"github.com/BurntSushi/toml"
	"net/http"
)
var (
	// db文件的本地地址，也就是要初始化一个 xx.db 文件
	dbLocation = flag.String("db-location","", "the path to the bold db database")
	httpAddr = flag.String("http-addr", "127.0.0.1:9090", "HTTP host and port")
	configFile = flag.String("config-file", "sharding.toml", "config file for static sharding")
	shard = flag.String("shard", "", "the name of the shard for the data")
)

func parseFlags(){
	flag.Parse()

	if *dbLocation == "" {
		xlog.Fatal("Must provide db-location")
	}

	if *shard == "" {
		xlog.Fatalf("Must Provide shard")
	}
}

func main(){
	parseFlags()

	//configContents, err := os.ReadFile(*configFile)
	//if err != nil {
	//	xlog.Fatalf("ReadFile(%q): %v", *configFile, err)
	//}
	//var c config.Config
	//if _, err := toml.DecodeFile(*configFile, &c); err != nil {
	//	xlog.Fatalf("toml.DecodeFile(%q): %v", *configFile, err)
	//}
	c, err := config.ParseFile(*configFile)
	if err != nil {
		return
	}

	xlog.Infof("%#v", &c)

	//targetShardIndex := -1
	//shardCount := len(c.Shard)
	//var addrs = make(map[int]string)
	//for _, v := range c.Shard {
	//	addrs[v.Index] = v.Address
	//	if v.Name == *shard {
	//		targetShardIndex = v.Index
	//	}
	//}
	//if targetShardIndex < 0 {
	//	xlog.Fatalf("shard %q was not found", *shard)
	//}
	//xlog.Infof("shard count is %d, current shard: %d", shardCount, targetShardIndex)

	shards, err := config.ParseShards(c.Shard,*shard)
	if err != nil {
		return
	}

	db, err := dbpkg.NewDataBase(dbLocation)
	if err != nil {
		return
	}
	defer db.Close()

	svr := handlers.NewServer(db, shards)
	// 到这里之后  我们就要考虑要怎么添加API了
	// 我们要考虑 如果一个人来使用我们的东西会怎么想 怎么做
	http.HandleFunc("/get", svr.GetHandler)
	http.HandleFunc("/set", svr.SetHandler)

	// hash(key) % <count> = <current index>


	xlog.Fatal(svr.ListenAndServe(*httpAddr))
}