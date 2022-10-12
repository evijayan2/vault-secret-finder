package internal

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Resource struct {
	client       *api.Client
	vaultAddress string
	VaultToken   string
	Logger       zerolog.Logger
	storage      map[string]interface{}
}

func (r Resource) SetAddress(add string) {
	r.vaultAddress = add
}

func New(vaultAddr string, logger zerolog.Logger) Resource {
	c, err := api.NewClient(
		&api.Config{
			Address: vaultAddr,
			HttpClient: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			},
			MaxRetries: 3,
		},
	)
	if err != nil {
		logger.Fatal().AnErr("err", err).
			Msg("error occured creating client")
	}

	if len(viper.GetString("VAULT_TOKEN")) > 0 {
		c.SetToken(viper.GetString("VAULT_TOKEN"))
	}
	vaultResource := Resource{
		VaultToken: viper.GetString("VAULT_TOKEN"),
		Logger:     logger,
		client:     c,
	}
	return vaultResource
}

// renewToken - renews a token
func (r *Resource) renewToken() error {

	if len(r.VaultToken) <= 0 {
		return errors.New("error renewing vault client token, no vault_token provided")
	}
	r.Logger.Debug().Msg("attempting renewal of token")

	resp, err := r.client.Logical().Write("auth/token/renew-self", nil)

	if err != nil {
		return err
	}

	if resp.Auth != nil {
		r.client.SetToken(resp.Auth.ClientToken)
		r.Logger.Debug().Msg("succesfully renewed token")
		return nil
	}

	return errors.New("error renewing token")
}

func (r *Resource) gatherAll(org string) {
	err := r.renewToken()
	if err != nil {
		r.Logger.Fatal().AnErr("err", err).
			Msg("error occured renewing token")
	}

	kvsecret, err := r.readList(org, "")
	if err != nil {
		r.Logger.Fatal().AnErr("err", err).
			Msg("error occured reading secret")
	}
	storage := make(map[string]interface{}, 0)

	b, _ := json.Marshal(kvsecret.Data)

	var m map[string][]string
	json.Unmarshal(b, &m)

	newFunction("", m, r, org, storage)
	r.storage = storage
}

func newFunction(key string, baseMap map[string][]string, r *Resource, org string, storage map[string]interface{}) {
	for _, v := range baseMap["keys"] {
		if strings.HasSuffix(v, "/") {
			secret, err := r.readList(org, v)
			if err != nil {
				r.Logger.Fatal().AnErr("err", err).
					Msg("error occured reading secret")
			}
			b, _ := json.Marshal(secret.Data)
			var m map[string][]string
			json.Unmarshal(b, &m)
			newFunction(v, m, r, org, storage)
		} else {
			value := ""
			if len(key) > 0 {
				value = fmt.Sprintf("%s%s", key, v)
			} else {
				value = v
			}
			s, e := r.readSecret(org, value)
			if e != nil {
				storage[fmt.Sprintf("secret/%s/%s", org, value)] = e.Error()
			} else {
				storage[fmt.Sprintf("secret/%s/%s", org, value)] = s.Data
			}
		}
	}
}

func (r *Resource) readSecret(org, key string) (*api.Secret, error) {
	secret, err := r.client.Logical().ReadWithContext(context.Background(), "secret/"+org+"/"+key)
	return secret, err
}

func (r *Resource) readList(org, key string) (*api.Secret, error) {
	secret, err := r.client.Logical().ListWithContext(context.Background(), "secret/"+org+"/"+key)
	return secret, err
}

func (r *Resource) ListSecrets(org string) {
	r.gatherAll(org)

	secretPaths := make([]string, 0, len(r.storage))
	for k := range r.storage {
		secretPaths = append(secretPaths, k)
	}

	for _, path := range secretPaths {
		fmt.Println("")
		fmt.Println(path)
		fmt.Println(strings.Repeat("-", len(path)))

		var mapValue map[string]interface{}
		b, _ := json.Marshal(r.storage[path])
		json.Unmarshal(b, &mapValue)

		fmt.Printf("%-20s  %-80v\n", "Key", "Value")
		fmt.Printf("%-20s  %-80s\n", strings.Repeat("-", 3), strings.Repeat("-", 5))
		for k, v := range mapValue {
			fmt.Printf("%-20s  %-80s\n", k, fmt.Sprintf("%s", v))
		}
	}
}

func (r *Resource) FindByKey(org, key string) {
	r.gatherAll(org)

	secretPaths := make([]string, 0, len(r.storage))
	for k := range r.storage {
		secretPaths = append(secretPaths, k)
	}

	var find Find
	for _, path := range secretPaths {
		var findObject FindObject

		var mapValue map[string]interface{}
		b, _ := json.Marshal(r.storage[path])
		json.Unmarshal(b, &mapValue)

		for k, v := range mapValue {
			if k == key {
				findObject.Path = path

				findObject.KeyValue = append(findObject.KeyValue, KeyValue{
					Key:   key,
					Value: fmt.Sprintf("%s", v),
				})
			}
		}
		find.FindObject = append(find.FindObject, findObject)
	}

	print(find)
}

func (r *Resource) FindByValue(org, value string) {
	r.gatherAll(org)

	secretPaths := make([]string, 0, len(r.storage))
	for k := range r.storage {
		secretPaths = append(secretPaths, k)
	}

	var find Find
	for _, path := range secretPaths {
		var findObject FindObject

		var mapValue map[string]interface{}
		b, _ := json.Marshal(r.storage[path])
		json.Unmarshal(b, &mapValue)

		for k, v := range mapValue {
			if v == value {
				findObject.Path = path

				findObject.KeyValue = append(findObject.KeyValue, KeyValue{
					Key:   k,
					Value: fmt.Sprintf("%s", v),
				})
			}
		}
		find.FindObject = append(find.FindObject, findObject)
	}

	print(find)
}

func print(find Find) {
	for _, fo := range find.FindObject {
		if len(fo.Path) > 0 {
			fmt.Println(fo.Path)
			fmt.Println(strings.Repeat("-", len(fo.Path)))

			fmt.Printf("%-20s  %-80v\n", "Key", "Value")
			fmt.Printf("%-20s  %-80s\n", strings.Repeat("-", 3), strings.Repeat("-", 5))
			for _, kv := range fo.KeyValue {
				fmt.Printf("%-20s  %-80s\n", kv.Key, kv.Value)
			}
			fmt.Println("")
		}
	}
}

type Find struct {
	FindObject []FindObject
}

type FindObject struct {
	Path     string
	KeyValue []KeyValue
}

type KeyValue struct {
	Key   string
	Value string
}
