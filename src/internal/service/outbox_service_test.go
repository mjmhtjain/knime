package service

import (
	"errors"
	"testing"

	"github.com/mjmhtjain/knime/src/internal/model"
	"github.com/mjmhtjain/knime/src/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestOutboxService_ConsumeOutboxMessages(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		setupMock   func(*mocks.IOutboxMessageRepository)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Success - Messages consumed correctly",
			setupMock: func(mockRepo *mocks.IOutboxMessageRepository) {
				mockMessages := []model.OutboxMessageEntity{
					{
						ID:      "msg1",
						Subject: "test-subject-1",
						Status:  "pending",
					},
					{
						ID:      "msg2",
						Subject: "test-subject-2",
						Status:  "pending",
					},
				}
				mockRepo.On("PushPendingMessages").Return(mockMessages, nil)
			},
			wantErr: false,
		},
		{
			name: "Error - Repository returns error",
			setupMock: func(mockRepo *mocks.IOutboxMessageRepository) {
				expectedErr := errors.New("repository error")
				mockRepo.On("PushPendingMessages").Return(nil, expectedErr)
			},
			wantErr:     true,
			expectedErr: errors.New("repository error"),
		},
		{
			name: "Success - Empty message list",
			setupMock: func(mockRepo *mocks.IOutboxMessageRepository) {
				mockRepo.On("PushPendingMessages").Return([]model.OutboxMessageEntity{}, nil)
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
			svc := &OutboxService{
				outboxRepository: mockRepo,
			}

			// Call the method being tested
			err := svc.ConsumeOutboxMessages()

			// Assertions
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify that all expected calls were made
			mockRepo.AssertExpectations(t)
		})
	}
}
