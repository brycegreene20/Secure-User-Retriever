import sys
import json
import requests
import hashlib
import time

BASE_URL = "http://localhost:8888"
AUTH_ENDPOINT = "/auth"
USERS_ENDPOINT = "/users"
MAX_RETRIES = 3


def fetch_auth_token():
    """Fetch the authentication token from /auth."""
    for attempt in range(MAX_RETRIES):
        try:
            response = requests.get(BASE_URL + AUTH_ENDPOINT)
            if response.status_code == 200:
                return response.headers.get("Badsec-Authentication-Token")
            else:
                sys.stderr.write(f"Auth failed with status {response.status_code}. Retrying...\n")
        except requests.RequestException as e:
            sys.stderr.write(f"Error fetching auth token: {e}. Retrying...\n")
        time.sleep(1)  # Retry delay
    sys.stderr.write("Failed to fetch auth token after retries. Exiting.\n")
    sys.exit(1)


def fetch_user_ids(auth_token):
    """Fetch the user IDs from /users using the checksum."""
    checksum = hashlib.sha256(f"{auth_token}{USERS_ENDPOINT}".encode()).hexdigest()
    headers = {"X-Request-Checksum": checksum}

    for attempt in range(MAX_RETRIES):
        try:
            response = requests.get(BASE_URL + USERS_ENDPOINT, headers=headers)
            if response.status_code == 200:
                return response.text.strip().split("\n")
            else:
                sys.stderr.write(f"Users endpoint failed with status {response.status_code}. Retrying...\n")
        except requests.RequestException as e:
            sys.stderr.write(f"Error fetching user IDs: {e}. Retrying...\n")
        time.sleep(1)  # Retry delay
    sys.stderr.write("Failed to fetch user IDs after retries. Exiting.\n")
    sys.exit(1)


def main():
    """Main function to retrieve and print the user IDs."""
    auth_token = fetch_auth_token()
    user_ids = fetch_user_ids(auth_token)
    print(json.dumps(user_ids))  # Output to stdout in JSON format
    sys.exit(0)


if __name__ == "__main__":
    main()
