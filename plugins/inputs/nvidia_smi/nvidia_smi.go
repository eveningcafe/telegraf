package nvidia_smi

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/internal"
	"github.com/influxdata/telegraf/plugins/inputs"
)

var (
	measurement = "nvidia_smi"
	metrics     = "fan.speed,memory.total,memory.used,memory.free,pstate,temperature.gpu,name,uuid,compute_mode,utilization.gpu,utilization.memory,index,power.draw,pcie.link.gen.current,pcie.link.width.current,encoder.stats.sessionCount,encoder.stats.averageFps,encoder.stats.averageLatency,clocks.current.graphics,clocks.current.sm,clocks.current.memory,clocks.current.video"
	metricNames = [][]string{
		{"fan_speed", "integer"},
		{"memory_total", "integer"},
		{"memory_used", "integer"},
		{"memory_free", "integer"},
		{"pstate", "tag"},
		{"temperature_gpu", "integer"},
		{"name", "tag"},
		{"uuid", "tag"},
		{"compute_mode", "tag"},
		{"utilization_gpu", "integer"},
		{"utilization_memory", "integer"},
		{"index", "tag"},
		{"power_draw", "float"},
		{"pcie_link_gen_current", "integer"},
		{"pcie_link_width_current", "integer"},
		{"encoder_stats_session_count", "integer"},
		{"encoder_stats_average_fps", "integer"},
		{"encoder_stats_average_latency", "integer"},
		{"clocks_current_graphics", "integer"},
		{"clocks_current_sm", "integer"},
		{"clocks_current_memory", "integer"},
		{"clocks_current_video", "integer"},
	}
)

// NvidiaSMI holds the methods for this plugin
type NvidiaSMI struct {
	BinPath string
	Timeout internal.Duration

	metrics string
}

// Description returns the description of the NvidiaSMI plugin
func (smi *NvidiaSMI) Description() string {
	return "Pulls statistics from nvidia GPUs attached to the host"
}

// SampleConfig returns the sample configuration for the NvidiaSMI plugin
func (smi *NvidiaSMI) SampleConfig() string {
	return `
  ## Optional: path to nvidia-smi binary, defaults to $PATH via exec.LookPath
  # bin_path = "/usr/bin/nvidia-smi"

  ## Optional: timeout for GPU polling
  # timeout = "5s"
`
}

// Gather implements the telegraf interface
func (smi *NvidiaSMI) Gather(acc telegraf.Accumulator) error {

	if _, err := os.Stat(smi.BinPath); os.IsNotExist(err) {
		return fmt.Errorf("nvidia-smi binary not at path %s, cannot gather GPU data", smi.BinPath)
	}

	data, err := smi.pollSMI()
	if err != nil {
		return err
	}

	err = gatherNvidiaSMI(data, acc)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	inputs.Add("nvidia_smi", func() telegraf.Input {
		return &NvidiaSMI{
			BinPath: "/usr/bin/nvidia-smi",
			Timeout: internal.Duration{Duration: 5 * time.Second},
			metrics: metrics,
		}
	})
}

func (smi *NvidiaSMI) pollSMI() (string, error) {
	// Construct and execute metrics query
	opts := []string{"--format=noheader,nounits,csv", fmt.Sprintf("--query-gpu=%s", smi.metrics)}
	ret, err := internal.CombinedOutputTimeout(exec.Command(smi.BinPath, opts...), smi.Timeout.Duration)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

func gatherNvidiaSMI(ret string, acc telegraf.Accumulator) error {
	// First split the lines up and handle each one
	scanner := bufio.NewScanner(strings.NewReader(ret))
	for scanner.Scan() {
		tags, fields, err := parseLine(scanner.Text())
		if err != nil {
			return err
		}
		acc.AddFields(measurement, fields, tags)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error scanning text %s", ret)
	}

	return nil
}

func parseLine(line string) (map[string]string, map[string]interface{}, error) {
	tags := make(map[string]string, 0)
	fields := make(map[string]interface{}, 0)

	// Next split up the comma delimited metrics
	met := strings.Split(line, ",")

	// Make sure there are as many metrics in the line as there were queried.
	if len(met) == len(metricNames) {
		for i, m := range metricNames {
			col := strings.TrimSpace(met[i])

			// Handle the tags
			if m[1] == "tag" {
				tags[m[0]] = col
				continue
			}

			// In some cases we may not be able to get data.
			// One such case is when the memory is overclocked.
			// nvidia-smi reads the max supported memory clock from the stock value.
			// If the current memory clock is greater than the max detected memory clock then we receive [Unknown Error] as a value.

			// For example, the stock max memory clock speed on a 2080 Ti is 7000 MHz which nvidia-smi detects.
			// The user has overclocked their memory using an offset of +1000 so under load the memory clock reaches 8000 MHz.
			// Now when nvidia-smi tries to read the current memory clock it fails and spits back [Unknown Error] as the value.
			// This value will break the parsing logic below unless it is accounted for here.
			if strings.Contains(col, "[Not Supported]") || strings.Contains(col, "[Unknown Error]") {
				continue
			}

			// Parse the integers
			if m[1] == "integer" {
				out, err := strconv.ParseInt(col, 10, 64)
				if err != nil {
					return tags, fields, err
				}
				fields[m[0]] = out
			}

			// Parse the floats
			if m[1] == "float" {
				out, err := strconv.ParseFloat(col, 64)
				if err != nil {
					return tags, fields, err
				}
				fields[m[0]] = out
			}
		}

		// Return the tags and fields
		return tags, fields, nil
	}

	// If the line is empty return an emptyline error
	return tags, fields, fmt.Errorf("Different number of metrics returned (%d) than expeced (%d)", len(met), len(metricNames))
}
