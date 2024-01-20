# notification-api
API with pubsub kafka to notification using push-notification(clevertap) , send email(sengrid) , sms (sms alibaba)

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Getting Started](#getting-started)
    - [Installation](#installation)
    - [Command](#command)
- [Step Running]
- [Contributing](#contributing)
- [License](#license)

## Introduction

This project is send payload to post into 3rd party (example: clevertap, sendgrid, sms alibaba) using queue , and then 3rd party will send notif to user
requirement :
1. Mysql
2. Kafka
3. Credential Clevertap
4. Credential Sengrid
5. Credential Alibaba SMS

## Features
This project will provide
- Create Job
- Get Job Status


## Getting Started
### Prerequisites

Go version 1.18 

## Step Running

### Run Fiber (API)

```bash
 $ go run . serve file://.env
```
### Run kafka consumer 
```bash
 $ go run . consumer file://.env
```
 note: available for scaling kafka consumer
## Contributing



