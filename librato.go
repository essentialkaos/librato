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
	"strings"
	"time"

	"pkg.re/essentialkaos/ek.v5/req"
	"pkg.re/essentialkaos/ek.v5/timeutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// VERSION contains current version of librato package and used as part of User-Agent
const VERSION = "3.0.1"

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
	execErrorHandler(errs []error)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Metrics struct
type Metrics struct {
	period          time.Duration
	maxQueueSize    int
	lastSendingDate int64
	initialized     bool
	queue           []Measurement
	engine          *req.Engine

	// Function executed if we have errors while sending data to Librato
	ErrorHandler func(errs []error)
}

// Collector struct
type Collector struct {
	period          time.Duration
	lastSendingDate int64
	collectFunc     func() []Measurement
	engine          *req.Engine

	// Function executed if we have errors while sending data to Librato
	ErrorHandler func(errs []error)
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

// ////////////////////////////////////////////////////////////////////////////////// //

type measurements struct {
	Gauges   []Gauge   `json:"gauges,omitempty"`
	Counters []Counter `json:"counters,omitempty"`
}

type paramsErrorMap struct {
	Params map[string][]string `json:"params"`
}

type requestErrorSlice struct {
	Request []string `json:"request"`
}

type systemErrorSlice struct {
	System []string `json:"system"`
}

type paramsErrors struct {
	Errors paramsErrorMap `json:"errors"`
}

type requestErrors struct {
	Errors requestErrorSlice `json:"errors"`
}

type systemErrors struct {
	Errors systemErrorSlice `json:"errors"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Access credentials
var (
	Mail  = ""
	Token = ""
)

// APIEndpoint contians URL of Librato API endpoint
var APIEndpoint = "https://metrics-api.librato.com"

// ////////////////////////////////////////////////////////////////////////////////// //

// List of sources
var sources []DataSource

// Some default errors
var (
	errAccessCredentials = []error{errors.New("Access credentials is not set")}
	errEmptyStreamName   = []error{errors.New("Stream name can't be empty")}
)

// Request engine
var engine = &req.Engine{}

// ////////////////////////////////////////////////////////////////////////////////// //

// NewMetrics create new metrics struct for async metrics sending
func NewMetrics(period time.Duration, maxQueueSize int) (*Metrics, error) {
	metrics := &Metrics{
		maxQueueSize:    maxQueueSize,
		period:          period,
		initialized:     true,
		queue:           make([]Measurement, 0),
		lastSendingDate: -1,
		engine:          &req.Engine{},
	}

	err := validateMetrics(metrics)

	if err != nil {
		return nil, err
	}

	if sources == nil {
		sources = make([]DataSource, 0)
		go sendingLoop()
	}

	sources = append(sources, metrics)

	return metrics, nil
}

// NewCollector create new metrics struct for async metrics collecting and sending
func NewCollector(period time.Duration, collectFunc func() []Measurement) *Collector {
	collector := &Collector{
		period:          period,
		collectFunc:     collectFunc,
		lastSendingDate: -1,
		engine:          &req.Engine{},
	}

	if sources == nil {
		sources = make([]DataSource, 0)
		go sendingLoop()
	}

	sources = append(sources, collector)

	return collector
}

// AddMetric synchronously send metric to librato
func AddMetric(m ...Measurement) []error {
	data := measurements{}

	var errs []error

	for _, metric := range m {
		err := metric.Validate()

		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errs
	}

	for _, metric := range m {
		switch metric.(type) {
		case Gauge:
			data.Gauges = append(data.Gauges, metric.(Gauge))

		case Counter:
			data.Counters = append(data.Counters, metric.(Counter))
		}
	}

	return execRequest(engine, req.POST, APIEndpoint+"/v1/metrics/", data)
}

// AddAnnotation synchronously send annotation to librato
func AddAnnotation(stream string, a Annotation) []error {
	if stream == "" {
		return errAccessCredentials
	}

	err := validateAnotation(a)

	if err != nil {
		return []error{err}
	}

	return execRequest(engine, req.POST, APIEndpoint+"/v1/annotations/"+stream, a)
}

// DeleteAnnotations synchronously remove annotation stream on librato
func DeleteAnnotations(stream string) []error {
	if stream == "" {
		return errEmptyStreamName
	}

	return execRequest(engine, req.DELETE, APIEndpoint+"/v1/annotations/"+stream, nil)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Add adds gauge to sending queue
func (mt *Metrics) Add(m ...Measurement) error {
	var err error

	err = validateMetrics(mt)

	if err != nil {
		return err
	}

	for _, metric := range m {
		err = metric.Validate()

		if err != nil {
			return err
		}
	}

	mt.queue = append(mt.queue, m...)

	if len(mt.queue) >= mt.maxQueueSize {
		mt.Send()
	}

	return nil
}

// Send sends metrics data to Librato service
func (mt *Metrics) Send() []error {
	if Mail == "" || Token == "" {
		return errAccessCredentials
	}

	err := validateMetrics(mt)

	if err != nil {
		return []error{err}
	}

	if len(mt.queue) == 0 {
		return nil
	}

	mt.lastSendingDate = time.Now().Unix()

	data := convertMeasurementSlice(mt.queue)

	mt.queue = make([]Measurement, 0)

	errs := execRequest(mt.engine, req.POST, APIEndpoint+"/v1/metrics/", data)

	mt.execErrorHandler(errs)

	return errs
}

// Send sends metrics data to Librato service
func (cl *Collector) Send() []error {
	if Mail == "" || Token == "" {
		return errAccessCredentials
	}

	measurements := cl.collectFunc()

	if len(measurements) == 0 {
		return nil
	}

	var errs []error

	for _, m := range measurements {
		err := m.Validate()

		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		cl.execErrorHandler(errs)
		return errs
	}

	cl.lastSendingDate = time.Now().Unix()

	data := convertMeasurementSlice(measurements)
	errs = execRequest(cl.engine, req.POST, APIEndpoint+"/v1/metrics/", data)

	cl.execErrorHandler(errs)

	return errs
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Validate validates gauge struct
func (g Gauge) Validate() error {
	return validateGauge(g)
}

// Validate validates gauge struct
func (c Counter) Validate() error {
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

// execErrorHandler exec error handler if present
func (mt *Metrics) execErrorHandler(errs []error) {
	if mt.ErrorHandler == nil || len(errs) == 0 {
		return
	}

	mt.ErrorHandler(errs)
}

// getPeriod return sending period
func (cl *Collector) getPeriod() time.Duration {
	return cl.period
}

// getLastSendingDate return last sending date
func (cl *Collector) getLastSendingDate() int64 {
	return cl.lastSendingDate
}

// execErrorHandler exec error handler if present
func (cl *Collector) execErrorHandler(errs []error) {
	if cl.ErrorHandler == nil || len(errs) == 0 {
		return
	}

	cl.ErrorHandler(errs)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// sendingLoop is function used for async sending data from data sources
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

			if period == 0 || lastSendTime <= 0 {
				go source.Send()
				continue
			}

			if period+lastSendTime <= now {
				go source.Send()
			}
		}
	}
}

// convertMeasurementSlice convert slice with measurements to struct
// with counters and gauges slices
func convertMeasurementSlice(data []Measurement) measurements {
	result := measurements{}

	for _, m := range data {
		switch m.(type) {
		case Gauge:
			if result.Gauges == nil {
				result.Gauges = make([]Gauge, 0)
			}

			result.Gauges = append(result.Gauges, m.(Gauge))

		case Counter:
			if result.Counters == nil {
				result.Counters = make([]Counter, 0)
			}

			result.Counters = append(result.Counters, m.(Counter))
		}
	}

	return result
}

// execRequest create and execute request to API
func execRequest(engine *req.Engine, method, url string, data interface{}) []error {
	if engine.UserAgent == "" {
		engine.SetUserAgent("go-ek-librato", VERSION)
	}

	request := req.Request{
		Method: method,
		URL:    url,

		BasicAuthUsername: Mail,
		BasicAuthPassword: Token,

		ContentType: req.CONTENT_TYPE_JSON,

		Close: true,
	}

	if data != nil {
		request.Body = data
	}

	resp, err := engine.Do(request)

	if err != nil {
		return []error{err}
	}

	if resp.StatusCode > 299 || resp.StatusCode == 0 {
		return extractErrors(resp.String())
	}

	resp.Discard()

	return nil
}

// validateMetrics validate metrics struct
func validateMetrics(m *Metrics) error {
	if !m.initialized {
		return errors.New("Metrics struct is not initialized")
	}

	return nil
}

// validateCounter validate counter struct
func validateCounter(c Counter) error {
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

// validateGauge validate gauge struct
func validateGauge(g Gauge) error {
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

// validateAnotation validate annotation struct
func validateAnotation(a Annotation) error {
	if a.Title == "" {
		return errors.New("Annotation property Title can't be empty")
	}

	return nil
}

// extractErrors extracts error descriptions from API response
func extractErrors(data string) []error {
	var err error

	// Data doesn't looks like JSON. Return raw data
	if !strings.HasPrefix(data, "{") {
		return []error{fmt.Errorf(data)}
	}

	switch {
	case strings.Contains(data, "\"params\":"):
		errStruct := &paramsErrors{}
		err = parseErrorsData(data, errStruct)

		if err != nil {
			return []error{fmt.Errorf("Can't parse errors data: %v", err)}
		}

		return mapToErrors(errStruct.Errors.Params)

	case strings.Contains(data, "\"request\":"):
		errStruct := &requestErrors{}
		err = parseErrorsData(data, errStruct)

		if err != nil {
			return []error{fmt.Errorf("Can't parse errors data: %v", err)}
		}

		return sliceToErrors(errStruct.Errors.Request)

	case strings.Contains(data, "\"system\":"):
		errStruct := &systemErrors{}
		err = parseErrorsData(data, errStruct)

		if err != nil {
			return []error{fmt.Errorf("Can't parse errors data: %v", err)}
		}

		return sliceToErrors(errStruct.Errors.System)

	default:
		return []error{fmt.Errorf("Unsupported errors data")}
	}
}

// sliceToErrors convert slice with strings to slice with errors
func sliceToErrors(data []string) []error {
	var result []error

	for _, err := range data {
		result = append(result, errors.New(err))
	}

	return result
}

// mapToErrors convert map with prop name and description to slice with errors
func mapToErrors(data map[string][]string) []error {
	var result []error

	for _, errSlice := range data {
		for _, err := range errSlice {
			result = append(result, errors.New(err))
		}
	}

	return result
}

// parseErrorsData parse error json data to struct
func parseErrorsData(data string, v interface{}) error {
	err := json.Unmarshal([]byte(data), v)

	if err != nil {
		return err
	}

	return nil
}
