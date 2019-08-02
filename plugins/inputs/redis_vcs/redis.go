package redis_vcs

import (
	"bufio"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/internal/tls"
	"github.com/influxdata/telegraf/plugins/inputs"
	"io"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Redis struct {
	Server   string
	Dbs      []string
	Keys     []string // 1.input -> client-db1:input .
	Password string
	tls.ClientConfig

	clients     []Client
	initialized bool
}

type Client interface {
	Info() *redis.StringCmd
	BaseTags() map[string]string
	GetKeyLen([]string) (map[string]int64, error)
}

type RedisClient struct {
	client *redis.Client
	tags   map[string]string
}
// input queue1. queue2,...
func (r *RedisClient) GetKeyLen(keys []string ) (map[string]int64, error){
	result := make(map[string]int64)
	var err error
	for _, key := range keys {
		result[key], err = r.client.LLen(key).Result()
	}
	return result, err
}

func (r *RedisClient) Info() *redis.StringCmd {
	return r.client.Info()
}


func (r *RedisClient) BaseTags() map[string]string {
	tags := make(map[string]string)
	for k, v := range r.tags {
		tags[k] = v
	}
	return tags
}

var sampleConfig = `
  ## specify servers via a url matching:
  ##  [protocol://][:password]@address[:port]
  ##  e.g.
  ##    tcp://localhost:6379
  ##    tcp://:password@192.168.99.100
  ##    unix:///var/run/redis.sock
  ##
  ## If no servers are specified, then localhost is used as the host.
  ## If no port is specified, 6379 is used
  server = "tcp://localhost:6379"
  dbs = ["1", "2"]
  keys = ["1.input_queue", "2.input_queue"]

  ## specify server password
  # password = "s#cr@t%"

  ## Optional TLS Config
  # tls_ca = "/etc/telegraf/ca.pem"
  # tls_cert = "/etc/telegraf/cert.pem"
  # tls_key = "/etc/telegraf/key.pem"
  ## Use TLS but skip chain & host verification
  # insecure_skip_verify = true
`

func (r *Redis) SampleConfig() string {
	return sampleConfig
}

func (r *Redis) Description() string {
	return "Read metrics from one or many redis servers"
}

var Tracking = map[string]string{
	"uptime_in_seconds": "uptime",
	"connected_clients": "clients",
	"role":              "replication_role",
}

func (r *Redis) init(acc telegraf.Accumulator) error {
	if r.initialized {
		return nil
	}
	if r.Server == "" {
		r.Server = "tcp://localhost:6379"
	}
	if !strings.HasPrefix(r.Server, "tcp://") && !strings.HasPrefix(r.Server, "unix://") {
		log.Printf("W! [inputs.redis]: server URL found without scheme; please update your configuration file")
		r.Server = "tcp://" + r.Server
	}
	u, err := url.Parse(r.Server)
	if err != nil {
		return fmt.Errorf("Unable to parse to address  %v", err)
	}
	password := ""
	if u.User != nil {
		pw, ok := u.User.Password()
		if ok {
			password = pw
		}
	}
	if len(r.Password) > 0 {
		password = r.Password
	}

	var address string
	if u.Scheme == "unix" {
		address = u.Path
	} else {
		address = u.Host
	}

	tlsConfig, err := r.ClientConfig.TLSConfig()
	if err != nil {
		return err
	}


	r.clients = make([]Client, len(r.Dbs))
	for i, db := range r.Dbs {
		intDB, err := strconv.Atoi(db)
		if err != nil {
			return err
		}
		client := redis.NewClient(
			&redis.Options{
				Addr:      address,
				Password:  password,
				Network:   u.Scheme,
				PoolSize:  1,
				TLSConfig: tlsConfig,
				DB: intDB,
			},
		)
		tags := map[string]string{}
		if u.Scheme == "unix" {
			tags["socket"] = u.Path
		} else {
			tags["server"] = u.Hostname()
			tags["port"] = u.Port()
		}
		tags["DB"] = db
		r.clients[i] = &RedisClient{
			client: client,
			tags:   tags,
		}
		r.initialized = true
	}
	return nil
}

// Reads stats from all configured servers accumulates stats.
// Returns one of the errors encountered while gather stats (if any).
func (r *Redis) Gather(acc telegraf.Accumulator) error {
	if !r.initialized {
		err := r.init(acc)
		if err != nil {
			return err
		}
	}

	var wg sync.WaitGroup
	for _, client := range r.clients {
		keys := getKey(r.Keys,client.BaseTags()["DB"])
		if keys == nil{
			continue
		}
		wg.Add(1)
		go func(client Client) {
			defer wg.Done()
			acc.AddError(r.gatherDB(client, acc, keys))
		}(client)
	}
	wg.Wait()
	acc.AddError(r.gatherServer(r.clients[0], acc))
	return nil
}
// input ([1.queue1, 2.queue2,1.queue3, 2.queue1..], 1) output (queue1, queue3)
func getKey(userInput []string, DB string) []string{
	var result []string
	for _, v := range userInput {
		r := regexp.MustCompile(`(.+?).(.+)`)
		arr :=  r.FindStringSubmatch(v)
		if arr[1]== DB{
			result = append(result, arr[2])
		}
	}
	return result
}

func (r *Redis) gatherServer(client Client, acc telegraf.Accumulator) error {
	info, err := client.Info().Result()
	if err != nil {
		return err
	}

	rdr := strings.NewReader(info)
	return gatherInfoOutputOS(rdr, acc, client.BaseTags())
}

func (r *Redis) gatherDB(client Client, acc telegraf.Accumulator, keys []string) error {

	//info, err := client.Info().Result()

	info, err := client.GetKeyLen(keys)
	if err != nil {
		return err
	}
	return gatherInfoOutputDB(info, acc, client.BaseTags())
}
// gatherInfoOutput gathers
func gatherInfoOutputOS(
	rdr io.Reader,
	acc telegraf.Accumulator,
	tags map[string]string,
) error {
	var section string
	var keyspace_hits, keyspace_misses int64

	scanner := bufio.NewScanner(rdr)
	fields := make(map[string]interface{})
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			if len(line) > 2 {
				section = line[2:]
			}
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}
		name := string(parts[0])

		if section == "Server" {
			if name != "lru_clock" && name != "uptime_in_seconds" && name != "redis_version" {
				continue
			}
		}

		if strings.HasPrefix(name, "master_replid") {
			continue
		}

		if name == "mem_allocator" {
			continue
		}

		if strings.HasSuffix(name, "_human") {
			continue
		}

		metric, ok := Tracking[name]
		if !ok {
			if section == "Keyspace" {
				kline := strings.TrimSpace(string(parts[1]))
				gatherKeyspaceLine(name, kline, acc, tags)
				continue
			}
			metric = name
		}

		val := strings.TrimSpace(parts[1])

		// Some percentage values have a "%" suffix that we need to get rid of before int/float conversion
		val = strings.TrimSuffix(val, "%")

		// Try parsing as int
		if ival, err := strconv.ParseInt(val, 10, 64); err == nil {
			switch name {
			case "keyspace_hits":
				keyspace_hits = ival
			case "keyspace_misses":
				keyspace_misses = ival
			case "rdb_last_save_time":
				// influxdb can't calculate this, so we have to do it
				fields["rdb_last_save_time_elapsed"] = time.Now().Unix() - ival
			}
			fields[metric] = ival
			continue
		}

		// Try parsing as a float
		if fval, err := strconv.ParseFloat(val, 64); err == nil {
			fields[metric] = fval
			continue
		}

		// Treat it as a string

		if name == "role" {
			tags["replication_role"] = val
			continue
		}

		fields[metric] = val
	}
	var keyspace_hitrate float64 = 0.0
	if keyspace_hits != 0 || keyspace_misses != 0 {
		keyspace_hitrate = float64(keyspace_hits) / float64(keyspace_hits+keyspace_misses)
	}
	fields["keyspace_hitrate"] = keyspace_hitrate
	acc.AddFields("redis", fields, tags)
	return nil
}
// gatherInfoOutput gathers
func gatherInfoOutputDB(
	info map[string]int64, // info[<queue_name>] = leng_queue
	acc telegraf.Accumulator,
	tags map[string]string,
) error {
	for k, v := range info {
		tags["queue_name"]= k
		fields := map[string]interface{}{
			"num_message": v,
		}
		acc.AddFields("redis_queuespace", fields, tags)
	}
	return nil
}

// Parse the special Keyspace line at end of redis stats
// This is a special line that looks something like:
//     db0:keys=2,expires=0,avg_ttl=0
// And there is one for each db on the redis instance
func gatherKeyspaceLine(
	name string,
	line string,
	acc telegraf.Accumulator,
	global_tags map[string]string,
) {
	if strings.Contains(line, "keys=") {
		fields := make(map[string]interface{})
		tags := make(map[string]string)
		for k, v := range global_tags {
			tags[k] = v
		}
		tags["database"] = name
		dbparts := strings.Split(line, ",")
		for _, dbp := range dbparts {
			kv := strings.Split(dbp, "=")
			ival, err := strconv.ParseInt(kv[1], 10, 64)
			if err == nil {
				fields[kv[0]] = ival
			}
		}
		acc.AddFields("redis_keyspace", fields, tags)
	}
}

func init() {
	inputs.Add("redis_vcs", func() telegraf.Input {
		return &Redis{}
	})
}
