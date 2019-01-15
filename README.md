# EC2 cleanup lambda

## Description
Looks for `DateCreated` tag (with unix timestamp for value) and `Ephemeral` tag (`True` as the value) to determine which instances are more than a day old, then terminates those instances.

## Usage

To deploy this, execute the `package.sh` script, with the only parameter being the name of an S3 bucket:

`./deploy.sh mys3bucketname`

# ec2-cleanup-lambda
