# Gmail Provider Setup

The Gmail provider allows Flexpane to display your Gmail inbox. This document explains how to set up Gmail API access.

## Prerequisites

1. A Google account with Gmail
2. Access to [Google Cloud Console](https://console.cloud.google.com/)

## Setup Steps

### 1. Create a Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Click "Create Project" or use an existing project
3. Note your project ID

### 2. Enable Gmail API

1. In your project, go to "APIs & Services" > "Library"
2. Search for "Gmail API"
3. Click on it and press "Enable"

### 3. Create OAuth 2.0 Credentials

1. Go to "APIs & Services" > "Credentials"
2. Click "Create Credentials" > "OAuth 2.0 Client ID"
3. If prompted, configure the OAuth consent screen first:
   - Choose "External" (unless you're in a Google Workspace)
   - Fill in required fields (app name, user support email, etc.)
   - Add your email to "Test users" during development
4. For Application Type, choose "Desktop application"
5. Give it a name (e.g., "Flexpane Gmail")
6. Download the JSON file and save it securely

### 4. Get Access and Refresh Tokens

You'll need to exchange your OAuth 2.0 credentials for access tokens. You can use Google's OAuth 2.0 Playground:

1. Go to [OAuth 2.0 Playground](https://developers.google.com/oauthplayground/)
2. Click the gear icon (⚙️) in the top right
3. Check "Use your own OAuth credentials"
4. Enter your Client ID and Client Secret from step 3
5. In the left sidebar, find "Gmail API v1" and select `https://www.googleapis.com/auth/gmail.readonly`
6. Click "Authorize APIs" and complete the OAuth flow
7. Click "Exchange authorization code for tokens"
8. Copy the `access_token` and `refresh_token`

### 5. Configure Flexpane

Edit `config/providers.json`:

```json
{
  "default": "gmail",
  "providers": {
    "mock": {
      "type": "mock"
    },
    "gmail": {
      "type": "gmail",
      "config": {
        "client_id": "your-client-id.googleusercontent.com",
        "client_secret": "your-client-secret",
        "access_token": "your-access-token",
        "refresh_token": "your-refresh-token"
      }
    }
  }
}
```

Replace the placeholder values with your actual credentials.

## Security Notes

- **Never commit credentials to version control**
- Keep your `client_secret`, `access_token`, and `refresh_token` private
- Consider using environment variables instead of config files for production
- The access token expires, but the refresh token can generate new access tokens

## Usage

Once configured, restart Flexpane and it will display your Gmail inbox in the email pane.

## Troubleshooting

### "Invalid credentials" errors
- Verify your client ID and secret are correct
- Check that your access token hasn't expired
- Ensure the Gmail API is enabled in your project

### "Insufficient permissions" errors  
- Make sure you authorized the `gmail.readonly` scope
- Check that your app is properly configured in OAuth consent screen

### No emails showing
- Check that you have emails in your inbox
- Verify the Gmail API quota hasn't been exceeded
- Look at application logs for specific error messages

## Limitations

- Currently only displays emails (no calendar support via Gmail API)
- Read-only access (cannot send emails or modify)  
- Limited to 10 most recent inbox messages
- For calendar events, you'd need to set up Google Calendar API separately

## Future Improvements

- Support for Google Calendar API integration
- OAuth 2.0 web flow for easier setup
- Better error handling and user feedback
- Support for multiple Gmail accounts