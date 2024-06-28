package avs

import "github.com/aerospike/avs-client-go/protos"

func createUserPassCredential(username, password string) *protos.Credentials {
	return &protos.Credentials{
		Username: username,
		Credentials: &protos.Credentials_PasswordCredentials{
			PasswordCredentials: &protos.PasswordCredentials{
				Password: password,
			},
		},
	}

}
