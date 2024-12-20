require 'net/http'    # Library for HTTP requests
require 'json'        # Library for JSON parsing and formatting
require 'digest'      # Library for SHA256 hashing

# Constants for the API endpoints and maximum retries
BASE_URL = 'http://localhost:8888'
AUTH_ENDPOINT = '/auth'
USERS_ENDPOINT = '/users'
MAX_RETRIES = 3

# Function to fetch the authentication token from /auth endpoint
def fetch_auth_token
  MAX_RETRIES.times do |attempt|
    begin
      uri = URI(BASE_URL + AUTH_ENDPOINT)             # Build the URI for /auth
      response = Net::HTTP.get_response(uri)          # Perform the GET request
      if response.code == '200'                       # Check if the response is successful
        return response['Badsec-Authentication-Token'] # Return the token from the header
      else
        warn "Auth failed with status #{response.code}. Retrying..."
      end
    rescue => e
      warn "Error fetching auth token: #{e.message}. Retrying..." # Log error and retry
    end
    sleep 1  # Wait before retrying
  end
  warn 'Failed to fetch auth token after retries. Exiting.'
  exit(1) # Exit with non-zero status if retries fail
end

# Function to fetch the user IDs from /users endpoint using the checksum
def fetch_user_ids(auth_token)
  # Generate the checksum using SHA256 hash of token + endpoint path
  checksum = Digest::SHA256.hexdigest("#{auth_token}#{USERS_ENDPOINT}")
  MAX_RETRIES.times do |attempt|
    begin
      uri = URI(BASE_URL + USERS_ENDPOINT)           # Build the URI for /users
      request = Net::HTTP::Get.new(uri)              # Create a GET request object
      request['X-Request-Checksum'] = checksum       # Add checksum to the headers
      response = Net::HTTP.start(uri.hostname, uri.port) { |http| http.request(request) }

      if response.code == '200'                      # Check if the response is successful
        return response.body.split("\n")             # Split response body into lines (user IDs)
      else
        warn "Users endpoint failed with status #{response.code}. Retrying..."
      end
    rescue => e
      warn "Error fetching user IDs: #{e.message}. Retrying..." # Log error and retry
    end
    sleep 1  # Wait before retrying
  end
  warn 'Failed to fetch user IDs after retries. Exiting.'
  exit(1) # Exit with non-zero status if retries fail
end

# Main function to orchestrate fetching the token and user IDs
def main
  auth_token = fetch_auth_token             # Fetch the authentication token
  user_ids = fetch_user_ids(auth_token)     # Fetch the user IDs using the token
  puts JSON.generate(user_ids)              # Output the user IDs as JSON
end

main if __FILE__ == $0  # Run the script only if it's executed directly
