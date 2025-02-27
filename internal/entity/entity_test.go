package entity

import (
	"testing"
	"time"
)

func TestClientOrder_SetContactType(t *testing.T) {
	type fields struct {
		ID          int
		Contact     string
		ContactType ContactType
		Message     string
		CreatedAt   time.Time
		IsSent      bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid email",
			fields: fields{
				Contact: "exempele@mail.ru",
			},
			wantErr: false,
		},
		{
			name: "valid contact",
			fields: fields{
				Contact: "89081066015",
			},
			wantErr: false,
		},
		{
			name: "valid telegram",
			fields: fields{
				Contact: "@telegram",
			},
			wantErr: false,
		},
		{
			name: "not valid email",
			fields: fields{
				Contact: "exempele-mail.ru",
			},
			wantErr: true,
		},
		{
			name: "not valid contact",
			fields: fields{
				Contact: "123",
			},
			wantErr: true,
		},
		{
			name: "not valid telegram",
			fields: fields{
				Contact: "telegram",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClientOrder{
				ID:          tt.fields.ID,
				Contact:     tt.fields.Contact,
				ContactType: tt.fields.ContactType,
				Message:     tt.fields.Message,
				CreatedAt:   tt.fields.CreatedAt,
				IsSent:      tt.fields.IsSent,
			}
			if err := c.SetContactType(); (err != nil) != tt.wantErr {
				t.Errorf("ClientOrder.SetContactType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientOrder_Validate(t *testing.T) {
	type fields struct {
		ID          int
		Contact     string
		ContactType ContactType
		Message     string
		CreatedAt   time.Time
		IsSent      bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid struct",
			fields: fields{
				Contact: "exempele@mail.ru",
				Message: "message",
			},
			wantErr: false,
		},
		{
			name: "not valid contact",
			fields: fields{
				Message: "апап",
			},
			wantErr: true,
		},
		{
			name: "not valid message",
			fields: fields{
				Contact: "exempele@mail.ru",
				Message: "fgdgdfgdfgdfgdfgdfgdgdgfLorem ipsum dolor sit amet, consectetur adipiscing elit. Proin id urna tellus. Nulla euismod urna massa, vitae consectetur justo rhoncus mollis. Morbi pulvinar quis nunc at ultrices. Sed ut dolor non tellus mollis mollis eget vel nulla. Praesent condimentum elit sit amet est eleifend, vel maximus massa commodo. Suspendisse potenti. Nunc semper scelerisque consequat. Mauris quis neque lobortis, laoreet purus eu, auctor tortor. Interdum et malesuada fames ac ante ipsum primis in faucibus accums..",
			},
			wantErr: true,
		},
		{
			name: "not valid message empty",
			fields: fields{
				Contact: "exempele@mail.ru",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClientOrder{
				ID:          tt.fields.ID,
				Contact:     tt.fields.Contact,
				ContactType: tt.fields.ContactType,
				Message:     tt.fields.Message,
				CreatedAt:   tt.fields.CreatedAt,
				IsSent:      tt.fields.IsSent,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ClientOrder.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
