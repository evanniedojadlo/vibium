#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "$0")" && pwd)"
v="${script_dir}/../clicker/bin/vibium"
out="${script_dir}/../saucedemo-record-cli.zip"

# Start daemon (visible browser, background process)
"$v" daemon start
trap '"$v" daemon stop 2>/dev/null || true' EXIT

# Set viewport
"$v" viewport 1280 720

# Start recording with screenshots
"$v" record start --screenshots --name "saucedemo-e2e" --title "SauceDemo E2E Test" --format jpeg --quality 0.1

# 1. Logging in
"$v" record group start "Logging in"
"$v" go "https://www.saucedemo.com"
"$v" find "#user-name"
"$v" fill "@e1" "standard_user"
"$v" find "#password"
"$v" fill "@e1" "secret_sauce"
"$v" find "#login-button"
"$v" click "@e1"
"$v" sleep 500
"$v" record group stop

# 2. Selecting products
"$v" record group start "Selecting products"
"$v" find "#add-to-cart-sauce-labs-backpack"
"$v" click "@e1"
"$v" find "#add-to-cart-sauce-labs-bike-light"
"$v" click "@e1"
"$v" find "#add-to-cart-sauce-labs-onesie"
"$v" click "@e1"
"$v" find ".shopping_cart_badge"
badge=$("$v" text "@e1")
if [ "$badge" != "3" ]; then
    echo "FAIL: Expected cart badge '3', got '$badge'" >&2
    exit 1
fi
echo "Cart badge: $badge"
"$v" record group stop

# 3. Reviewing cart
"$v" record group start "Reviewing cart"
"$v" find ".shopping_cart_link"
"$v" click "@e1"
"$v" sleep 300
"$v" find "#remove-sauce-labs-bike-light"
"$v" click "@e1"
"$v" record group stop

# 4. Checking out
"$v" record group start "Checking out"
"$v" find "#checkout"
"$v" click "@e1"
"$v" find "#first-name"
"$v" fill "@e1" "Test"
"$v" find "#last-name"
"$v" fill "@e1" "User"
"$v" find "#postal-code"
"$v" fill "@e1" "90210"
"$v" find "#continue"
"$v" click "@e1"
"$v" sleep 300
"$v" record group stop

# 5. Completing order
"$v" record group start "Completing order"
"$v" find "#finish"
"$v" click "@e1"
"$v" sleep 500
"$v" find ".complete-header"
confirmation=$("$v" text "@e1")
if [[ "$confirmation" != *"Thank you"* ]]; then
    echo "FAIL: Unexpected confirmation: '$confirmation'" >&2
    exit 1
fi
echo "Confirmation: $confirmation"
"$v" record group stop

# 6. Logging out
"$v" record group start "Logging out"
"$v" find "#react-burger-menu-btn"
"$v" click "@e1"
"$v" sleep 400
"$v" find "#logout_sidebar_link"
"$v" click "@e1"
"$v" sleep 300
"$v" find "#login-button"
loginBtn=$("$v" text "@e1")
echo "Back on login page: $loginBtn"
"$v" record group stop

# Stop recording & save
"$v" record stop -o "$out"
echo "Recording saved → $out"
