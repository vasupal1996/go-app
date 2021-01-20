package auth

import (
	"go-app/server/config"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func getTestConfig() *config.TokenAuthConfig {
	c := config.GetConfigFromFile("test")
	return &c.TokenAuthConfig
}

var testTokenAuth = NewTokenAuthentication(getTestConfig())

func getTestUserClaim() *UserClaim {
	uc := UserClaim{
		ID:   uuid.NewV4().String(),
		Type: "user",
	}
	return &uc
}

func TestTokenAuthentication_SignToken(t *testing.T) {
	type fields struct {
		Config *config.TokenAuthConfig
		User   *UserAuth
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Basic Auth",
			fields: fields{
				Config: getTestConfig(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tAuth := NewTokenAuthentication(tt.fields.Config)
			tAuth.SetClaim(getTestUserClaim())
			got, err := tAuth.SignToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenAuthentication.SignToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotEmpty(t, got)
			assert.NotNil(t, tAuth.User.JWTToken)
		})
	}
}

func TestTokenAuthentication_VerifyToken(t *testing.T) {
	uc := getTestUserClaim()

	testTokenAuth.SetClaim(uc)
	tokenString, _ := testTokenAuth.SignToken()

	testConfigInvalidSignature := getTestConfig()
	testConfigInvalidSignature.JWTSignKey = "abccadnced"

	type args struct {
		tokenString string
	}
	type fields struct {
		Config *config.TokenAuthConfig
		User   *UserAuth
	}
	tests := []struct {
		name          string
		args          args
		fields        fields
		wantErr       bool
		wantErrString string
		wantClaim     Claim
	}{
		{
			name:    "Basic Token Verify",
			wantErr: false,
			args: args{
				tokenString: tokenString,
			},
			fields: fields{
				Config: getTestConfig(),
				User:   &UserAuth{},
			},
			wantClaim: uc,
		},
		{
			name:    "Invalid Token String",
			wantErr: true,
			args: args{
				tokenString: "ZXlKaGJHY2lPaUpJ5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SnBaQ0k2SWpObE9XUm1NVGcxTFRaa01EUXROREZqTXkwNVl6VmhMVGMzWW1Vd1pEYzNNbVJpTkNJc0luUjVjR1VpT2lKMWMyVnlJbjAuSTZHajBPYzBaYTBzZUVvX2VzX29EZnJmNDE3V1p2bGJJcF9tOTU1Nk9NTQ==",
			},
			fields: fields{
				Config: getTestConfig(),
				User:   &UserAuth{UserClaim: &UserClaim{}},
			},
			wantErrString: "illegal base64 data at input byte 208",
			wantClaim:     &UserClaim{},
		},
		{
			name:    "Invalid Token Signature",
			wantErr: true,
			args: args{
				tokenString: "ZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SnBaQ0k2SWpNNE1qQmhZMk13TFdWbU0ySXRORFExTnkxaFlqUmhMVE13T0Rsall6VXhOV0U0WXlJc0luUjVjR1VpT2lKMWMyVnlJbjAuVEFoS0Vhenh6UjU3T3VZb0FKZVVfeVVKT3ktU0tzSXRBN0FNZlptRG15SQ==",
			},
			fields: fields{
				Config: testConfigInvalidSignature,
				User:   &UserAuth{UserClaim: &UserClaim{}},
			},
			wantErrString: "signature is invalid",
			wantClaim:     &UserClaim{},
		},
		{
			name:    "Token Expired",
			wantErr: true,
			args: args{
				tokenString: "ZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SnBaQ0k2SW1NeVpETmtNMlkwTFRjelpHRXROR001TmkwNE1URmlMVE01TmpGaU5qUmtPVEZtTUNJc0luUjVjR1VpT2lKMWMyVnlJaXdpWlhod0lqb3hOakV3T1RBNU5UTXhmUS5aeU9sVXVSMkhKUzUteFJjcjdfd000dXdmMU90NkdrdmhiYjdpbi0yR3dV",
			},
			fields: fields{
				Config: getTestConfig(),
				User:   &UserAuth{UserClaim: &UserClaim{}},
			},
			wantErrString: "token is expired by",
			wantClaim:     &UserClaim{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tAuth := &TokenAuthentication{
				Config: tt.fields.Config,
				User:   tt.fields.User,
			}
			if err := tAuth.VerifyToken(tt.args.tokenString); (err != nil) != tt.wantErr {
				t.Errorf("TokenAuthentication.VerifyToken() error = %v, wantErr %v", err, tt.wantErr)
				assert.Contains(t, tt.wantErrString, err.Error())
			}
			assert.Equal(t, tt.wantClaim, tAuth.GetClaim())

		})
	}
}

func TestTokenAuthentication_GetClaimWithTokenString(t *testing.T) {
	uc := getTestUserClaim()
	testTokenAuth.SetClaim(uc)
	tokenString, _ := testTokenAuth.SignToken()
	type fields struct {
		Config *config.TokenAuthConfig
		User   *UserAuth
	}
	type args struct {
		tokenString string
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   Claim
	}{
		{
			name: "Basic Get Claim",
			fields: fields{
				Config: getTestConfig(),
				User:   &UserAuth{},
			},
			want: uc,
			args: args{
				tokenString: tokenString,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tAuth := &TokenAuthentication{
				Config: tt.fields.Config,
				User:   tt.fields.User,
			}
			tAuth.VerifyToken(tt.args.tokenString)
			if got := tAuth.GetClaim(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TokenAuthentication.GetClaim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenAuthentication_GetClaim(t *testing.T) {
	uc := getTestUserClaim()
	type fields struct {
		Config *config.TokenAuthConfig
		User   *UserAuth
	}
	tests := []struct {
		name   string
		fields fields
		want   Claim
	}{
		{
			name: "Basic Get Claim",
			fields: fields{
				Config: getTestConfig(),
				User: &UserAuth{
					UserClaim: uc,
				},
			},
			want: uc,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tAuth := &TokenAuthentication{
				Config: tt.fields.Config,
				User:   tt.fields.User,
			}
			if got := tAuth.GetClaim(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TokenAuthentication.GetClaim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenAuthentication_SetClaim(t *testing.T) {
	uc := getTestUserClaim()
	type fields struct {
		Config *config.TokenAuthConfig
		User   *UserAuth
	}
	type args struct {
		uc Claim
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Claim
	}{
		{
			name: "Basic Set Claim",
			fields: fields{
				Config: getTestConfig(),
				User:   &UserAuth{},
			},
			args: args{
				uc: uc,
			},
			want: uc,
		},
		{
			name: "Basic Set Claim With Nil Claim",
			fields: fields{
				Config: getTestConfig(),
				User:   &UserAuth{UserClaim: &UserClaim{}},
			},
			args: args{
				uc: nil,
			},
			want: &UserClaim{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tAuth := &TokenAuthentication{
				Config: tt.fields.Config,
				User:   tt.fields.User,
			}
			tAuth.SetClaim(tt.args.uc)
			assert.Equal(t, tt.want, tAuth.GetClaim())
		})
	}
}
