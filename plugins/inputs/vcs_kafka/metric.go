package vcs_kafka

import (
	"fmt"
	"sort"
	"time"
)

type Tag struct {
	Key   string
	Value string
}
type Field struct {
	Key   string
	Value interface{}
}
type metric struct {
	name   string
	tags   []*Tag
	fields []*Field
	tm     time.Time
}

func NewMetric(
	name string,
	tags map[string]string,
	fields map[string]interface{},
	tm time.Time,
) *metric {
	m := &metric{
		name:   name,
		tags:   nil,
		fields: nil,
		tm:     tm,
	}
	if len(tags) > 0 {
		m.tags = make([]*Tag, 0, len(tags))
		for k, v := range tags {
			m.tags = append(m.tags,
				&Tag{Key: k, Value: v})
		}
		sort.Slice(m.tags, func(i, j int) bool { return m.tags[i].Key < m.tags[j].Key })
	}
	m.fields = make([]*Field, 0, len(fields))
	for k, v := range fields {
		v := convertField(v)
		if v == nil {
			continue
		}
		m.AddField(k, v)
	}
	return m
}
func FmtTags(tags []*Tag) string {
	out := ""
	for _, tag := range tags {
		out = out + fmt.Sprintf(",%s=%s", tag.Key, tag.Value)
	}
	return out
}
func FmtField(fields []*Field) string {
	out := ""
	for i, field := range fields {
		if i > 0 {
			out = out + ","
		}
		out = out + fmt.Sprintf("%s=%v", field.Key, field.Value)
	}
	return out
}
func (m *metric) String() string {
	return fmt.Sprintf("%s%s %s %d", m.name, FmtTags(m.tags), FmtField(m.fields), m.tm.UnixNano())
}
func (m *metric) Name() string {
	return m.name
}
func (m *metric) Tags() map[string]string {
	tags := make(map[string]string, len(m.tags))
	for _, tag := range m.tags {
		tags[tag.Key] = tag.Value
	}
	return tags
}
func (m *metric) TagList() []*Tag {
	return m.tags
}
func (m *metric) Fields() map[string]interface{} {
	fields := make(map[string]interface{}, len(m.fields))
	for _, field := range m.fields {
		fields[field.Key] = field.Value
	}
	return fields
}
func (m *metric) FieldList() []*Field {
	return m.fields
}
func (m *metric) Time() time.Time {
	return m.tm
}
func (m *metric) SetName(name string) {
	m.name = name
}
func (m *metric) AddPrefix(prefix string) {
	m.name = prefix + m.name
}
func (m *metric) AddSuffix(suffix string) {
	m.name = m.name + suffix
}
func (m *metric) AddTag(key, value string) {
	for i, tag := range m.tags {
		if key > tag.Key {
			continue
		}
		if key == tag.Key {
			tag.Value = value
			return
		}
		m.tags = append(m.tags, nil)
		copy(m.tags[i+1:], m.tags[i:])
		m.tags[i] = &Tag{Key: key, Value: value}
		return
	}
	m.tags = append(m.tags, &Tag{Key: key, Value: value})
}
func (m *metric) HasTag(key string) bool {
	for _, tag := range m.tags {
		if tag.Key == key {
			return true
		}
	}
	return false
}
func (m *metric) GetTag(key string) (string, bool) {
	for _, tag := range m.tags {
		if tag.Key == key {
			return tag.Value, true
		}
	}
	return "", false
}
func (m *metric) RemoveTag(key string) {
	for i, tag := range m.tags {
		if tag.Key == key {
			copy(m.tags[i:], m.tags[i+1:])
			m.tags[len(m.tags)-1] = nil
			m.tags = m.tags[:len(m.tags)-1]
			return
		}
	}
}
func (m *metric) AddField(key string, value interface{}) {
	for i, field := range m.fields {
		if key == field.Key {
			m.fields[i] = &Field{Key: key, Value: convertField(value)}
		}
	}
	m.fields = append(m.fields, &Field{Key: key, Value: convertField(value)})
}
func (m *metric) HasField(key string) bool {
	for _, field := range m.fields {
		if field.Key == key {
			return true
		}
	}
	return false
}
func (m *metric) GetField(key string) (interface{}, bool) {
	for _, field := range m.fields {
		if field.Key == key {
			return field.Value, true
		}
	}
	return nil, false
}
func (m *metric) RemoveField(key string) {
	for i, field := range m.fields {
		if field.Key == key {
			copy(m.fields[i:], m.fields[i+1:])
			m.fields[len(m.fields)-1] = nil
			m.fields = m.fields[:len(m.fields)-1]
			return
		}
	}
}
func (m *metric) SetTime(t time.Time) {
	m.tm = t
}

// Convert field to a supported type or nil if unconvertible
func convertField(v interface{}) interface{} {
	switch v := v.(type) {
	case float64:
		return v
	case int64:
		return v
	case string:
		return v
	case bool:
		return v
	case int:
		return int64(v)
	case uint:
		return uint64(v)
	case uint64:
		return uint64(v)
	case []byte:
		return string(v)
	case int32:
		return int64(v)
	case int16:
		return int64(v)
	case int8:
		return int64(v)
	case uint32:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint8:
		return uint64(v)
	case float32:
		return float64(v)
	default:
		return nil
	}
}
