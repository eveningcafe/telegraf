package vcs_kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"log"
	"os"
	"regexp"
	"sync"
)

var sampleConfig = `
  ## kafka servers
  brokers = ["localhost:9092"]
  ## Kafka protocol version
  version = "v0.8.2.0"
  ## Consumer group name
  group = "correlation"
  ## Regex to filter groups
  filter_groups = ""
  ## Topic to consume
  topic = ""
  ## Print detail for all partitions or not
  detail = false
  ## Comma separated list of partitions to limit offsets to, or all
  partitions = "all"
`

type Kafka struct {
	Topic        string
	Brokers      []string
	Partitions   []int32
	Group        string
	FilterGroups string
	Detail       bool
	Version      string
	Debug        bool

	filter  *regexp.Regexp
	version sarama.KafkaVersion
	conf    *sarama.Config
	client  sarama.Client
}

func (Kafka) SampleConfig() string {
	return sampleConfig
}

func (k *Kafka) Description() string {
	return "Read metrics from Kafka topic(s)"
}

func (k *Kafka) parseArgs() error {
	var err error
	if kafka.filter, err = regexp.Compile(kafka.FilterGroups); err != nil {
		failf("filter regexp invalid err=%v", err)
	}
}

func (k *Kafka) createClient() (sarama.Client, error) {
	var err error
	kafka.parseArgs()
	cfg := sarama.NewConfig()
	cfg.Version = .version
	client, err := sarama.NewClient(kafka.Brokers, kafka.saramaConfig())
	return client, err
}
func (k *Kafka) connect(broker *sarama.Broker) error {
	if ok, _ := broker.Connected(); ok {
		return nil
	}
	cfg := sarama.NewConfig()
	cfg.Version = cmd.version
	cfg.ClientID = "kafka-monitor-tool"

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

func (k *Kafka) findGroupsOnBroker(broker *sarama.Broker, results chan findGroupResult, errs chan error) {
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

func (k *Kafka) findGroups(brokers []*sarama.Broker) []string {
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

func (k *Kafka) Gather(telegraf.Accumulator) error {
	kafka := &Kafka{}
	client, err := kafka.createClient()
	if (err != nil) {
		return err
	}
	brokers := client.Brokers()
	if kafka.Debug {
		var addrs []string
		for _, b := range brokers {
			addrs = append(addrs, b.Addr())
		}
		log.Printf("D! Found brokers: %v\n", addrs)
	}

	groups := []string{k.Group}
	if k.Group == "" {
		groups = []string{}
		for _, g := range k.findGroups(brokers) {
			if cmd.filter.MatchString(g) {
				groups = append(groups, g)
			}
		}
	}
	if kafka.Debug {
		log.Printf("D! Found groups: %v\n", groups)
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

func init() {
	inputs.Add("vcs_kafka", func() telegraf.Input {
		return &Kafka{
			partitions: []int32{},
		}
	})
}
