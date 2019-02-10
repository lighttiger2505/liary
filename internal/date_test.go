package internal

import (
	"reflect"
	"testing"
	"time"
)

func TestGetDateRange(t *testing.T) {
	type args struct {
		start time.Time
		end   time.Time
	}
	tests := []struct {
		name string
		args args
		want []time.Time
	}{
		{
			name: "Nomal",
			args: args{
				start: time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.Local),
				end:   time.Date(2019, time.Month(1), 3, 0, 0, 0, 0, time.Local),
			},
			want: []time.Time{
				time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.Local),
				time.Date(2019, time.Month(1), 2, 0, 0, 0, 0, time.Local),
				time.Date(2019, time.Month(1), 3, 0, 0, 0, 0, time.Local),
			},
		},
		{
			name: "Straddling year",
			args: args{
				start: time.Date(2018, time.Month(12), 31, 0, 0, 0, 0, time.Local),
				end:   time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.Local),
			},
			want: []time.Time{
				time.Date(2018, time.Month(12), 31, 0, 0, 0, 0, time.Local),
				time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.Local),
			},
		},
		{
			name: "Straddling month",
			args: args{
				start: time.Date(2019, time.Month(1), 31, 0, 0, 0, 0, time.Local),
				end:   time.Date(2019, time.Month(2), 1, 0, 0, 0, 0, time.Local),
			},
			want: []time.Time{
				time.Date(2019, time.Month(1), 31, 0, 0, 0, 0, time.Local),
				time.Date(2019, time.Month(2), 1, 0, 0, 0, 0, time.Local),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDateRange(tt.args.start, tt.args.end); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDateRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
