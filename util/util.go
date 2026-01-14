package util

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"slices"
	"strings"
)

type MetricLabels map[string]string

func (l MetricLabels) String() string {
	var labelStrs []string
	for key, value := range l {
		labelStrs = append(labelStrs, fmt.Sprintf("%s=%q", key, value))
	}
	slices.Sort(labelStrs)
	labelSection := strings.Join(labelStrs, ", ")
	return labelSection
}

type Metric struct {
	name   string
	labels MetricLabels
	value  string
}

func NewMetric(name string, labels map[string]string, value string) Metric {
	return Metric{
		name:   name,
		labels: labels,
		value:  value,
	}
}

func NewMetricFromObject(obj any) (Metric, error) {
	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return Metric{}, fmt.Errorf("invalid type %s", reflect.TypeOf(obj))
	}

	tags := make(map[string]string)
	reflectToMap(tags, rv)
	metricName, okName := tags["metric_name"]
	if !okName {
		return Metric{}, fmt.Errorf("missing metric_name tag in struct")
	}
	delete(tags, "metric_name")

	metricValue, okValue := tags["metric_value"]
	if !okValue {
		return Metric{}, fmt.Errorf("missing metric_value tag in struct")
	}
	delete(tags, "metric_value")

	return NewMetric(metricName, tags, metricValue), nil
}

func (m Metric) Name() string {
	return formatString(m.name)
}

func (m Metric) Labels() MetricLabels {
	return m.labels
}

func (m Metric) Value() string {
	return m.value
}

func (m Metric) IsValid() bool {
	return m.name != "" && m.value != "" && m.labels != nil
}

func (m Metric) String() string {
	if !m.IsValid() {
		return ""
	}

	result := fmt.Sprintf("%s{%s} %s", m.Name(), m.Labels().String(), m.Value())
	return result
}

// Print writes the Prometheus-formatted representation of obj to w.
// obj must be a pointer.
func Print(obj any, w io.Writer) error {
	metrics, err := NewMetricFromObject(obj)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s\n", metrics.String())
	return err
}

// collectTagsRec recursively collects struct field tags into tagsMapping.
func reflectToMap(tagsMapping map[string]string, val reflect.Value) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		switch f.Kind() {
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
			reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32,
			reflect.Float64, reflect.String:
			label := val.Type().Field(i).Tag.Get("label")
			if label == "" {
				continue
			}
			value := fmt.Sprintf("%v", f.Interface())
			tagsMapping[formatString(label)] = value

		case reflect.Map:
			for _, key := range f.MapKeys() {
				k := formatString(key.String())
				v := f.MapIndex(key).String()
				tagsMapping[k] = v
			}

		case reflect.Slice:
			for j := 0; j < f.Len(); j++ {
				reflectToMap(tagsMapping, f.Index(j))
			}

		case reflect.Struct:
			reflectToMap(tagsMapping, f)
		}
	}
}

func formatString(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		var err error
		if isAlphaNumeric(r) {
			_, err = b.WriteRune(r)

		} else {
			_, err = b.WriteRune('_')
		}

		if err != nil {
			log.Printf("error formatting string: %v", err)
		}
	}
	return b.String()
}

func isAlphaNumeric(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
}
