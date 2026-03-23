#!/bin/sh

CONFIG_FILE='/config.ini'
MESSAGE='¡Todo es verdad señor, esto es la radio!'

SERVER_IP='server'
SERVER_PORT='12345'

if [ -f "$CONFIG_FILE" ]; then
  ip_line=$(grep -E '^[[:space:]]*SERVER_IP[[:space:]]*=' "$CONFIG_FILE" 2>/dev/null | head -n 1 || true)
  port_line=$(grep -E '^[[:space:]]*SERVER_PORT[[:space:]]*=' "$CONFIG_FILE" 2>/dev/null | head -n 1 || true)

  ip_val=$(printf '%s' "$ip_line" | cut -d= -f2- | tr -d ' \t\r')
  port_val=$(printf '%s' "$port_line" | cut -d= -f2- | tr -d ' \t\r')

  [ -n "$ip_val" ] && SERVER_IP="$ip_val"
  [ -n "$port_val" ] && SERVER_PORT="$port_val"
fi

[ "$SERVER_IP" = 'server' ] && SERVER_IP='127.0.0.1'

response=$(printf '%s' "$MESSAGE" | nc -w 3 "$SERVER_IP" "$SERVER_PORT" 2>/dev/null)
status=$?

if [ "$status" -eq 0 ] && [ "$response" = "$MESSAGE" ]; then
  echo 'action: test_echo_server | result: success'
  exit 0
fi

echo 'action: test_echo_server | result: fail'
exit 1
