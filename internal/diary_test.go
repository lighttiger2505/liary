package internal

import (
	"testing"
	"time"
)

func TestDiaryPath(t *testing.T) {
	type args struct {
		targetTime time.Time
		dirPath    string
		suffix     string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "no suffix",
			args: args{
				targetTime: time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC),
				dirPath:    "/tmp/diary",
				suffix:     "",
			},
			want:    "/tmp/diary/2000/02/01.md",
			wantErr: false,
		},
		{
			name: "with suffix",
			args: args{
				targetTime: time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC),
				dirPath:    "/tmp/diary",
				suffix:     "hoge",
			},
			want:    "/tmp/diary/2000/02/01-hoge.md",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DiaryPath(tt.args.targetTime, tt.args.dirPath, tt.args.suffix)
			if (err != nil) != tt.wantErr {
				t.Errorf("DiaryPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DiaryPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
