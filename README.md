# doug
Solidity contract artifact release manager

![Build status](https://travis-ci.org/BadgeForce/doug.svg?branch=master)

## Overview 
Exposing the JSON ABI for your smart contracts enables opensource contributions, development and makes collaboration easier for your daaps (decentralized applications) . Even if 
you are not developing an daap for public collaboration it is useful to expose these JSON files within your dev organization. These JSON files are the product of [truffle](http://truffleframework.com/) compilation when using truffle to build daaps. This project is a [github webhook](https://developer.github.com/v3/repos/releases/#get-a-single-release) server written in Go that will upload your configured JSON ABI files to S3 for public consumption whenever you release a new version of your daap.

## Usage 

#### Configuration 

- aws: Your s3 regions and bucket configuration 
The server will attempt to upload each JSON file to all your S3 regions

- projects: This is the configuration for different projects your server will handle
    - name: The name of the project on github 
    - artifacts: File names of the JSON ABI's you need uploaded

```linux
[aws]
regions = ["us-east-1"]
bucket = "badgeforce-artifacts"

[[projects]]
name = "badgeforce"
artifacts = ["Issuer.json", "Holder.json"]

[[projects]]
name = "super-daap"
artifacts = ["ABI.json", "Token.json"]
```

Build and run

```linux 
$ go build 
$ ./doug
```

Doug will upload your JSON ABI files to your S3 bucket in this form

    bucket-name/project-name/version/file-name