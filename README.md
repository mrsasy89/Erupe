# Erupe Custom
This repository is a fork of [Mezeporta/Erupe](https://github.com/Mezeporta/Erupe), a Monster Hunter Frontier server emulator.
The goal of this fork is to provide a version of Erupe integrated with a custom launcher (originally ButterClient, now mhf-custom-launcher), with a focus on Linux support and easier installation.

## Credits
Original server project: Mezeporta/Erupe

Launcher integration and related code changes: [LilButter](https://github.com/mrsasy89/ButterClient)

This fork (Linux fixes, configuration, documentation): mrsasy89

All code is subject to the license in LICENSE.txt in this repository and in the original project.

## Client compatibility
### Platforms

PC Windows

PC Linux (I am currently working on a native Linux version of MHFZ)

### Client versions (ClientMode)

All versions after HR compression (G10–ZZ) have been tested and work correctly with this fork when used together with the custom launcher.

If you have an installed copy of Monster Hunter Frontier on an old drive, please consider helping to preserve it for archival purposes.

### Requirements
Go (recent version, ideally 1.21 or newer)

PostgreSQL

A database initialized with the Erupe schema and the required patch schemas

mhf-custom-launcher configured to point to your server

### Quick installation
Clone the repository

bash
``git clone https://github.com/mrsasy89/Erupe.git
cd Erupe``
Prepare the PostgreSQL database

Create a database (for example erupe) and a dedicated user.

Restore the provided .backup or .sql dump.

Run all scripts in patch-schema in the correct order (usually by filename).

Configure config.json

Set:

PostgreSQL host, port, database name, user and password

IP/ports for Entrance, Channel and any other servers

Logging options and command options as you prefer

Make sure the launcher-related options match your mhf-custom-launcher configuration.

Build or run the server

bash
``go build -o erupe-server .
./erupe-server``
or, for testing:

bash
``go run .``
Configure the launcher

Install / configure mhf-custom-launcher.

Set the server address and ports so they match your config.json (Entrance and Channel endpoints).

Differences from Mezeporta/Erupe
This fork keeps the core behaviour of Erupe but adds:

Integration and configuration tweaks for the custom launcher (ButterClient → mhf-custom-launcher).

Adjustments to make building and running the server on Linux easier (configuration handling, paths, console behaviour).

Small changes to default configuration and SQL scripts to simplify private server setup.

For full details about the server architecture, schema types (init, update, patch, bundled) and general Erupe features, please refer to the original documentation in Mezeporta/Erupe.
