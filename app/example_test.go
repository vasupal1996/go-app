package app

import (
	"go-app/server/config"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func getTestConfig() *config.Config {
	return config.GetConfigFromFile("test")
}

func TestExampleImpl_SayHello(t *testing.T) {
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)
	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		wantErr bool
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Ok",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ExampleConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				wantErr: false,
			},
			want:    "Hello",
			wantErr: false,
		},
		{
			name: "Expecting Err",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ExampleConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				wantErr: true,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ExampleImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			got, err := e.SayHello(tt.args.wantErr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExampleImpl.SayHello() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExampleImpl.SayHello() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExampleImpl_SaveHello(t *testing.T) {
	app := NewTestApp(getTestConfig())
	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *SaveHelloOpts
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *SaveHelloResp
	}{
		{
			name: "Ok",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.ExampleConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &SaveHelloOpts{
					Name: "Namaste",
				},
			},
			want: &SaveHelloResp{
				Name: "Namaste",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ExampleImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			got, err := e.SaveHello(tt.args.opts)
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.False(t, got.ID.IsZero())
		})
	}
}
