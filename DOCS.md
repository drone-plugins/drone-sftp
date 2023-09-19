Use this plugin to publish artifacts from the build to an SFTP server
You can override the default configuration with the following parameters:

* `host` - Server host
* `port` - Server port, default to 22
* `username` - Server username, default to blank
* `password` - Password for password-based authentication
* `key` - Private key for public-key-based authentication
* `key_path` - Private key path for public-key-based authentication
* `passphrase` - Passphrase of your key for public-key-based authentication (optional)
* `destination_path` - Target path on the server, default to '/'
* `files` - List of files to upload

All file paths must be relative to current project sources

## Examples

Sample configuration using a password in your .drone.yml file:

```yaml
publish:
  sftp:
    host: sftp.company.com
    port: 2222
    username: user
    password: pa$$word
    files:
      - *.nupkg
```

Sample configuration using a private key saved as a secret in your .drone.yml file:

```yaml
publish:
  sftp:
    host: sftp.company.com
    port: 2222
    username: user
    key:
      from_secret: sftp_private_key
    passphrase: my_passphrase
    files:
      - *.nupkg
```
