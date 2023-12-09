# IAMme IAMme

![](https://github.com/notdodo/IAMme-IAMme/actions/workflows/gosec.yml/badge.svg)
![](https://github.com/notdodo/IAMme-IAMme/actions/workflows/gobuild.yml/badge.svg)

## ChatGPT definition

> "IAMme IAMme" seems to be an unusual and distinctive name. It appears to emphasize "IAM," which stands for Identity and Access Management, and could work as a creative and memorable name for a tool related to Okta or similar IAM systems.

Visualize your Okta tenant.

![Screenshot_20231208_151619](https://github.com/notdodo/IAMme-IAMme/assets/6991986/9e67f882-59c7-45ea-a847-6276b3943ca5)

IAMme is a tool designed to visualize the connections between entities within an Okta tenant, including:

- Users
- Groups
- Group Rules
- Group Memberships

## Usage

To get started, set up the required environment variables by copying the example file:

`cp .env_example .env`

Edit the new `.env` file according to your preferences.

Start the Neo4j database:

`make start`

Dump Okta information into the database:

`make build`

`./iamme dump -u your-okta.okta.com -c "<YourOktaAPIToken" --verbose --debug`

Always leverage on the `--help` flag for getting useful information about command line options.

## Why

The primary aim of IAMme is to provide Okta administrators with a powerful tool for visualizing and analyzing their Okta tenant. While it currently lacks a specific use case, it is developed as a production-grade open-source tool with a focus on security, incorporating a proper development life cycle and architectural design.

IAMme can be particularly helpful for Okta administrators in identifying:

- lots of empty groups
- high numbers of identical rules
- incorrectly configured rules

Although similar information can be obtained using tools like `steampipe``, IAMme offers the capability to correlate and visualize these details using Neo4j.

Feel free to leverage IAMme for:

- Understanding potential security issues
- Analyzing the Okta tenant design

Please note that this project is evolving and open to contributions and improvements.
