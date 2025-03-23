package service

import (
	"errors"
	"testing"

	"github.com/mjmhtjain/knime/src/internal/model"
	"github.com/mjmhtjain/knime/src/internal/obj"
	"github.com/mjmhtjain/knime/src/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMessageService_SaveMessage(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		message     *obj.Message
		setupMock   func(*mocks.IOutboxMessageRepository)
		wantErr     bool
		expectedErr string
	}{
		// Original test cases
		{
			name:    "Save message with valid subject",
			message: obj.NewMessage("Test Subject", "Test Body"),
			setupMock: func(mockRepo *mocks.IOutboxMessageRepository) {
				mockRepo.On("Create", mock.AnythingOfType("*model.OutboxMessageEntity")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "Save message with empty subject",
			message: obj.NewMessage("", "Test Body"),
			setupMock: func(mockRepo *mocks.IOutboxMessageRepository) {
				// No mock expectation needed since validation will fail before repository is called
			},
			wantErr:     true,
			expectedErr: "message subject is empty",
		},
		{
			name:    "Save message with nil body",
			message: obj.NewMessage("Test Subject", nil),
			setupMock: func(mockRepo *mocks.IOutboxMessageRepository) {
				// No mock expectation needed since validation will fail before repository is called
			},
			wantErr:     true,
			expectedErr: "message body is nil",
		},

		// Additional mock-based test cases
		{
			name:    "Repository returns database error",
			message: obj.NewMessage("Test Subject", "Test Body"),
			setupMock: func(mockRepo *mocks.IOutboxMessageRepository) {
				mockRepo.On("Create", mock.AnythingOfType("*model.OutboxMessageEntity")).
					Return(errors.New("database connection failed"))
			},
			wantErr:     true,
			expectedErr: "database connection failed",
		},
		{
			name:    "Validates message entity in repository",
			message: obj.NewMessage("Validation Test", "Test Body"),
			setupMock: func(mockRepo *mocks.IOutboxMessageRepository) {
				mockRepo.On("Create", mock.MatchedBy(func(entity *model.OutboxMessageEntity) bool {
					return entity.Subject == "Validation Test"
				})).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock repository for each test
			mockRepo := mocks.NewIOutboxMessageRepository(t)

			// Set up the mock expectations
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			// Create the service with the mock repository
			service := &MessageService{repo: mockRepo}

			// Call the method being tested
			err := service.SaveMessage(tt.message)

			// Assertions
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Equal(t, tt.expectedErr, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
