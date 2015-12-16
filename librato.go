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

// Measurement is interface for different type of measurements
type Measurement interface {
	Validate() error
}

// DataSource is interface for diferent type of data source
type DataSource interface {
	Send() []error

	getPeriod() time.Duration
	getLastSendingDate() int64
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Metrics struct
type Metrics struct {
	period          time.Duration
	maxQueueSize    int
	lastSendingDate int64
	initialized     bool
	queue           []Measurement
}

// Annotations struct
type Annotations struct {
	stream      string
	queue       []*Annotation
	initialized bool
}

// Gauge struct
type Gauge struct {

	// Each metric has a name that is unique to its class of metrics e.g. a gauge name
	// must be unique amongst gauges. The name identifies a metric in subsequent API
	// calls to store/query individual measurements and can be up to 255 characters
	// in length. Valid characters for metric names are 'A-Za-z0-9.:-_'. The metric
	// namespace is case insensitive.
	Name string `json:"name"`

	// The numeric value of an individual measurement. Multiple formats are
	// supported (e.g. integer, floating point, etc) but the value must be numeric.
	Value interface{} `json:"value"`

	// The epoch time at which an individual measurement occurred with a maximum
	// resolution of seconds.
	MeasureTime int64 `json:"measure_time,omitempty"`

	// Source is an optional property that can be used to subdivide a common
	// gauge/counter amongst multiple members of a population. For example
	// the number of requests/second serviced by an application could be broken
	// up amongst a group of server instances in a scale-out tier by setting
	// the hostname as the value of source.
	//
	// Source names can be up to 255 characters in length and must be composed
	// of the following 'A-Za-z0-9.:-_'. The word all is a reserved word and
	// cannot be used as a user source. The source namespace is case insensitive.
	Source string `json:"source,omitempty"`

	// Indicates the request corresponds to a multi-sample measurement. This is
	// useful if measurements are taken very frequently in a closed loop and the
	// metric value is only periodically reported. If count is set, then sum must
	// also be set in order to calculate an average value for the recorded metric
	// measurement. Additionally min, max, and sum_squares may also be set when
	// count is set. The value parameter should not be set if count is set.
	Count interface{} `json:"count,omitempty"`

	// If count was set, sum must be set to the summation of the individual
	// measurements. The combination of count and sum are used to calculate an
	// average value for the recorded metric measurement.
	Sum interface{} `json:"sum,omitempty"`

	// If count was set, min can be used to report the smallest individual
	// measurement amongst the averaged set.
	Min interface{} `json:"min,omitempty"`

	// If count was set, max can be used to report the largest individual
	// measurement amongst the averaged set.
	Max interface{} `json:"max,omitempty"`

	// If count was set, sum_squares report the summation of the squared
	// individual measurements. If sum_squares is set, a standard deviation
	// can be calculated for the recorded metric measurement.
	SumSquares interface{} `json:"sum_squares,omitempty"`
}

// Counter struct
type Counter struct {
	// Each metric has a name that is unique to its class of metrics e.g. a gauge name
	// must be unique amongst gauges. The name identifies a metric in subsequent API
	// calls to store/query individual measurements and can be up to 255 characters
	// in length. Valid characters for metric names are 'A-Za-z0-9.:-_'. The metric
	// namespace is case insensitive.
	Name string `json:"name"`

	// The numeric value of an individual measurement. Multiple formats are
	// supported (e.g. integer, floating point, etc) but the value must be numeric.
	Value interface{} `json:"value"`

	// The epoch time at which an individual measurement occurred with a maximum
	// resolution of seconds.
	MeasureTime int64 `json:"measure_time,omitempty"`

	// Source is an optional property that can be used to subdivide a common
	// gauge/counter amongst multiple members of a population. For example
	// the number of requests/second serviced by an application could be broken
	// up amongst a group of server instances in a scale-out tier by setting
	// the hostname as the value of source.
	//
	// Source names can be up to 255 characters in length and must be composed
	// of the following 'A-Za-z0-9.:-_'. The word all is a reserved word and
	// cannot be used as a user source. The source namespace is case insensitive.
	Source string `json:"source,omitempty"`
}

// Annotation struct
type Annotation struct {

	// The title of an annotation is a string and may contain spaces. The title should
	// be a short, high-level summary of the annotation e.g. v45 Deployment. The title
	// is a required parameter to create an annotation.
	Title string `json:"title"`

	// A string which describes the originating source of an annotation when that
	// annotation is tracked across multiple members of a population.
	// Examples: foo3.bar.com, user-123, 77025.
	Source string `json:"source,omitempty"`

	// The description contains extra meta-data about a particular annotation. The
	// description should contain specifics on the individual annotation e.g.
	// Deployed 9b562b2: shipped new feature foo! A description is not required to
	// create an annotation.
	Desc string `json:"description,omitempty"`

	// An optional list of references to resources associated with the particular
	// annotation. For example, these links could point to a build page in a CI
	// system or a changeset description of an SCM. Each link has a tag that
	// defines the link\'s relationship to the annotation.
	Links []string `json:"links,omitempty"`

	// The unix timestamp indicating the the time at which the event referenced by this
	// annotation started. By default this is set to the current time if not specified.
	StartTime int64 `json:"start_time,omitempty"`

	// The unix timestamp indicating the the time at which the event referenced by
	// this annotation ended. For events that have a duration, this is a useful way
	// to annotate the duration of the event.
	EndTime int64 `json:"end_time,omitempty"`
}

type mesdata struct {
	Gauges   []*Gauge   `json:"gauges,omitempty"`
	Counters []*Counter `json:"counters,omitempty"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Access credentials
var (
	Mail  = ""
	Token = ""
)

// APIEndpoint contians URL of Librato API endpoint
var APIEndpoint = "https://metrics-api.librato.com"

// AsyncSending enable async data sending (enabled by default)
var AsyncSending = true

// List of sources
var sources []DataSource

// ////////////////////////////////////////////////////////////////////////////////// //

// NewMetrics create new metrics struct
func NewMetrics(period time.Duration, maxQueueSize int) (*Metrics, error) {
	metrics := &Metrics{
		maxQueueSize:    maxQueueSize,
		period:          period,
		lastSendingDate: -1,
		initialized:     true,
		queue:           make([]Measurement, 0),
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
func (mt *Metrics) Add(m Measurement) error {
	var err error

	err = validateMetrics(mt)

	if err != nil {
		return err
	}

	err = m.Validate()

	if err != nil {
		return err
	}

	mt.queue = append(mt.queue, m)

	if len(mt.queue) >= mt.maxQueueSize {
		mt.Send()
	}

	return nil
}

// Send sends metrics data to Librato service
func (mt *Metrics) Send() []error {
	if Mail == "" || Token == "" {
		return []error{errors.New("Access credentials is not set")}
	}

	var err error

	err = validateMetrics(mt)

	if err != nil {
		return []error{err}
	}

	if len(mt.queue) == 0 {
		return []error{}
	}

	mt.lastSendingDate = time.Now().Unix()

	return mt.sendData()
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
		URL: APIEndpoint + "/v1/annotations/" + an.stream,

		BasicAuthUsername: Mail,
		BasicAuthPassword: Token,

		AutoDiscard: true,
	}.Delete()

	if err != nil {
		return fmt.Errorf("Error while sending request: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Service return error code %d", resp.StatusCode)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Validate validates gauge struct
func (g *Gauge) Validate() error {
	return validateGauge(g)
}

// Validate validates gauge struct
func (c *Counter) Validate() error {
	return validateCounter(c)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getPeriod return sending period
func (mt *Metrics) getPeriod() time.Duration {
	return mt.period
}

// getLastSendingDate return last sending date
func (mt *Metrics) getLastSendingDate() int64 {
	return mt.lastSendingDate
}

// sendData send json encoded metrics data to Librato service
func (mt *Metrics) sendData() []error {
	jsonData, err := json.MarshalIndent(convertQueue(mt.queue), "", "  ")

	mt.queue = make([]Measurement, 0)

	if err != nil {
		return []error{err}
	}

	resp, err := req.Request{
		Method: req.POST,
		URL:    APIEndpoint + "/v1/metrics/",

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
func (an *Annotations) getPeriod() time.Duration {
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
			URL: APIEndpoint + "/v1/annotations/" + an.stream,

			BasicAuthUsername: Mail,
			BasicAuthPassword: Token,

			ContentType: "application/json",
			Body:        jsonData,

			AutoDiscard: true,
		}.Post()

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
			period := timeutil.DurationToSeconds(source.getPeriod())
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

func convertQueue(queue []Measurement) *mesdata {
	result := &mesdata{}

	now := time.Now().Unix()

	for _, m := range queue {
		switch m.(type) {
		case *Gauge:
			if result.Gauges == nil {
				result.Gauges = make([]*Gauge, 0)
			}

			gauge := m.(*Gauge)

			if gauge.MeasureTime != 0 {
				gauge.MeasureTime = now
			}

			result.Gauges = append(result.Gauges, gauge)

		case *Counter:
			if result.Counters == nil {
				result.Counters = make([]*Counter, 0)
			}

			counter := m.(*Counter)

			if counter.MeasureTime != 0 {
				counter.MeasureTime = now
			}

			result.Counters = append(result.Counters, counter)
		}
	}

	return result
}

func validateMetrics(m *Metrics) error {
	if !m.initialized {
		return errors.New("Metrics struct is not initialized")
	}

	return nil
}

func validateCounter(c *Counter) error {
	if c.Name == "" {
		return errors.New("Counter property Name can't be empty")
	}

	if len(c.Name) > 255 {
		return errors.New("Length of counter property Name must be 255 or fewer characters")
	}

	switch c.Value.(type) {
	case int, int32, int64, uint, uint32, uint64, float32, float64:
	default:
		return errors.New("Counter property Value can't be non-numeric")
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
