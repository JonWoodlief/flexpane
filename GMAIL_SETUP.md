# Gmail Integration Setup

Flexplane now supports Gmail integration to show your real calendar and email data from Google!

## Quick Start (Demo Mode)

By default, Flexplane runs with mock data - no setup required. Just run:

```bash
go run main.go
```

## Gmail Integration Setup

To enable Gmail integration with your real Google data:

### 1. Create Google OAuth Credentials

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Gmail API and Calendar API:
   - Go to "APIs & Services" → "Library"
   - Search for and enable "Gmail API"
   - Search for and enable "Google Calendar API"
4. Create OAuth 2.0 credentials:
   - Go to "APIs & Services" → "Credentials"
   - Click "Create Credentials" → "OAuth 2.0 Client IDs"
   - Choose "Web application"
   - Add `http://localhost:3000/auth/callback` to "Authorized redirect URIs"
   - Copy the Client ID (you don't need the client secret for this setup)

### 2. Configure Flexplane

Set the Google Client ID as an environment variable:

```bash
export GOOGLE_CLIENT_ID="your-client-id-here"
go run main.go
```

Or run directly:

```bash
GOOGLE_CLIENT_ID="your-client-id-here" go run main.go
```

### 3. Sign In

1. Open http://localhost:3000
2. Click "Sign in with Google"
3. Complete the OAuth flow
4. Your real Gmail and Calendar data will now appear!

## Security Notes

- This setup follows Google's OAuth 2.0 best practices
- No client secret is required, making it safe for distributed applications
- Tokens are currently logged for demo purposes - in production, implement secure token storage
- Only read-only access to Gmail and Calendar is requested

## Features

- ✅ Real Gmail email preview (last 10 emails from inbox)
- ✅ Real Google Calendar events (today's events)  
- ✅ Automatic OAuth token refresh
- ✅ Graceful fallback to mock data if not configured
- ✅ No setup required for demo/development use

## Troubleshooting

**"Gmail provider not configured" error**: Make sure `GOOGLE_CLIENT_ID` environment variable is set.

**OAuth errors**: Verify that the redirect URI `http://localhost:3000/auth/callback` is added to your Google OAuth client configuration.

**No data showing**: Check that you have events in your Google Calendar and emails in your Gmail inbox.