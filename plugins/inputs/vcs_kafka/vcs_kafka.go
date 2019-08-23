package vcs_kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
)

const (
	clientID = "kafka-vcs-monitor"
)

var sampleConfig = `
  ## kafka servers
  brokers = ["localhost:9092"]
  ## Kafka protocol version
  version = "v0.8.2.0"
  ## Regex to filter groups
  filter_consummer_groups = ""
  ## Topic to consume
  topic = ""
  ## Print detail for all partitions or not
  detail = false
  ## Comma separated list of partitions to limit offsets to, or all
  partitions = "all"
  ## Print verbose debug
  debug = false
`

type Kafka struct {
	Brokers               []string
	FilterTopics          string
	FilterConsummerGroups string
	Detail                bool
	Version               string
	Debug                 bool

	filterTopics          *regexp.Regexp
	filterConsummerGroups *regexp.Regexp
	clientVersion         sarama.KafkaVersion
	conf                  *sarama.Config
	client                sarama.Client
}
type findGroupResult struct {
	done  bool
	group string
}

type groupOffset struct {
	Partition int32
	Offset    *int64
	Lag       *int64
}

type tagMap map[string]string
type fieldMap map[string]interface{}

func (Kafka) SampleConfig() string {
	return sampleConfig
}

func (k *Kafka) Description() string {
	return "Read metrics from Kafka topic(s)"
}

func (k *Kafka) parseArgs() error {
	var err error
	k.filterConsummerGroups, err = regexp.Compile(k.FilterConsummerGroups)
	if err != nil {
		return fmt.Errorf("config of filter consummer group regexp invalid err=%v", err)
	}
	k.filterTopics, err = regexp.Compile(k.FilterTopics)
	if err != nil {
		return fmt.Errorf("config of filter topics regexp invalid err=%v", err)
	}
	k.clientVersion, err = kafkaVersion(k.Version)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return err
}
func kafkaVersion(s string) (sarama.KafkaVersion, error) {
	var dflt sarama.KafkaVersion
	var err error
	switch s {
	case "v0.8.2.0":
		dflt = sarama.V0_8_2_0
	case "v0.8.2.1":
		dflt = sarama.V0_8_2_1
	case "v0.8.2.2":
		dflt = sarama.V0_8_2_2
	case "v0.9.0.0":
		dflt = sarama.V0_9_0_0
	case "v0.9.0.1":
		dflt = sarama.V0_9_0_1
	case "v0.10.0.0":
		dflt = sarama.V0_10_0_0
	case "v0.10.0.1":
		dflt = sarama.V0_10_0_1
	case "v0.10.1.0":
		dflt = sarama.V0_10_1_0
	case "v0.10.1.1":
		dflt = sarama.V0_10_1_1
	case "v0.10.2.0":
		dflt = sarama.V0_10_2_0
	case "v0.10.2.1":
		dflt = sarama.V0_10_2_1
	case "v0.11.0.0":
		dflt = sarama.V0_11_0_0
	case "v0.11.0.1":
		dflt = sarama.V0_11_0_1
	case "v0.11.0.2":
		dflt = sarama.V0_11_0_2
	case "v1.0.0.0":
		dflt = sarama.V1_0_0_0
	case "v1.1.0.0":
		dflt = sarama.V1_1_0_0
	default:
		err = fmt.Errorf("unsupported kafka version %#v - supported: v0.8.2.0, v0.8.2.1, v0.8.2.2, v0.9.0.0, v0.9.0.1, v0.10.0.0, v0.10.0.1, v0.10.1.0, v0.10.1.1, v0.10.2.0, v0.10.2.1, v0.11.0.0, v0.11.0.1, v0.11.0.2, v1.0.0.0, v1.1.0.0", s)
	}
	return dflt, err
}

//func (k *Kafka) connect(broker *sarama.Broker) error {
//	if ok, _ := broker.Connected(); ok {
//		return nil
//	}
//	cfg := sarama.NewConfig()
//	cfg.Version = cmd.version
//	cfg.ClientID = "kafka-monitor-tool"
//
//	if err := broker.Open(cmd.saramaConfig()); err != nil {
//		return err
//	}
//	connected, err := broker.Connected()
//	if err != nil {
//		return err
//	}
//	if !connected {
//		return fmt.Errorf("Failed to connect broker %#v", broker.Addr())
//	}
//	return nil
//}

//func (k *Kafka) findGroupsOnBroker(broker *sarama.Broker, results chan findGroupResult, errs chan error) {
//	var (
//		err  error
//		resp *sarama.ListGroupsResponse
//	)
//	if err = cmd.connect(broker); err != nil {
//		errs <- fmt.Errorf("Failed to connect to broker %#v err=%s\n", broker.Addr(), err)
//	}
//	if resp, err = broker.ListGroups(&sarama.ListGroupsRequest{}); err != nil {
//		errs <- fmt.Errorf("Failed to list brokers on %#v err=%v", broker.Addr(), err)
//	}
//	if resp.Err != sarama.ErrNoError {
//		errs <- fmt.Errorf("Failed to list brokers on %#v err=%v", broker.Addr(), resp.Err)
//	}
//	for name := range resp.Groups {
//		results <- findGroupResult{group: name}
//	}
//	results <- findGroupResult{done: true}
//}
//
//func (k *Kafka) findGroups(brokers []*sarama.Broker) []string {
//	var (
//		doneCount int
//		groups    = []string{}
//		results   = make(chan findGroupResult)
//		errs      = make(chan error)
//	)
//	for _, broker := range brokers {
//		go k.findGroupsOnBroker(broker, results, errs)
//	}
//awaitGroups:
//	for {
//		if doneCount == len(brokers) {
//			return groups
//		}
//		select {
//		case err := <-errs:
//			failf("failed to find groups err=%v", err)
//		case res := <-results:
//			if res.done {
//				doneCount++
//				continue awaitGroups
//			}
//			groups = append(groups, res.group)
//		}
//	}
//}

func saramaConfig(version sarama.KafkaVersion) *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = version
	cfg.ClientID = clientID
	return cfg
}
func (k *Kafka) connect(broker *sarama.Broker) error {
	if ok, _ := broker.Connected(); ok {
		return nil
	}
	if err := broker.Open(saramaConfig(k.clientVersion)); err != nil {
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

func (k *Kafka) findGroupsOnBroker(broker *sarama.Broker, results chan findGroupResult, errs chan error) {
	var (
		err  error
		resp *sarama.ListGroupsResponse
	)
	if err = k.connect(broker); err != nil {
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

func (k *Kafka) listGroups(brokers []*sarama.Broker) ([]string, error) {
	var (
		doneCount int
		groups    = []string{}
		results   = make(chan findGroupResult)
		errs      = make(chan error)
	)
	for _, broker := range brokers {
		go k.findGroupsOnBroker(broker, results, errs)
	}

awaitGroups:
	for {
		if doneCount == len(brokers) {
			return groups, nil
		}
		select {
		case err := <-errs:
			return nil, fmt.Errorf("failed to find groups err=%v", err)
		case res := <-results:
			if res.done {
				doneCount++
				continue awaitGroups
			}
			groups = append(groups, res.group)
		}
	}

}

func (k *Kafka) Gather(acc telegraf.Accumulator) error {
	err := k.parseArgs()
	if (err != nil) {
		return err
	}
	k.client, err = sarama.NewClient(k.Brokers, saramaConfig(k.clientVersion))
	if (err != nil) {
		return err
	}
	brokers := k.client.Brokers()
	if k.Debug {
		var addrs []string
		for _, b := range brokers {
			addrs = append(addrs, b.Addr())
		}
		log.Printf("D! Found brokers: %v\n", addrs)
	}

	allGroups, err := k.listGroups(brokers)
	if err != nil {
		return err
	}
	if k.Debug {
		log.Printf("D! all groups: %v\n", allGroups)
	}
	groups := []string{}
	for _, g := range allGroups {
		if k.filterConsummerGroups.MatchString(g) {
			groups = append(groups, g)
		}
	}
	if k.Debug {
		log.Printf("D! found groups: %v\n", groups)
	}

	allTopics, err := k.client.Topics()
	if err != nil {
		return err
	}
	if k.Debug {
		log.Printf("D! all groups: %v\n", allTopics)
	}
	topics := []string{}
	for _, t := range allTopics {
		if k.filterTopics.MatchString(t) {
			topics = append(groups, t)
		}
	}
	if k.Debug {
		log.Printf("D! found groups: %v\n", groups)
	}

	topicPartitions := map[string][]int32{}
	for _, topic := range topics {

		parts, err := k.client.Partitions(topic)
		if err != nil {
			return fmt.Errorf("fail to list partitions of topic %s %v", topic, err)
		}
		if k.Debug {
			fmt.Fprintf(os.Stderr, "found partitions=%v for topic=%v\n", parts, topic)
		}
		topicPartitions[topic] = parts
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(groups) * len(topics))
	for _, grp := range groups {
		for top, parts := range topicPartitions {
			go func(grp string, topic string, partitions []int32) {
				if topic == "__consumer_offsets" { // ingor consummer_offsets
					wg.Done()
					return
				}

				kafkalag, KafkaLagDetail, err := k.getKafkaLag(grp, topic, partitions)

				if err != nil{
					wg.Done()
					return
				}
				acc.AddFields("kafka_lag", fieldMap{
					"min_lag": kafkalag.min_lag,
					"max_lag": kafkalag.max_lag,
				}, tagMap{
					"group": kafkalag.consummerGroup,
					"topic": kafkalag.topic,})

				for _, v := range KafkaLagDetail{
					acc.AddFields("kafka_lag_detail", fieldMap{
						"current_offset": v.current_offset,
						"lag": v.lag,
					}, tagMap{
						"group": v.consummerGroup,
						"topic": v.topic,
						"partition": v.partition,
					})
				}

				wg.Done()
			}(grp, top, parts)
		}
	}
	wg.Wait()

	return nil
}

type KafkaLag struct {
	consummerGroup string
	topic          string
	min_lag        int64
	max_lag        int64
	lag            int64
}

type KafkaLagDetail struct {
	consummerGroup string
	topic          string
	partition      string
	current_offset int64
	lag            int64
}

func (k *Kafka) getKafkaLag(consummerGroup string, top string, parts []int32) (KafkaLag, []KafkaLagDetail, error) {
	var lag int64
	var partLagMin int64 = 9223372036854775807 //max int64
	var partLagMax int64
	ok := false
	var kafkaLag KafkaLag
	var kafkaLagDetails []KafkaLagDetail
	for _, part := range parts {
		gOff, err := k.fetchGroupOffset(consummerGroup, top, part)
		if err != nil {
			return KafkaLag{}, nil, err
		}
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
			if k.Detail {
				kafkaLagDetails = append(kafkaLagDetails, KafkaLagDetail{
					consummerGroup: consummerGroup,
					partition:      strconv.Itoa(int(part)),
					topic:          top,
					current_offset: gOffset,
					lag:            gLag,
				})
			}
		}
	}
	if ok {
		kafkaLag = KafkaLag{
			consummerGroup: consummerGroup,
			topic:          top,
			min_lag:        partLagMin,
			max_lag:        partLagMax,
		}
	}
	return kafkaLag, kafkaLagDetails, nil
}
func (k *Kafka) fetchGroupOffset(grp string, top string, part int32) (groupOffset, error) {
	var (
		err           error
		offsetManager sarama.OffsetManager
	)
	// if cmd.verbose {
	// 	fmt.Fprintf(os.Stderr, "fetching offset information for group=%v topic=%v partition=%v\n", grp, top, part)
	// }
	if offsetManager, err = sarama.NewOffsetManagerFromClient(grp, k.client); err != nil {
		return groupOffset{}, fmt.Errorf("failed to create client err=%v", err)
	}
	defer logClose("offset manager", offsetManager)
	pom, err := offsetManager.ManagePartition(top, part)
	if err != nil {
		return groupOffset{}, fmt.Errorf("failed to manage partition group=%s topic=%s partition=%d err=%v", grp, top, part, err)
	}
	defer logClose("partition offset manager", pom)
	groupOff, _ := pom.NextOffset()
	// we haven't reset it, and it wasn't set before - lag depends on client's config
	if groupOff == sarama.OffsetNewest || groupOff == sarama.OffsetOldest {
		return groupOffset{Partition: part}, nil
	}
	partOff, err := k.resolveOffset(top, part, sarama.OffsetNewest)
	if err != nil {
		return groupOffset{}, err
	}
	lag := partOff - groupOff
	return groupOffset{Partition: part, Offset: &groupOff, Lag: &lag}, nil
}
func logClose(name string, c io.Closer) {
	if err := c.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to close %#v err=%v", name, err)
	}
}
func (k *Kafka) resolveOffset(top string, part int32, off int64) (int64, error) {
	resolvedOff, err := k.client.GetOffset(top, part, off)
	if err != nil {
		return 0, fmt.Errorf("failed to get offset to reset to for partition=%d err=%v", part, err)
	}
	// if cmd.verbose {
	// 	fmt.Fprintf(os.Stderr, "resolved offset %v for topic=%s partition=%d to %v\n", off, top, part, resolvedOff)
	// }
	return resolvedOff, nil
}
func init() {
	inputs.Add("vcs_kafka", func() telegraf.Input {
		return &Kafka{
			Brokers: []string{"localhost:9092"},
			Version: "v0.10.0.0",
			Detail: true,

		}
	})
}
