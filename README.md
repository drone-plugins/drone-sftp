# drone-sftp

Drone plugin for uploading files to an SFTP server

# Usage

        node index.js <<EOF
        {
            "repo": {
                "clone_url": "git://github.com/drone/drone",
                "full_name": "drone/drone"
            },
            "build": {
                "number": 1,
                "event": "push",
                "branch": "master",
                "commit": "436b7a6e2abaddfd35740527353e78a227ddcb2c",
                "ref": "refs/heads/master",
                "status": "success"
            },
            "workspace": {
                "root": "/drone/src",
                "path": "/drone/src/github.com/drone/drone"
            },
            "vargs": {
                "host": "sftp.company.com",
                "port": 2222,
                "username": "username",
                "password": "password", 
                "files": [
                    "*.nuget"
                ]
            }
        }
        EOF

# Docker

Build the Docker container:

    docker build -t athieriot/drone-stp .
