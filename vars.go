package goenvvars

import (
	"os"
	"strconv"
)

type EnvVar struct {
	Key           string
	Value         string
	Found         bool
	optional      bool
	allowFallback func() bool
}

type envVarOpt func(*EnvVar)

type fallbackOpt func(*EnvVar)

var DefaultAllowFallback = defaultAllowFallback

func New(key string, opts ...envVarOpt) *EnvVar {
	ev := new(EnvVar)
	ev.Key = key
	ev.allowFallback = DefaultAllowFallback
	ev.Value, ev.Found = os.LookupEnv(key)

	for _, opt := range opts {
		opt(ev)
	}

	return ev
}

func (ev *EnvVar) OverrideAllowFallback(af func() bool) *EnvVar {
	ev.allowFallback = af
	return ev
}

func (ev *EnvVar) Optional() *EnvVar {
	ev.optional = true
	return ev
}

func parseFallbacks[T any](ev *EnvVar, fallbacks ...T) T {
	canFallback := len(fallbacks) > 0 && ev.allowFallback()

	if ev.optional {
		var result T
		if canFallback {
			result = fallbacks[0]
		}
		return result
	}

	if canFallback {
		return fallbacks[0]
	}
	panic("Missing required environment variable: " + ev.Key)
}

func (ev *EnvVar) String(fallbacks ...string) string {
	if ev.Found {
		return ev.Value
	}
	return parseFallbacks(ev, fallbacks...)
}

func (ev *EnvVar) Bool(fallbacks ...bool) bool {
	if ev.Found {
		value, err := strconv.ParseBool(ev.Value)
		if err != nil {
			panic("Invalid boolean environment variable: " + ev.Value)
		}
		return value
	}
	return parseFallbacks(ev, fallbacks...)
}

func (ev *EnvVar) Int(fallbacks ...int) int {
	if ev.Found {
		value, err := strconv.Atoi(ev.Value)
		if err != nil {
			panic("Invalid integer environment variable: " + ev.Value)
		}
		return value
	}
	return parseFallbacks(ev, fallbacks...)
}

func (ev *EnvVar) Float64(fallbacks ...float64) float64 {
	if ev.Found {
		value, err := strconv.ParseFloat(ev.Value, 64)
		if err != nil {
			panic("Invalid float64 environment variable: " + ev.Value)
		}
		return value
	}
	return parseFallbacks(ev, fallbacks...)
}

// Returns true if the environment variable with the given key is set and non-empty
func Presence(key string) bool {
	val, ok := os.LookupEnv(key)
	return ok && val != ""
}

func defaultAllowFallback() bool {
	return !IsProd()
}
