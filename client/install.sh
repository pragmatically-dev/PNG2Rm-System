#!/bin/bash

BIN_DIR="/home/root/.local/bin/png2rm"
CONFIG_DIR="/home/root/.config/png2rm"
SYSTEMD_DIR="/etc/systemd/system"

CLIENT_SRC="./client"
CONFIG_SRC="./config.yaml"
SERVICE_FILE="client.service"

SERVICE_CONTENT="[Unit]
Description=Client Service
After=home.mount

[Service]
ExecStart=$BIN_DIR/client
Restart=on-failure
RestartSec=5
Enviroment="HOME=/home/root"

[Install]
WantedBy=multi-user.target"

mkdir "$BIN_DIR"
mkdir "$CONFIG_DIR"
mkdir "/home/root/rmdoc"

cp "$CLIENT_SRC" "$BIN_DIR"
if [ $? -ne 0 ]; then
  echo "Error copying client to $BIN_DIR"
  exit 1
fi


cp "$CONFIG_SRC" "$CONFIG_DIR"
if [ $? -ne 0 ]; then
  echo "Error copying config.yaml to $CONFIG_DIR"
  exit 1
fi

echo "Files copied successfully."


read -p "Do you want to edit the config.yaml file? (y/n) " edit_answer

if [[ $edit_answer =~ ^[Yy]$ ]]; then
  nano "$CONFIG_DIR/config.yaml"
  if [ $? -ne 0 ]; then
    echo "Error opening config.yaml with nano"
    exit 1
  fi
fi


echo "$SERVICE_CONTENT" > "$SERVICE_FILE"
mv "$SERVICE_FILE" "$SYSTEMD_DIR"
if [ $? -ne 0 ]; then
  echo "Error creating systemd service file"
  exit 1
fi


systemctl daemon-reload


read -p "Do you want to enable the client service? (y/n) " enable_answer

if [[ $enable_answer =~ ^[Yy]$ ]]; then
  systemctl enable client.service
  if [ $? -ne 0 ]; then
    echo "Error enabling client service"
    exit 1
  fi
  echo "Client service enabled."
fi

read -p "Do you want to start the client service now? (y/n) " start_answer

if [[ $start_answer =~ ^[Yy]$ ]]; then
  systemctl start client.service
  if [ $? -ne 0 ]; then
    echo "Error starting client service"
    exit 1
  fi
  echo "Client service started."
fi

echo "Installation completed successfully."
