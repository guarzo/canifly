# Can I Fly

## Description

Provides a way to easily see which ships your characters can fly

## Installation

To install and run the project, ensure you have the following prerequisites:

- Go 1.22.3 or newer installed. You can download and install it from [the official Go website](https://golang.org/dl/).
- Npm version 10.9.0 or newer installed
- EVE Online Developer Application credentials. Set up an application and retrieve the `EVE_CLIENT_ID` and `EVE_CLIENT_SECRET` from [EVE Online Developers](https://developers.eveonline.com/applications).
- Your own base64 encoded secret key for the application (if you lose this, you'll need to clear the data directory and all users will need to reauthenticate).  One option is to use the following command to generate a secret key:

```sh
openssl rand -base64 32 
```

Set the `EVE_CLIENT_ID`, `EVE_CLIENT_SECRET`, `EVE_CALLBACK_URL` and `SECRET_KEY` environment variables in an .env file at the project root


## Usage

To run the application, use the following command:

```sh
npm start
```

## License

This project is licensed under the ISC License. See the LICENSE file in the repository for more information.
