# myservice (ID: myservice)

Description of my service

## Contents

- [Installation](#Installation)
  - [MESG SDK](#MESG-SDK)
  - [Deploy the Service](#Service)
- [Definitions](#Definitions)
  - [Events](#Events)
    - [started](#started)
  - [Tasks](#Tasks)
    - [taskX](#taskX)

## Installation

### MESG SDK

This service requires [MESG SDK](https://github.com/mesg-foundation/engine) to be installed first.

You can install MESG SDK by running the following command or [follow the installation guide](https://docs.mesg.com/guide/start-here/installation.html).

```bash
npm install -g mesg-cli
```

### Deploy the Service

To deploy this service, go to [this service page](https://marketplace.mesg.com/services/myservice) on the [MESG Marketplace](https://marketplace.mesg.com) and click the button "get/buy this service".

## Definitions

### Events

<h4 id="started">started</h4>

Event key: `started`



| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **x** | `x` | `Boolean` |  |

### Tasks

<h4 id="taskX">taskX</h4>

Task key: `taskX`



##### Inputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **foo** | `foo` | `String` |  |
| **bar** | `bar` | `String` |  |
  
##### Outputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **message** | `message` | `String` |  |
