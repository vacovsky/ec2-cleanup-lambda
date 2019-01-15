#!/bin/sh



# USAGE:      ./package.sh herheh5ynjey5nje5   
#  Might need to create a bucket first:   aws s3 create-bucket --bucket="herheh5ynjey5nje5"
BUCKETNAME=$1


# Build the lambda code
GOOS=linux go build -o lifecycle

# Zip up the lambda in preparation for deployment
zip lifecycle.zip ./lifecycle
rm -rf lifecycle

# push artifact to S3
aws s3 create-bucket --bucket="$BUCKETNAME"
aws s3 mv lifecycle.zip s3://"$BUCKETNAME"

# try to create stack with CloudFormation
aws cloudformation create-stack --stack-name PipelineInstanceLifecycleLambda --template-body file://stack.json --capabilities CAPABILITY_NAMED_IAM --parameters  ParameterKey=S3BucketName,ParameterValue="$BUCKETNAME"

# try to update stack with CloudFormation
aws cloudformation update-stack --stack-name PipelineInstanceLifecycleLambda --template-body file://stack.json --capabilities CAPABILITY_NAMED_IAM --parameters  ParameterKey=S3BucketName,ParameterValue="$BUCKETNAME"
