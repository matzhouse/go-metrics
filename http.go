package metrics

import (
	"net/http"
	"runtime"
)

type sysmetrics struct {
	goroutines int
	version    string
}

type httpmetrics struct {
	metrics []httpmetric
	system  sysmetrics
}

type httpmetric struct {
	Type  string
	Count int
	Value string
}

func Http() {

	port := "127.0.0.1:6007"

	fmt.Println("Starting metrics server on port ", port)

	http.HandleFunc("/", metricshandler)
	http.ListenAndServe(port, nil)

}

func metricshandler(w http.ResponseWriter, r *http.Request) {

	m := data()

	o, err := json.Marshal(m)

	if err != nil {
		fmt.Fprintf(w, "Error in Marshalling JSON")
	} else {
		fmt.Fprintf(w, "%s", o)
	}

}

func data() (output httpmetrics) {

	series := *[]httpmetric

	r.Each(func(name string, i interface{}) {
		now := getCurrentTime()
		switch metric := i.(type) {
		case metrics.Counter:

			series = append(series, &metric{"counter", metric.Count()})

		case metrics.Gauge:
			series = append(series, &influxClient.Series{
				Name:    fmt.Sprintf("%s.value", name),
				Columns: []string{"time", "value"},
				Points: [][]interface{}{
					{now, metric.Value()},
				},
			})
		case metrics.GaugeFloat64:
			series = append(series, &influxClient.Series{
				Name:    fmt.Sprintf("%s.value", name),
				Columns: []string{"time", "value"},
				Points: [][]interface{}{
					{now, metric.Value()},
				},
			})
		case metrics.Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			series = append(series, &influxClient.Series{
				Name: fmt.Sprintf("%s.histogram", name),
				Columns: []string{"time", "count", "min", "max", "mean", "std-dev",
					"50-percentile", "75-percentile", "95-percentile",
					"99-percentile", "999-percentile"},
				Points: [][]interface{}{
					{now, h.Count(), h.Min(), h.Max(), h.Mean(), h.StdDev(),
						ps[0], ps[1], ps[2], ps[3], ps[4]},
				},
			})
		case metrics.Meter:
			m := metric.Snapshot()
			series = append(series, &influxClient.Series{
				Name: fmt.Sprintf("%s.meter", name),
				Columns: []string{"count", "one-minute",
					"five-minute", "fifteen-minute", "mean"},
				Points: [][]interface{}{
					{m.Count(), m.Rate1(), m.Rate5(), m.Rate15(), m.RateMean()},
				},
			})
		case metrics.Timer:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			series = append(series, &influxClient.Series{
				Name: fmt.Sprintf("%s.timer", name),
				Columns: []string{"count", "min", "max", "mean", "std-dev",
					"50-percentile", "75-percentile", "95-percentile",
					"99-percentile", "999-percentile", "one-minute", "five-minute", "fifteen-minute", "mean-rate"},
				Points: [][]interface{}{
					{h.Count(), h.Min(), h.Max(), h.Mean(), h.StdDev(),
						ps[0], ps[1], ps[2], ps[3], ps[4],
						h.Rate1(), h.Rate5(), h.Rate15(), h.RateMean()},
				},
			})
		}
		if err := client.WriteSeries(series); err != nil {
			log.Println(err)
		}
	})

	output.metrics = series

	output.system = systemdata()

	return

}

func metricdata() (md []httpmetric) {

}

func systemdata() (sd sysmetrics) {
	sd = sysmetrics{
		runtime.NumGoroutine(),
		runtime.Version(),
	}
}
