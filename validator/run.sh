#!/bin/sh

MESSAGE='¡Todo es verdad señor, esto es la radio!'

SERVER_IP='server'
SERVER_PORT='12345'

response=$(printf '%s' "$MESSAGE" | nc -w 3 "$SERVER_IP" "$SERVER_PORT" 2>/dev/null)
status=$?

if [ "$status" -eq 0 ] && [ "$response" = "$MESSAGE" ]; then
  echo 'action: test_echo_server | result: success'
  exit 0
fi

echo 'action: test_echo_server | result: fail'
exit 1
