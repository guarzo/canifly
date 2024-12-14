# Can I Fly

[![Latest Release](https://img.shields.io/github/v/release/guarzo/canifly)](https://github.com/guarzo/canifly/releases/latest)
[![Build & Test](https://github.com/guarzo/canifly/actions/workflows/test.yml/badge.svg)](https://github.com/guarzo/canifly/actions/workflows/test.yaml)

## Overview

**Can I Fly** is an application that helps EVE Online players quickly determine which ships their characters can pilot and what training they need to reach their goals. By integrating directly with EVE's APIs, **Can I Fly** lets you manage characters, skill plans, and account mappings from a single, intuitive interface. The tool also helps you configure local EVE profiles, sync settings, and back up or restore them effortlessly.

## Key Features

- **Character Management:**  
  Add, remove, and organize your EVE Online characters by account, role, or location. Use the character overview page to see training statuses, skill points, and quick details at a glance.

- **Skill Plans:**  
  View, create, and delete skill plans. Instantly see which characters qualify for specific plans, which are pending, and which require more training. Easily copy skill plans for sharing or backup.

- **Mapping & Syncing:**  
  Associate character files with local user files for seamless syncing of in-game settings. Drag and drop characters to user files, and use the dashboard and mapping pages to ensure all accounts and characters have the right configuration.

- **Integrated Instructions & Tooltips:**  
  Each page (e.g., Character Overview, Skill Plans) includes a built-in instructions toggle. Simply click the help icon to show or hide guidance on how to use that page’s features.

## Prerequisites

- **Go:** Version 1.22.3 or newer. [Download Go](https://golang.org/dl/)
- **npm:** Version 22.2.0 or newer.
- **EVE Developer Credentials:**  
  Create an application via [EVE Online Developers](https://developers.eveonline.com/applications) to get:
    - `EVE_CLIENT_ID`
    - `EVE_CLIENT_SECRET`

- **Callback & Secret Key:**  
  Set `EVE_CALLBACK_URL` to the callback URL configured in your EVE developer application.  
  Generate a secret key for encryption:
  ```sh
  openssl rand -base64 32
  ```
  Use this output as `SECRET_KEY`.

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

After the server starts, visit `http://localhost:3000` in your browser to access **Can I Fly**.

## Usage Tips

- **Character Overview Page:**
    - Use the green plus icon to add new characters by account name.
    - Toggle grouping by account, role, or location.
    - Click on a character’s name for detailed skill and training info.
    - Sort order can be toggled to quickly find characters or accounts.

- **Skill Plans Page:**
    - Switch between "Character" and "Skill Plan" views using the toggle button.
    - View missing skills or pending training time at a glance.
    - Copy or delete skill plans via the action column icons.
    - Add new skill plans using the skill plan icon in the header.

- **Mapping & Syncing:**
    - On the mapping page, drag characters onto user files for easy association.
    - On the sync page, align multiple profiles, backing up and restoring configurations across all characters and accounts.

**Built-in Instructions:**  
Each page has a "Help" icon to show/hide instructions. Keep the UI clutter-free when you know the workflows, or show guidance if you need a refresher.

## Contributing

We welcome contributions! If you find a bug, have a feature request, or want to improve documentation, feel free to open an issue or submit a pull request.

## License

This project is licensed under the ISC License. For details, see the [LICENSE](./LICENSE) file.

