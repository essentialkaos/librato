// Package librato provides methods and structs for working with Librato Metrics API
package librato

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2015 Essential Kaos                         //
//      Essential Kaos Open Source License <http://essentialkaos.com/ekol?en>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/essentialkaos/ek/req"
	"github.com/essentialkaos/ek/timeutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// DataSource is interface for diferent type of data source
type DataSource interface {
	Send() []error

	getPeriod() int64
	getLastSendingDate() int64
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Metrics struct
type Metrics struct {
	Period       int64
	MaxQueueSize int

	queue           []*Gauge
	lastSendingDate int64
	initialized     bool
}

// Annotations struct
type Annotations struct {
	stream      string
	queue       []*Annotation
	initialized bool
}

// Gauge struct
type Gauge struct {
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
	MeasureTime int64       `json:"measure_time,omitempty"`
	Source      string      `json:"source,omitempty"`
	Count       interface{} `json:"count,omitempty"`
	Sum         interface{} `json:"sum,omitempty"`
	Min         interface{} `json:"min,omitempty"`
	Max         interface{} `json:"max,omitempty"`
	SumSquares  interface{} `json:"sum_squares,omitempty"`
}

// Annotation struct
type Annotation struct {
	Title     string   `json:"title"`
	Source    string   `json:"source,omitempty"`
	Desc      string   `json:"description,omitempty"`
	Links     []string `json:"links,omitempty"`
	StartTime int64    `json:"start_time,omitempty"`
	EndTime   int64    `json:"end_time,omitempty"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Access credentials
var (
	Mail  = ""
	Token = ""
)

// URL of Librato API endpoint
var ApiEndpoint = "https://metrics-api.librato.com"

// AsyncSending enable async data sending (enabled by default)
var AsyncSending = true

// List of sources
var sources []DataSource

// ////////////////////////////////////////////////////////////////////////////////// //

// NewMetrics create new metrics struct
func NewMetrics(period time.Duration) (*Metrics, error) {
	metrics := &Metrics{
		MaxQueueSize: 60,
		Period:       timeutil.DurationToSeconds(period),

		queue:           make([]*Gauge, 0),
		lastSendingDate: -1,
		initialized:     true,
	}

	err := validateMetrics(metrics)

	if err != nil {
		return nil, err
	}

	if sources == nil && AsyncSending {
		sources = make([]DataSource, 0)
		go sendingLoop()
	}

	sources = append(sources, metrics)

	return metrics, nil
}

// NewAnnotations create new annotations struct
func NewAnnotations(stream string) (*Annotations, error) {
	annotations := &Annotations{
		stream:      stream,
		queue:       make([]*Annotation, 0),
		initialized: true,
	}

	err := validateAnotations(annotations)

	if err != nil {
		return nil, err
	}

	if sources == nil && AsyncSending {
		sources = make([]DataSource, 0)
		go sendingLoop()
	}

	sources = append(sources, annotations)

	return annotations, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Add adds gauge to sending queue
func (m *Metrics) Add(g *Gauge) error {
	var err error

	err = validateMetrics(m)

	if err != nil {
		return err
	}

	err = validateGauge(g)

	if err != nil {
		return err
	}

	if g.MeasureTime == 0 {
		g.MeasureTime = time.Now().Unix()
	}

	m.queue = append(m.queue, g)

	if len(m.queue) >= m.MaxQueueSize {
		m.Send()
	}

	return nil
}

// Send sends metrics data to Librato service
func (m *Metrics) Send() []error {
	if Mail == "" || Token == "" {
		return []error{errors.New("Access credentials is not set")}
	}

	var err error

	err = validateMetrics(m)

	if err != nil {
		return []error{err}
	}

	if len(m.queue) == 0 {
		return []error{}
	}

	m.lastSendingDate = time.Now().Unix()

	return m.sendData()
}

// Add adds new annotation to stream
func (an *Annotations) Add(a *Annotation) error {
	var err error

	err = validateAnotations(an)

	if err != nil {
		return err
	}

	err = validateAnotation(a)

	if err != nil {
		return err
	}

	an.queue = append(an.queue, a)

	return nil
}

// Send sends annotations data to Librato service
func (an *Annotations) Send() []error {
	if Mail == "" || Token == "" {
		return []error{errors.New("Access credentials is not set")}
	}

	var err error

	err = validateAnotations(an)

	if err != nil {
		return []error{err}
	}

	if len(an.queue) == 0 {
		return []error{}
	}

	return an.sendData()
}

// Delete remove annotations stream
func (an *Annotations) Delete() error {
	if Mail == "" || Token == "" {
		return errors.New("Access credentials is not set")
	}

	var err error

	err = validateAnotations(an)

	if err != nil {
		return err
	}

	resp, err := req.Request{
		Method: req.DELETE,
		URL:    ApiEndpoint + "/v1/annotations/" + an.stream,

		BasicAuthUsername: Mail,
		BasicAuthPassword: Token,

		AutoDiscard: true,
	}.Do()

	if err != nil {
		return fmt.Errorf("Error while sending request: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Service return error code %d", resp.StatusCode)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getPeriod return sending period
func (m *Metrics) getPeriod() int64 {
	return m.Period
}

// getLastSendingDate return last sending date
func (m *Metrics) getLastSendingDate() int64 {
	return m.lastSendingDate
}

// sendData send json encoded metrics data to Librato service
func (m *Metrics) sendData() []error {
	curQueue := m.queue
	m.queue = make([]*Gauge, 0)

	data := struct {
		Gauges []*Gauge `json:"gauges"`
	}{curQueue}

	jsonData, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		return []error{err}
	}

	resp, err := req.Request{
		Method: req.POST,
		URL:    ApiEndpoint + "/v1/metrics/",

		BasicAuthUsername: Mail,
		BasicAuthPassword: Token,

		ContentType: "application/json",
		Body:        jsonData,

		AutoDiscard: true,
	}.Do()

	if err != nil {
		return []error{err}
	}

	if resp.StatusCode != 200 {
		return []error{fmt.Errorf("Service return error code %d", resp.StatusCode)}
	}

	return []error{}
}

// getPeriod return sending period
func (an *Annotations) getPeriod() int64 {
	return 0
}

// getLastSendingDate return last sending date
func (an *Annotations) getLastSendingDate() int64 {
	return -1
}

// sendData send json encoded annotations data to Librato service
func (an *Annotations) sendData() []error {
	curQueue := an.queue
	an.queue = make([]*Annotation, 0)

	var errs []error

	for _, a := range curQueue {
		jsonData, err := json.MarshalIndent(a, "", "  ")

		if err != nil {
			errs = append(errs, err)
			continue
		}

		resp, err := req.Request{
			Method: req.POST,
			URL:    ApiEndpoint + "/v1/annotations/" + an.stream,

			BasicAuthUsername: Mail,
			BasicAuthPassword: Token,

			ContentType: "application/json",
			Body:        jsonData,

			AutoDiscard: true,
		}.Do()

		if err != nil {
			errs = append(errs, err)
			continue
		}

		if resp.StatusCode != 200 {
			errs = append(errs, fmt.Errorf("Service return error code %d", resp.StatusCode))
			continue
		}
	}

	return errs
}

// ////////////////////////////////////////////////////////////////////////////////// //

func sendingLoop() {
	for {
		time.Sleep(time.Second)

		if len(sources) == 0 {
			continue
		}

		now := time.Now().Unix()

		for _, source := range sources {
			period := source.getPeriod()
			lastSendTime := source.getLastSendingDate()

			if period == 0 || lastSendTime == -1 {
				source.Send()
				continue
			}

			if period+lastSendTime < now {
				source.Send()
				continue
			}
		}
	}
}

func validateMetrics(m *Metrics) error {
	if !m.initialized {
		return errors.New("Metrics struct is not initialized")
	}

	return nil
}

func validateGauge(g *Gauge) error {
	if g.Name == "" {
		return errors.New("Gauge property Name can't be empty")
	}

	if len(g.Name) > 255 {
		return errors.New("Length of gauge property Name must be 255 or fewer characters")
	}

	switch g.Value.(type) {
	case int, int32, int64, uint, uint32, uint64, float32, float64:
	default:
		return errors.New("Gauge property Value can't be non-numeric")
	}

	switch g.Count.(type) {
	case int, int32, int64, uint, uint32, uint64, float32, float64, nil:
	default:
		return errors.New("Gauge property Count can't be non-numeric")
	}

	switch g.Sum.(type) {
	case int, int32, int64, uint, uint32, uint64, float32, float64, nil:
	default:
		return errors.New("Gauge property Sum can't be non-numeric")
	}

	switch g.Min.(type) {
	case int, int32, int64, uint, uint32, uint64, float32, float64, nil:
	default:
		return errors.New("Gauge property Min can't be non-numeric")
	}

	switch g.Max.(type) {
	case int, int32, int64, uint, uint32, uint64, float32, float64, nil:
	default:
		return errors.New("Gauge property Max can't be non-numeric")
	}

	switch g.SumSquares.(type) {
	case int, int32, int64, uint, uint32, uint64, float32, float64, nil:
	default:
		return errors.New("Gauge property SumSquares can't be non-numeric")
	}

	return nil
}

func validateAnotations(a *Annotations) error {
	if !a.initialized {
		return errors.New("Annotations struct is not initialized")
	}

	if a.stream == "" {
		return errors.New("Annotations must have non-empty property stream")
	}

	return nil
}

func validateAnotation(a *Annotation) error {
	if a.Title == "" {
		return errors.New("Annotation property Title can't be empty")
	}

	return nil
}
