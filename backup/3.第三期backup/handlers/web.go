package handlers

import (
	"fmt"
	"github.com/ameamezhou/distributekv/config"
	"github.com/ameamezhou/distributekv/dbpkg"
	"github.com/ameamezhou/distributekv/xlog"
	"io"
	"net/http"
)


// DbWebServer contains HTTP method handlers to be used for the database
type DbWebServer struct {
	db *dbpkg.DataBase
	//shardIndex int
	//shardCount int
	//address    map[int]string
	shards *config.Shards
}

// NewServer creates a new Server instance with HTTP handlers to be used to get and set Values.
func NewServer(db *dbpkg.DataBase, s *config.Shards) *DbWebServer {
	return &DbWebServer{
		db: db,
		shards: s,
	}
}

// 通过这种方式可以做 hash 分散存储到对应的 分片中
//func (s *DbWebServer) getShard (key string) int {
//	h := fnv.New64()
//	h.Write([]byte(key))
//	return int(h.Sum64() % uint64(s.shardCount))
//}

func (s *DbWebServer) redirect(writer http.ResponseWriter, r *http.Request, shard int) {
	url := "http://" + s.shards.Addrs[shard] + r.RequestURI
	xlog.Debugf("redirecting from shard %d to shard %d (%q)", s.shards.CurIndex, shard, url)
	resp, err := http.Get(url)
	if err != nil {
		writer.WriteHeader(500)
		xlog.Errorf("Error redirecting the request: %v", err)
		return
	}
	defer resp.Body.Close()
	io.Copy(writer, resp.Body)
}

// GetHandler handles "get" endpoint, read requests from the database.
func (s *DbWebServer) GetHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	key := request.Form.Get("key")

	shard := s.shards.Index(key)
	value, err := s.db.GetValue(key)

	// 如果非本分片，就转发到另一个分片
	if shard != s.shards.CurIndex {
		s.redirect(writer, request, shard)
		return
	}

	fmt.Fprintf(writer, "shard=%d, current shard = %d, addrs = %q, value = %q, err = %v", shard, s.shards.CurIndex, s.shards.Addrs[shard], value, err)
}

// SetHandler handles write requests from the database.
func (s *DbWebServer) SetHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	key := request.Form.Get("key")
	value := request.Form.Get("value")

	shard := s.shards.Index(key)
	// 如果非本分片，就转发到另一个分片
	if shard != s.shards.CurIndex {
		s.redirect(writer, request, shard)
		return
	}

	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(writer, "err = %v, current shard = %d, shardId = %d", err, s.shards.CurIndex, shard)
}

func (s *DbWebServer) ListenAndServe(addr string) error{
	xlog.Debug("Server start!")
	return http.ListenAndServe(addr, nil)
}