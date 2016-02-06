## Overview

A simple Go utility that provides a way to do some housekeeping on an AWS account by deleting resources older than X number of days.

## Setup


### Setup Requirements

The project comes with 1 yaml configuration file as example:
  - config.yaml
  
  Which should be like:

  ```
  ---
resources:
  ec2:
    delete_older_than: 30
  ebs:
    delete_older_than: 45
    whitelist:
      - vol-a1df
  snapshots:
    delete_older_than: 90
    whitelist:
      - snap-xyz
      - snap-abc

  ```

Currently supported resources to remove are ec2 instances, ebs volumes and snapshots represented in the config.yaml as
ec2, ebs and snapshots respectively. Removing any of the 3 resource types from the config will ignore the resources completely. Whitelisted resources by resource id will also be ignored irrespectively of how old they are.


If you run it outside an AWS instance you will need to have an AWSCLI profile file for the accounts you want to run the tool on. Otherwise you can use an IAM profile role when run inside an AWS instance on the same account you want to clean.


By default the tool does a dry run warning you of the resources it will delete based on the set configuration.
The -destroy-resources parameter will destroy all resources matched by the configuration rules without confirmation!!!



### Beginning with awscleaner

## Usage

The tool accepts the following parameters:

```
Please provide the following required parameters:
  -cfgfile string
      Location of custom config yaml file (default "config.yaml")
  -config string
      Location of aws config profiles to use
  -destroy-resources
      If set all resources matching the config rules will be destroyed without warning!!!
  -profile string
      Name of the profile to use when an AWS config with profiles is used 
  -region string
      AWS Region to use (default "eu-west-1")
  -userole
      Use instance IAM role instead of aws config file
```

To run it on the cli for an account with profile name project-dev:

```
awscleaner -config /path/to/aws/profile/config -profile project-dev

```

The above can be performed using an IAM instance role by using the parameter -userole and with the following access profile in IAM:

```
"Action": [
    "ec2:Describe*",
    "ec2:Delete*"
],

```

### Limitations

This tool currently supports the most costly resources to clean up. 
EBS volumes attached to an instance will be ignored.
Snapshots bound to an active AMI will not be deleted and fail with an error.







