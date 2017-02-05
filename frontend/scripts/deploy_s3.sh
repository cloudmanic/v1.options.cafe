#!/bin/bash

cd ../
ng build
aws s3 sync build s3://cdn.options.cafe/app --acl=public-read
cd scripts