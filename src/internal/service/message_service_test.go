package service

import (
	"testing"

	"github.com/mjmhtjain/knime/src/config"
	"github.com/mjmhtjain/knime/src/internal/obj"
)

func createTestDBConfig() *config.OutboxDBConfig {
	return config.NewOutboxDBConfig(
		"localhost",
		"5432",
		"test_user",
		"test_password",
		"test_db",
	)
}

func TestMessageService_SaveMessage(t *testing.T) {
	type args struct {
		msg *obj.Message
	}
	tests := []struct {
		name    string
		s       *MessageService
		args    args
		wantErr bool
	}{
		{
			name: "Save message with valid subject",
			s:    NewMessageService(createTestDBConfig()),
			args: args{
				msg: obj.NewMessage("Test Subject", "Test Body"),
			},
			wantErr: false,
		},
		{
			name: "Save message with empty subject",
			s:    NewMessageService(createTestDBConfig()),
			args: args{
				msg: obj.NewMessage("", "Test Body"),
			},
			wantErr: true,
		},
		{
			name: "Save message with nil body",
			s:    NewMessageService(createTestDBConfig()),
			args: args{
				msg: obj.NewMessage("Test Subject", nil),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.s
			if err := s.SaveMessage(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("MessageService.SaveMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
