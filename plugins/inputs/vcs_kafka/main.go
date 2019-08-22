package vcs_kafka
import (
"flag"
"fmt"
"io"
"os"
"regexp"
"strconv"
"strings"
"sync"
"time"
"github.com/Shopify/sarama"
)
const (
	allPartitionsHuman = "all"
)
type groupCmd struct {
	brokers    []string
	group      string
	filter     *regexp.Regexp
	topic      string
	partitions []int32
	verbose    bool
	detail     bool
	version    sarama.KafkaVersion
	client     sarama.Client
}
type group struct {
	Name    string
	Topic   string
	Offsets []groupOffset
}
type groupOffset struct {
	Partition int32
	Offset    *int64
	Lag       *int64
}
type printContext struct {
	output interface{}
	done   chan struct{}
}
func (cmd *groupCmd) run(args []string) {
	var err error
	cmd.parseArgs(args)
	if cmd.client, err = sarama.NewClient(cmd.brokers, cmd.saramaConfig()); err != nil {
		failf("failed to create client err=%v", err)
	}
	brokers := cmd.client.Brokers()
	if cmd.verbose {
		var addrs []string
		for _, b := range brokers {
			addrs = append(addrs, b.Addr())
		}
		fmt.Fprintf(os.Stderr, "found brokers: %v\n", addrs)
	}
	groups := []string{cmd.group}
	if cmd.group == "" {
		groups = []string{}
		for _, g := range cmd.findGroups(brokers) {
			if cmd.filter.MatchString(g) {
				groups = append(groups, g)
			}
		}
	}
	if cmd.verbose {
		fmt.Fprintf(os.Stderr, "found groups: %v\n", groups)
	}
	topics := []string{cmd.topic}
	if cmd.topic == "" {
		topics = cmd.fetchTopics()
	}
	if cmd.verbose {
		fmt.Fprintf(os.Stderr, "found topics: %v\n", topics)
	}
	topicPartitions := map[string][]int32{}
	for _, topic := range topics {
		parts := cmd.partitions
		if len(parts) == 0 {
			parts = cmd.fetchPartitions(topic)
			if cmd.verbose {
				fmt.Fprintf(os.Stderr, "found partitions=%v for topic=%v\n", parts, topic)
			}
		}
		topicPartitions[topic] = parts
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(groups) * len(topics))
	for _, grp := range groups {
		for top, parts := range topicPartitions {
			go func(grp, topic string, partitions []int32) {
				cmd.printGroupTopicOffset(grp, topic, partitions)
				wg.Done()
			}(grp, top, parts)
		}
	}
	wg.Wait()
}
func (cmd *groupCmd) printGroupTopicOffset(grp, top string, parts []int32) {
	var lag int64
	var partLagMin int64 = 9223372036854775807 //max int64
	var partLagMax int64
	ok := false
	now := time.Now()
	for _, part := range parts {
		gOff := cmd.fetchGroupOffset(grp, top, part)
		if gOff.Lag != nil && gOff.Offset != nil {
			ok = true
			gLag := int64(*gOff.Lag)
			gOffset := int64(*gOff.Offset)
			lag += gLag
			if gLag < partLagMin {
				partLagMin = gLag
			}
			if gLag > partLagMax {
				partLagMax = gLag
			}
			if cmd.detail {
				tags := map[string]string{}
				fields := make(map[string]interface{})
				m := NewMetric("kafka_lag_detail", tags, fields, now)
				m.AddTag("group", grp)
				m.AddTag("partition", strconv.Itoa(int(part)))
				m.AddTag("topic", top)
				m.AddField("current_offset", gOffset)
				m.AddField("lag", gLag)
				fmt.Println(m)
			}
		}
	}
	if ok && !cmd.detail {
		tags := map[string]string{}
		fields := make(map[string]interface{})
		m := NewMetric("kafka_lag", tags, fields, now)
		m.AddTag("group", grp)
		m.AddTag("topic", top)
		m.AddField("min_lag", partLagMin)
		m.AddField("max_lag", partLagMax)
		m.AddField("lag", lag)
		fmt.Println(m)
	}
}
func (cmd *groupCmd) resolveOffset(top string, part int32, off int64) int64 {
	resolvedOff, err := cmd.client.GetOffset(top, part, off)
	if err != nil {
		failf("failed to get offset to reset to for partition=%d err=%v", part, err)
	}
	// if cmd.verbose {
	// 	fmt.Fprintf(os.Stderr, "resolved offset %v for topic=%s partition=%d to %v\n", off, top, part, resolvedOff)
	// }
	return resolvedOff
}
func (cmd *groupCmd) fetchGroupOffset(grp string, top string, part int32) groupOffset {
	var (
		err           error
		offsetManager sarama.OffsetManager
	)
	// if cmd.verbose {
	// 	fmt.Fprintf(os.Stderr, "fetching offset information for group=%v topic=%v partition=%v\n", grp, top, part)
	// }
	if offsetManager, err = sarama.NewOffsetManagerFromClient(grp, cmd.client); err != nil {
		failf("failed to create client err=%v", err)
	}
	defer logClose("offset manager", offsetManager)
	pom, err := offsetManager.ManagePartition(top, part)
	if err != nil {
		failf("failed to manage partition group=%s topic=%s partition=%d err=%v", grp, top, part, err)
	}
	defer logClose("partition offset manager", pom)
	groupOff, _ := pom.NextOffset()
	// we haven't reset it, and it wasn't set before - lag depends on client's config
	if groupOff == sarama.OffsetNewest || groupOff == sarama.OffsetOldest {
		return groupOffset{Partition: part}
	}
	partOff := cmd.resolveOffset(top, part, sarama.OffsetNewest)
	lag := partOff - groupOff
	return groupOffset{Partition: part, Offset: &groupOff, Lag: &lag}
}
func (cmd *groupCmd) fetchTopics() []string {
	tps, err := cmd.client.Topics()
	if err != nil {
		failf("failed to read topics err=%v", err)
	}
	return tps
}
func (cmd *groupCmd) fetchPartitions(top string) []int32 {
	ps, err := cmd.client.Partitions(top)
	if err != nil {
		failf("failed to read partitions for topic=%s err=%v", top, err)
	}
	return ps
}
type findGroupResult struct {
	done  bool
	group string
}
func (cmd *groupCmd) findGroups(brokers []*sarama.Broker) []string {
	var (
		doneCount int
		groups    = []string{}
		results   = make(chan findGroupResult)
		errs      = make(chan error)
	)
	for _, broker := range brokers {
		go cmd.findGroupsOnBroker(broker, results, errs)
	}
awaitGroups:
	for {
		if doneCount == len(brokers) {
			return groups
		}
		select {
		case err := <-errs:
			failf("failed to find groups err=%v", err)
		case res := <-results:
			if res.done {
				doneCount++
				continue awaitGroups
			}
			groups = append(groups, res.group)
		}
	}
}
func (cmd *groupCmd) findGroupsOnBroker(broker *sarama.Broker, results chan findGroupResult, errs chan error) {
	var (
		err  error
		resp *sarama.ListGroupsResponse
	)
	if err = cmd.connect(broker); err != nil {
		errs <- fmt.Errorf("Failed to connect to broker %#v err=%s\n", broker.Addr(), err)
	}
	if resp, err = broker.ListGroups(&sarama.ListGroupsRequest{}); err != nil {
		errs <- fmt.Errorf("Failed to list brokers on %#v err=%v", broker.Addr(), err)
	}
	if resp.Err != sarama.ErrNoError {
		errs <- fmt.Errorf("Failed to list brokers on %#v err=%v", broker.Addr(), resp.Err)
	}
	for name := range resp.Groups {
		results <- findGroupResult{group: name}
	}
	results <- findGroupResult{done: true}
}
func (cmd *groupCmd) connect(broker *sarama.Broker) error {
	if ok, _ := broker.Connected(); ok {
		return nil
	}
	if err := broker.Open(cmd.saramaConfig()); err != nil {
		return err
	}
	connected, err := broker.Connected()
	if err != nil {
		return err
	}
	if !connected {
		return fmt.Errorf("Failed to connect broker %#v", broker.Addr())
	}
	return nil
}
func (cmd *groupCmd) saramaConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = cmd.version
	cfg.ClientID = "kafka-monitor-tool"
	return cfg
}
func (cmd *groupCmd) parseArgs(as []string) {
	var (
		err  error
		args = cmd.parseFlags(as)
	)
	cmd.topic = args.topic
	cmd.group = args.group
	cmd.verbose = args.verbose
	cmd.detail = args.detail
	cmd.version = kafkaVersion(args.version)
	switch args.partitions {
	case "", "all":
		cmd.partitions = []int32{}
	default:
		pss := strings.Split(args.partitions, ",")
		for _, ps := range pss {
			p, err := strconv.ParseInt(ps, 10, 32)
			if err != nil {
				failf("partition id invalid err=%v", err)
			}
			cmd.partitions = append(cmd.partitions, int32(p))
		}
	}
	if cmd.partitions == nil {
		failf(`failed to interpret partitions flag %#v. Should be a comma separated list of partitions or "all".`, args.partitions)
	}
	if cmd.filter, err = regexp.Compile(args.filter); err != nil {
		failf("filter regexp invalid err=%v", err)
	}
	if args.brokers == "" {
		args.brokers = "localhost:9092"
	}
	cmd.brokers = strings.Split(args.brokers, ",")
	for i, b := range cmd.brokers {
		if !strings.Contains(b, ":") {
			cmd.brokers[i] = b + ":9092"
		}
	}
}
type groupArgs struct {
	topic      string
	brokers    string
	partitions string
	group      string
	filter     string
	verbose    bool
	detail     bool
	version    string
}
func (cmd *groupCmd) parseFlags(as []string) groupArgs {
	var args groupArgs
	flags := flag.NewFlagSet("group", flag.ExitOnError)
	flags.StringVar(&args.brokers, "brokers", "", "Comma separated list of brokers. Port defaults to 9092 when omitted (defaults to localhost:9092).")flags.StringVar(&args.group, "group", "", "Consumer group name.")
	flags.StringVar(&args.filter, "filter", "", "Regex to filter groups.")
	flags.StringVar(&args.topic, "topic", "", "Topic to consume")
	flags.BoolVar(&args.verbose, "verbose", false, "More verbose logging.")
	flags.BoolVar(&args.detail, "detail", false, "Print detail for all partitions.")
	flags.StringVar(&args.version, "version", "", "Kafka protocol version")
	flags.StringVar(&args.partitions, "partitions", allPartitionsHuman, "comma separated list of partitions to limit offsets to, or all")
	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of group:")
		flags.PrintDefaults()
		os.Exit(2)
	}
	_ = flags.Parse(as)
	return args
}
func main() {
	cmd := &groupCmd{}
	cmd.run(os.Args[1:])
}
func kafkaVersion(s string) sarama.KafkaVersion {
	dflt := sarama.V0_10_0_0
	switch s {
	case "v0.8.2.0":
		return sarama.V0_8_2_0
	case "v0.8.2.1":
		return sarama.V0_8_2_1
	case "v0.8.2.2":
		return sarama.V0_8_2_2
	case "v0.9.0.0":
		return sarama.V0_9_0_0
	case "v0.9.0.1":
		return sarama.V0_9_0_1
	case "v0.10.0.0":
		return sarama.V0_10_0_0
	case "v0.10.0.1":
		return sarama.V0_10_0_1
	case "v0.10.1.0":
		return sarama.V0_10_1_0
	case "v0.10.1.1":
		return sarama.V0_10_1_1
	case "v0.10.2.0":
		return sarama.V0_10_2_0
	case "v0.10.2.1":
		return sarama.V0_10_2_1
	case "v0.11.0.0":
		return sarama.V0_11_0_0
	case "v0.11.0.1":
		return sarama.V0_11_0_1
	case "v0.11.0.2":
		return sarama.V0_11_0_2
	case "v1.0.0.0":
		return sarama.V1_0_0_0
	case "v1.1.0.0":
		return sarama.V1_1_0_0
	case "":
		return dflt
	}
	failf("unsupported kafka version %#v - supported: v0.8.2.0, v0.8.2.1, v0.8.2.2, v0.9.0.0, v0.9.0.1, v0.10.0.0, v0.10.0.1, v0.10.1.0, v0.10.1.1, v0.10.2.0, v0.10.2.1, v0.11.0.0, v0.11.0.1, v0.11.0.2, v1.0.0.0, v1.1.0.0", s)
	return dflt
}
func failf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
func logClose(name string, c io.Closer) {
	if err := c.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to close %#v err=%v", name, err)
	}
}
