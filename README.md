<!-- PROJECT LOGO -->
<br />
<div align="center">
  <h1>RedCarbon Agent</h1>

  <a href="https://github.com/redcarbon-dev/redcarbon-agent">
    <img src="https://github.com/redcarbon-dev.png" alt="Logo" width="80" height="80">
  </a>

  <p align="center">
    <br />
    <a href="https://github.com/redcarbon-dev/redcarbon-agent/blob/main/README.md"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/redcarbon-dev/redcarbon-agent/issues">Report Bug</a>
    ·
    <a href="https://github.com/redcarbon-dev/redcarbon-agent/issues">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project">About The Project</a></li>
    <li><a href="#getting-started">Getting Started</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->

## About The Project

This is the repository containing the source code of the RedCarbon Agent.

<!-- GETTING STARTED -->

## Getting Started

To build the Agent you need:

- [Go](https://go.dev/) 1.24 or later. You'll also need to set your `$GOPATH` and have `$GOPATH/bin` in your path.

## Installation

To install the RedCarbon Agent, you mast download the latest release you can find at the following [page](https://github.com/redcarbon-dev/redcarbon-agent/releases/latest).

We support Linux and MacOS. The binary is statically linked, so you don't need to install any dependencies.

After the download, you can configure the agent by using the 'profile.sh' utility script. This script is used to handle profiles within the agent. You can create, delete and list profiles.

Profiles are associated with a specific agent token.

For adding a new profile, you can use the following command:

```bash
./profile.sh add -t <agent_token> <profile_name>
```

Next, you can start the agent by using the 'svc.sh' utility script. This script is used to install, start, stop and remove the agent as a systemd service daemon.

The script must be run with root privileges, but it is recommended to specify a user to run the agent. This user must have the right permissions to access the downloaded files.

To install the agent as a systemd service, you can use the following command:

```bash
sudo ./svc.sh install <user>
```

## Update

The agent supports automatic updates. The agent will check for updates every 24 hours and will download the latest version if available.

### Manual Update

We suggest to use the automatic update feature, but if you want you can update the agent manually, in case of problems.

For executing the update, you need to stop the agent before and then substitute the binary with the new one. After that, you can start the agent again.

It will keep the same configuration as before, so you don't need to worry about it.
