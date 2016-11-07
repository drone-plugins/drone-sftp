Use this plugin to publish artifacts from the build to an SFTP server
You can override the default configuration with the following parameters:

* `host` - Server host
* `port` - Server port, default to 22
* `username` - Server username, default to blank
* `password` - Password for password-based authentication
* `destination_path` - Target path on the server, default to '/'
* `files` - List of files to upload

All file paths must be relative to current project sources

## Example

The following is a sample configuration in your .drone.yml file:

```yaml
publish:
  sftp:
    image: plugins/drone-sftp
    host: sftp.company.com
    port: 2222
    username: user
    password: pa$$word
    files: 
      - *.nupkg
```
