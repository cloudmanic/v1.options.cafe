#!/bin/bash

# TODO: update how we do the version

cd ../
aws s3 sync build s3://cdn.options.cafe/app --acl=public-read
cd scripts