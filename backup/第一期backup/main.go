package main

import (
	"flag"
	"fmt"
	"github.com/ameamezhou/distributekv/dbpkg"
	"github.com/ameamezhou/distributekv/xlog"
	"net/http"
)
var (
	// db文件的本地地址，也就是要初始化一个 xx.db 文件
	dbLocation = flag.String("db-location","", "the path to the bold db database")
	httpAddr = flag.String("http-addr", "127.0.0.1:9090", "HTTP host and port")
)

func parseFlags(){
	flag.Parse()

	if *dbLocation == "" {
		xlog.Fatal("Must provide db-location")
	}
}

func main(){
	parseFlags()
	db, err := dbpkg.NewDataBase(dbLocation)
	if err != nil {
		return
	}
	defer db.Close()
	// 到这里之后  我们就要考虑要怎么添加API了
	// 我们要考虑 如果一个人来使用我们的东西会怎么想 怎么做
	http.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		key := request.Form.Get("key")
		value, err := db.GetValue(key)
		fmt.Fprintf(writer, "value = %q, err = %v", value, err)
	})

	http.HandleFunc("/set", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		key := request.Form.Get("key")
		value := request.Form.Get("value")
		err := db.SetKey(key, []byte(value))
		fmt.Fprintf(writer, "err = %v", err)
	})

	xlog.Fatal(http.ListenAndServe(*httpAddr, nil))
}