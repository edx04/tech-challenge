# Tech challenge

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

## Setup process

### Clone the repo

 
```shell
git clone https://github.com/edx04/tech-challenge.git

cd tech-challenge
```



### Installing dependencies & building the target 

For this challenge we use the built-in `sam build` to automatically download all the dependencies and package our build target.   
Read more about [SAM Build here](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-build.html) 

 
```shell
sam build
```

## Deployment


To deploy your application for the first time, run the following in your shell:

```bash
sam deploy --guided
```

The command will package and deploy your application to AWS, with a series of prompts:

* **Stack Name**: The name of the stack to deploy to CloudFormation. This should be unique to your account and region, and a good starting point would be something matching your project name.
* **AWS Region**: The AWS region you want to deploy your app to.
* **Confirm changes before deploy**: If set to yes, any change sets will be shown to you before execution for manual review. If set to no, the AWS SAM CLI will automatically deploy application changes.
* **Allow SAM CLI IAM role creation**: Many AWS SAM templates, including this example, create AWS IAM roles required for the AWS Lambda function(s) included to access AWS services. By default, these are scoped down to minimum required permissions. To deploy an AWS CloudFormation stack which creates or modifies IAM roles, the `CAPABILITY_IAM` value for `capabilities` must be provided. If permission isn't provided through this prompt, to deploy this example you must explicitly pass `--capabilities CAPABILITY_IAM` to the `sam deploy` command.
* **Save arguments to samconfig.toml**: If set to yes, your choices will be saved to a configuration file inside the project, so that in the future you can just re-run `sam deploy` without parameters to deploy changes to your application.

You can find your API Gateway Endpoint URL in the output values displayed after deployment.

### Upload csv file to s3

After deploy the aplication you must upload the csv file to s3 bucket

```bash
aws s3 cp data/txns.csv  s3://$(aws s3 ls | grep -oh "transactions-.*$")
```

### Verify email account

For sending emails we use Amazon Simple Email Service (SES), is a AWS’s server-less service, easy to set up, cost-effective solution to send and receive high-volume emails. 

AWS SES requires that every email address (or the domain of that address) must be verified before it can be used in “From”, “Source”, “Sender”, or “Return-Path” address.

Verify Sender address:

```bash
   aws ses verify-email-identity --email-address youremailaddress@example.com

```

Because this aplication is in the sandbox the recipient address also must be verified

```bash
   aws ses verify-email-identity --email-address recipient@example.com

```

## Test the application

for test the apllication you nedt to send the next POST request to the endpoint that you obtain in the deployment step

### POST /email

#### Parameters

| Name     | Required | Type   |
|----------|----------|--------|
| name     | required | string |
| lastName | required | string |
| sender   | required | string |
| reciver  | required | string |


**Example**


Body

{"name": "Edgar","lastName": "Arellano","reciver":"egst-x-306@hotmail.com","sender":"ed04alb@gmail.com"}

Response 
{
    Email id 01000184783aafb5-fcb5e753-610c-4a83-b6e1-43c8495317d2-000000
}






