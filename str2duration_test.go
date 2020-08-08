package str2duration

import (
	"testing"
	"time"
)

func TestParseString(t *testing.T) {

	DisableCheck = false

	for i, tt := range []struct {
		dur      string
		expected time.Duration
	}{
		//This times are returned with time.Duration string
		{"1h", time.Duration(time.Hour)},
		{"1m", time.Duration(time.Minute)},
		{"1s", time.Duration(time.Second)},
		{"1ms", time.Duration(time.Millisecond)},
		{"1µs", time.Duration(time.Microsecond)},
		{"1us", time.Duration(time.Microsecond)},
		{"1ns", time.Duration(time.Nanosecond)},
		{"4.000000001s", time.Duration(4*time.Second + time.Nanosecond)},
		{"1h0m4.000000001s", time.Duration(time.Hour + 4*time.Second + time.Nanosecond)},
		{"1h1m0.01s", time.Duration(61*time.Minute + 10*time.Millisecond)},
		{"1h1m0.123456789s", time.Duration(61*time.Minute + 123456789*time.Nanosecond)},
		{"1.00002ms", time.Duration(time.Millisecond + 20*time.Nanosecond)},
		{"1.00000002s", time.Duration(time.Second + 20*time.Nanosecond)},
		{"693ns", time.Duration(693 * time.Nanosecond)},
		{"10s1us693ns", time.Duration(10*time.Second + time.Microsecond + 693*time.Nanosecond)},

		//This times aren't returned with time.Duration string, but are easily readable and can be parsed too!
		{"1ms1ns", time.Duration(time.Millisecond + 1*time.Nanosecond)},
		{"1s20ns", time.Duration(time.Second + 20*time.Nanosecond)},
		{"60h8ms", time.Duration(60*time.Hour + 8*time.Millisecond)},
		{"96h63s", time.Duration(96*time.Hour + 63*time.Second)},

		//And works with days and weeks!
		{"2d3s96ns", time.Duration(48*time.Hour + 3*time.Second + 96*time.Nanosecond)},
		{"1w2d3s96ns", time.Duration(168*time.Hour + 48*time.Hour + 3*time.Second + 96*time.Nanosecond)},
		{"1w2d3s3µs96ns", time.Duration(168*time.Hour + 48*time.Hour + 3*time.Second + 3*time.Microsecond + 96*time.Nanosecond)},

		//And can be case insensitive
		{"2D3S96NS", time.Duration(48*time.Hour + 3*time.Second + 96*time.Nanosecond)},

		//This cases are invalid
		{"2.3D3S96NS", time.Duration(48*time.Hour + 3*time.Second + 96*time.Nanosecond)},
		{"2D3S3.66.SMS", time.Duration(48*time.Hour + 3*time.Second + 96*time.Nanosecond)},
	} {
		durationFromString, err := Str2Duration(tt.dur)
		if err != nil {
			t.Logf("index %d -> in: %s returned: %s\tnot equal to %s", i, tt.dur, err.Error(), tt.expected.String())

		} else if tt.expected != durationFromString {
			t.Errorf("index %d -> in: %s returned: %s\tnot equal to %s", i, tt.dur, durationFromString.String(), tt.expected.String())
		}
	}
}

//TestParseDuration test if string returned by a duration is equal to string returned with the package
func TestParseDuration(t *testing.T) {

	DisableCheck = true

	for i, duration := range []time.Duration{
		time.Duration(time.Hour + time.Minute + time.Second + time.Millisecond + time.Microsecond + time.Nanosecond),
		time.Duration(time.Minute + time.Second + time.Millisecond + time.Microsecond + time.Nanosecond),
		time.Duration(time.Hour + time.Minute + time.Second + time.Millisecond + time.Microsecond + time.Nanosecond),
		time.Duration(time.Minute + time.Second + time.Millisecond + time.Microsecond + time.Nanosecond),
		time.Duration(time.Second + time.Millisecond + time.Microsecond + time.Nanosecond),
		time.Duration(time.Millisecond + time.Microsecond + time.Nanosecond),
		time.Duration(time.Microsecond + time.Nanosecond),
		time.Duration(time.Nanosecond),
		time.Duration(time.Hour + time.Minute + time.Second + time.Millisecond + time.Microsecond),
		time.Duration(time.Minute + time.Second + time.Millisecond + time.Microsecond),
		time.Duration(time.Second + time.Millisecond + time.Microsecond),
		time.Duration(time.Millisecond + time.Microsecond),
		time.Duration(time.Microsecond),
		time.Duration(time.Hour + time.Minute + time.Second + time.Millisecond),
		time.Duration(time.Minute + time.Second + time.Millisecond),
		time.Duration(time.Second + time.Millisecond),
		time.Duration(time.Millisecond),
		time.Duration(time.Hour + time.Minute + time.Second),
		time.Duration(time.Minute + time.Second),
		time.Duration(time.Second),
		time.Duration(time.Hour + time.Minute),
		time.Duration(time.Minute),
		time.Duration(time.Hour),
		time.Duration(time.Millisecond + time.Nanosecond),
		time.Duration(1001001 * time.Microsecond),
		time.Duration(1001 * time.Nanosecond),
		time.Duration(61 * time.Minute),
		time.Duration(61 * time.Second),
		time.Duration(time.Microsecond + 16*time.Nanosecond),
	} {
		durationFromString, _ := Str2Duration(duration.String())
		if duration.String() != durationFromString.String() {
			t.Errorf("index %d -> %s not equal to %s", i, duration.String(), durationFromString.String())
		}
	}
}
