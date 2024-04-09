package date

import (
	"reflect"
	"testing"
	"time"
)

func TestStartOfWeek(t *testing.T) {
	type args struct {
		now time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "should success",
			args: args{
				now: time.Date(2022, 3, 27, 23, 0, 0, 0, Location()),
			},
			want: time.Date(2022, 3, 21, 0, 0, 0, 0, Location()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StartOfWeek(tt.args.now); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StartOfWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLastWeek(t *testing.T) {
	type args struct {
		now time.Time
	}
	tests := []struct {
		name      string
		args      args
		wantStart time.Time
		wantEnd   time.Time
	}{
		{
			name: "should success 2022/03/27",
			args: args{
				now: time.Date(2022, 3, 27, 23, 59, 59, int(time.Second-1), Location()),
			},
			wantStart: time.Date(2022, 3, 14, 0, 0, 0, 0, Location()),
			wantEnd:   time.Date(2022, 3, 20, 23, 59, 59, int(time.Second-1), Location()),
		},
		{
			name: "should success 2020/01/03",
			args: args{
				now: time.Date(2022, 01, 3, 0, 0, 0, 0, Location()),
			},
			wantStart: time.Date(2021, 12, 27, 0, 0, 0, 0, Location()),
			wantEnd:   time.Date(2022, 1, 2, 23, 59, 59, int(time.Second-1), Location()),
		},
		{
			name: "should success 2022/03/28",
			args: args{
				now: time.Date(2022, 3, 28, 6, 0, 0, 0, Location()),
			},
			wantStart: time.Date(2022, 3, 21, 0, 0, 0, 0, Location()),
			wantEnd:   time.Date(2022, 3, 27, 23, 59, 59, int(time.Second-1), Location()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd := LastWeek(tt.args.now)
			if !reflect.DeepEqual(gotStart, tt.wantStart) {
				t.Errorf("LastWeek() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
			if !reflect.DeepEqual(gotStart.Weekday(), time.Monday) {
				t.Errorf("LastWeek() gotStart.Weekday() = %v, want %v", gotStart.Weekday(), time.Monday)
			}
			if !reflect.DeepEqual(gotEnd, tt.wantEnd) {
				t.Errorf("LastWeek() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
			if !reflect.DeepEqual(gotEnd.Weekday(), time.Sunday) {
				t.Errorf("LastWeek() gotEnd.Weekday() = %v, want %v", gotEnd.Weekday(), time.Sunday)
			}
		})
	}
}

func TestLastMonth(t *testing.T) {
	type args struct {
		now time.Time
	}
	tests := []struct {
		name      string
		args      args
		wantStart time.Time
		wantEnd   time.Time
	}{
		{
			name: "should success 2022/04/01",
			args: args{
				now: time.Date(2022, 4, 1, 0, 0, 0, 0, Location()),
			},
			wantStart: time.Date(2022, 3, 1, 0, 0, 0, 0, Location()),
			wantEnd:   time.Date(2022, 3, 31, 23, 59, 59, int(time.Second-1), Location()),
		},
		{
			name: "should success 2022/03/01",
			args: args{
				now: time.Date(2022, 3, 1, 0, 0, 0, 0, Location()),
			},
			wantStart: time.Date(2022, 2, 1, 0, 0, 0, 0, Location()),
			wantEnd:   time.Date(2022, 2, 28, 23, 59, 59, int(time.Second-1), Location()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd := LastMonth(tt.args.now)
			if !reflect.DeepEqual(gotStart, tt.wantStart) {
				t.Errorf("LastMonth() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
			if !reflect.DeepEqual(gotEnd, tt.wantEnd) {
				t.Errorf("LastMonth() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		layout string
		value  string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "should success",
			args: args{
				layout: "2006-01-02",
				value:  "2024-01-01",
			},
			want:    time.Date(2024, 1, 1, 0, 0, 0, 0, Location()),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.layout, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
