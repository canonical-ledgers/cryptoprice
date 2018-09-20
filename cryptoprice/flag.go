package main

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/posener/complete"
)

var (
	fromSymbol string
	toSymbol   string
	datetime   timeT
	unix       uint64
	debug      bool
	timeout    time.Duration

	flags = complete.Flags{
		"-from":     complete.PredictAnything,
		"-to":       complete.PredictAnything,
		"-datetime": complete.PredictAnything,
		"-timeout":  complete.PredictAnything,
		"-debug":    complete.PredictNothing,

		"-y":                   complete.PredictNothing,
		"-installcompletion":   complete.PredictNothing,
		"-uninstallcompletion": complete.PredictNothing,
	}

	completion *complete.Complete
)

func init() {
	flag.StringVar(&fromSymbol, "from", "BTC", "Currency to convert from")
	flag.StringVar(&toSymbol, "to", "USD", "Currency to convert to")
	flag.BoolVar(&debug, "debug", false, "Enable debug output")
	flag.DurationVar(&timeout, "timeout", 5*time.Second,
		"Timeout for HTTP request to CryptoCompare API")
	flag.Var(&datetime, "datetime", "Time for the price lookup")

	// Add flags for self installing the CLI completion tool
	completion = complete.New(os.Args[0], complete.Command{Flags: flags})
	completion.CLI.InstallName = "installcompletion"
	completion.CLI.UninstallName = "uninstallcompletion"
	completion.AddFlags(nil)
}

type timeT time.Time

func (t timeT) String() string {
	return time.Time(t).String()
}

func (t *timeT) Set(s string) error {
	if ts, err := time.Parse("2006-01-02T15:04:05Z", s); err == nil {
		*t = timeT(ts)
		return nil
	}
	if ts, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
		h, n, s := ts.Clock()
		y, m, d := ts.Date()
		*t = timeT(time.Date(y, m, d, h, n, s, 0, time.Local))
		return nil
	}
	if ts, err := time.Parse("15:04:05", s); err == nil {
		h, n, s := ts.Clock()
		y, m, d := time.Now().Date()
		*t = timeT(time.Date(y, m, d, h, n, s, 0, time.Local))
		return nil
	}
	if ts, err := time.Parse("15:04", s); err == nil {
		h, n, _ := ts.Clock()
		y, m, d := time.Now().Date()
		*t = timeT(time.Date(y, m, d, h, n, 0, 0, time.Local))
		return nil
	}
	if unix, err := strconv.ParseUint(s, 10, 64); err == nil {
		*t = timeT(time.Unix(int64(unix), 0))
		return nil
	}

	return nil
}

func (t timeT) Get() interface{} {
	return t
}
