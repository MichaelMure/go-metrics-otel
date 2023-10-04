package metricsotel

import "testing"

func TestName(t *testing.T) {
	for _, tc := range []struct {
		name       string
		meter      string
		instrument string
	}{
		{"foo.has_total", "foo", "has_total"},
		{"foo.bar.has_total", "foo.bar", "has_total"},
		{"has_total", "default", "has_total"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCreator(tc.name, "helptext").(*creator)
			if c.meterName != tc.meter {
				t.Errorf("unexpected meter name: expected \"%s\", got \"%s\"", tc.meter, c.meterName)
			}
			if c.instrumentName != tc.instrument {
				t.Errorf("unexpected instrument name: expected \"%s\", got \"%s\"", tc.instrument, c.instrumentName)
			}
		})
	}
}
