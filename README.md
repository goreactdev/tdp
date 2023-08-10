# TON Developers Platform
This repository contains the comprehensive system of the TON Developers Platform, encompassing a robust backend and frontend for users and administrators, an intuitive admin panel, as well as a versatile TON Developers Platform Bot and SBT Minting System. 

## Description

The TON Developers Platform System is designed to foster an active and vibrant developer community within the TON ecosystem. By facilitating user profiles, developer rankings, hackathons, event integrations, and more, we aim to provide a seamless user experience while ensuring security, scalability, and performance. 

## Features

- **Backend and Frontend Development for Users**: A robust backend system catering to user profiles, developer rankings, hackathon and event integrations, and bot integration for users.
- **Admin Panel**: An efficient admin panel focusing on user management, SBT token management, activity and achievement management, notification management, reporting and analytics.
- **TON Developers Platform Bot**: A versatile bot handling account management, notifications, and integration with the TON Developers Platform System.
- **SBT Minting System**: A secure system integrated into the admin panel for minting SBT tokens for various activities and events, providing a user-friendly interface for managing the minting process.

## Getting Started

```bash
git clone https://github.com/goreactdev/tdp
```
## Installation

### For production:
Run the following command to set up the system for a production environment:
```
docker-compose -f docker-compose.prod.yml up -d
```

### For development:
Run the following command to set up the system for a development environment:
```
docker-compose -f docker-compose.dev.yml up -d
```

## Documentation

For detailed technical specifications, please refer to the respective subrepositories: [Backend](backend), [Frontend](frontend), and [Admin](admin).

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
