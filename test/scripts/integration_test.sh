#!/usr/bin/env bash

MODIFIED_IMAGES_DIR="./modified_images"
DOWNLOADED_IMAGES_DIR="./downloaded_images"

mkdir -p "$DOWNLOADED_IMAGES_DIR"

BASE_URL="http://image-previewer:8080/fill"

countdown=5

while [ $countdown -gt 0 ]; do
    echo "$countdown сек"
    sleep 1
    countdown=$((countdown - 1))
done

for img in "$MODIFIED_IMAGES_DIR"/*; do
  filename=$(basename "$img")

  if [[ $filename =~ gopher_([0-9]+)x([0-9]+) ]]; then
    width=${BASH_REMATCH[1]}
    height=${BASH_REMATCH[2]}

    request_url="$BASE_URL/$width/$height/nginx:80/images/_gopher_original_1024x504.jpg"
    downloaded_image="$DOWNLOADED_IMAGES_DIR/$filename"

    curl -s -o "$downloaded_image" "$request_url"

    if [ -f "$downloaded_image" ]; then
      modified_width=$(identify -format "%w" "$img")
      modified_height=$(identify -format "%h" "$img")

      downloaded_width=$(identify -format "%w" "$downloaded_image")
      downloaded_height=$(identify -format "%h" "$downloaded_image")

      if [ "$modified_width" -eq "$downloaded_width" ] && [ "$modified_height" -eq "$downloaded_height" ]; then
        echo "Match found: $filename (Width: $modified_width, Height: $modified_height)"
      else
        echo "No match: $filename (Modified - Width: $modified_width, Height: $modified_height; Downloaded - Width: $downloaded_width, Height: $downloaded_height)"
      fi
    else
      echo "Failed to download image: $request_url"
    fi
  else
    echo "Filename format not recognized: $filename"
  fi
done
