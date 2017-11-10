#!/bin/bash

cd ../
ng build -prod
aws s3 sync dist s3://cdn.options.cafe/app --acl=public-read
cd scripts