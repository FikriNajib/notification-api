
# Notification API

This project is send notif to 3rd party (example: clevertap)

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

This project is send payload to post into 3rd party (example: clevertap) using queue , and then 3rd party will send notif to user
requirement :
1. Mysql
2. Kafka

## Features
This project will provide
- Create Job
- Get Job Status


## Getting Started
### Prerequisites

Golang version 1.18 



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



