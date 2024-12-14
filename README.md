# Can I Fly

[![Latest Release](https://img.shields.io/github/v/release/guarzo/canifly)](https://github.com/guarzo/canifly/releases/latest)
[![Build & Test](https://github.com/guarzo/canifly/actions/workflows/test.yaml/badge.svg)](https://github.com/guarzo/canifly/actions/workflows/test.yaml)

## Overview

**Can I Fly** is an application that helps EVE Online players quickly determine which ships their characters can pilot. By integrating with EVE Online’s APIs, the tool provides an easy and intuitive interface to view skill plans, training statuses, and ship requirements.

## Features

- **Character Identity Processing:** Fetch and update character information including skills, training queue, and location.
- **Skill Plan Integration:** Determine which skill plans a character has qualified for, is pending, or requires additional training.
- **Mapping & Sync Tools:** Associate EVE accounts and characters with local profiles, syncing settings, and directories.
- **User-Friendly Interface:** Easily select roles, add or remove characters, and adjust skill plans through a simple UI.

## Prerequisites

- **Go:** Version 1.22.3 or newer. [Download Go](https://golang.org/dl/)
- **npm:** Version 22.2.0 or newer.
- **EVE Developer Credentials:**  
  Sign up at [EVE Online Developers](https://developers.eveonline.com/applications) to create an application and obtain:
    - `EVE_CLIENT_ID`
    - `EVE_CLIENT_SECRET`
- **Callback & Secret Key:**  
  Set `EVE_CALLBACK_URL` to the callback URL configured in your EVE developer application.  
  Generate a secret key for encryption:
  ```sh
  openssl rand -base64 32
  ```
  Use this output as `SECRET_KEY`
## Environment Setup

Create a `.env` file at the project root with the following variables:

```
EVE_CLIENT_ID=<your_client_id>
EVE_CLIENT_SECRET=<your_client_secret>
EVE_CALLBACK_URL=<your_callback_url>
SECRET_KEY=<your_generated_secret_key>
```

## Installation

1. **Clone the Repository:**
   ```sh
   git clone https://github.com/guarzo/canifly.git
   cd canifly
   ```

2. **Install Dependencies:**
   ```sh
   npm install
   ```

3. **Build and Run:**
   ```sh
   npm start
   ```

## Usage

- Once running, open your browser and navigate to the URL specified by the application (usually `http://localhost:3000`).
- Log in with your EVE Online credentials (through SSO) to view and manage your character data.
- From the dashboard, you can review skill plans, check if you can fly certain ships, and manage account and character associations.

## Contributing

Contributions are welcome! Please open an issue or create a pull request if you find a bug or want to propose a new feature.

## License

This project is licensed under the ISC License. For details, see the [LICENSE](./LICENSE) file.





