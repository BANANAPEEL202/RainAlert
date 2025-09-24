# RainAlert
A lightweight Lambda function that sends a push notification when rain/slow is expected within the next couple of hours. 

## Prerequisites
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sso.html)
- [SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html)

## Configuration
The following configuration fields must be set in ```cmd/lambda/config.json```
- ```latitude/longitude```: coordinates of the location of interest
- ```location```: human-friendly name for the location
- ```timezone```: IANA [timezone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones) identifier for the location
- ```forecast_range_hrs```: future-looking range of each forecast
- ```ntfy_times```: one or more hours of the day (0–23) at which notifications should be sent. Can either be a single int or list of ints
- ```ntfy_topic```: [ntfy](https://docs.ntfy.sh/) notificaton topic. All ntfy topics are public, so choose something that is not easily guessable
- ```ignore_no_rain```: if true, will not send a notification if no rain is expected within the next ```forecast_range_hrs```

## Deploy
1. **Build the project**
   ```
   $ make
   ```
2. **Deploy with SAM**
  - First time deployment
    ```
    $ sam deploy --guided
    ```
  - Subsequent deployments
    ```
    $ sam deploy
    ```

## Usage
- Install the [ntfy](https://ntfy.sh/) app (iOS/Android/Desktop).
- Subscribe to the ```ntfy_topic``` configured in config.json

## Other notes
RainAlert uses the following AWS services and operates well under the free-tier limits for each service:
- [Lambda](https://aws.amazon.com/lambda/): runs hourly
- Eventbridge: schedules/triggers lambda function
- CloudFormation: infrastructure as code service
- Cloudwatch: minimal logging
