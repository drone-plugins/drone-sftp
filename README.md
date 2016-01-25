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
                "path": "/drone/src/athieriot/drone-sftp"
            },
            "vargs": {
                "host": "sftp.company.com",
                "port": 2222,
                "username": "username",
                "password": "pa$$word", 
                "destination_path": "/share"
                "files": [
                    "*.md"
                ]
            }
        }
        EOF

# Docker

Build the Docker container:

    docker build -t athieriot/drone-stp .
