#!/bin/bash

# Flexplane UI Screenshot Tool
# Automatically opens Firefox to localhost:8080 and captures screenshot for Claude feedback

URL="http://localhost:3000"
SCREENSHOT_PATH="../screenshots/current-ui.png"

echo "📸 Taking Flexplane screenshot..."

# Check if server is running
if ! curl -s $URL > /dev/null; then
    echo "❌ Server not running at $URL"
    echo "💡 Start server with: go run main.go"
    exit 1
fi

# Open Firefox (or focus if already open)
echo "🦊 Opening Firefox..."
open -a Firefox "$URL"

# Wait for page to load
sleep 3

# Take screenshot of Firefox window
echo "📷 Capturing screenshot..."
screencapture -l$(osascript -e 'tell app "Firefox" to id of window 1') "$SCREENSHOT_PATH"

if [ $? -eq 0 ]; then
    echo "✅ Screenshot saved: $SCREENSHOT_PATH"
    echo "🔄 Ready for Claude feedback loop!"
else
    echo "❌ Screenshot failed"
    exit 1
fi