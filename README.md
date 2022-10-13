# vault-secret-finder

This command helps to find the key/value in your vault org.

Currently we find data from KV1/Legacy vault format.

## Install in Mac using brew

brew tap evijayan2/vault-secret-finder
brew install vault-secret-finder

## Prerequsite

* Successfly logged into Vault
* Vault token should set in environment as `export VAULT_TOKEN=jscjksdhcjhdwjkch`


## Usage

### To get help

`vault-secret-finder --help`

```
Vault Secret Finder to search/list the vault org.

Usage:
  vault-secret-finder [flags]
  vault-secret-finder [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  find        Find by key or value
  help        Help about any command
  list        List all secrets under the Org

Flags:
      --debug               Enable the log entry to debug mode
  -h, --help                help for vault-secret-finder
      --org string          Vault Org (required)
  -a, --vault-addr string   Vault Address (required)

Use "vault-secret-finder [command] --help" for more information about a command.
```

### To list all secrets & key/value from the Vault ORG

You can set the Vault host address in command param or set as environment variable `export VAULT_ADDR=https://vault.vijay.com`

`vault-secret-finder --org vijay list`

this command will list all the secrets under the org name vijay.

### To find the secret by key 

`vault-secret-finder --org vijay find --key access-key`

This will list all the matching key named `access-key` in vault org

### To find the secret by value

`vault-secret-finder --org vijay find --value bar`

This will list all the matching value named `bar` in vault org

## Future planned items

* Search using KV2
* Find by Value include contains
* Find by Key regex
* Enable autocompletions



