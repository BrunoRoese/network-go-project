#!/bin/bash

# This script tests whether PDF files are correctly included in the Docker image

# Create a sample PDF file in the resources directory
echo "Creating a sample PDF file in the resources directory..."
echo "This is a sample PDF file content." > resources/sample.pdf
echo "It's not actually a valid PDF, but it's here to demonstrate" >> resources/sample.pdf
echo "how PDF files can be included in the Docker image." >> resources/sample.pdf

# Build the Docker image
echo "Building the Docker image..."
docker build -t socket-app-test .

# Verify that the PDF file is included in the image
echo "Verifying that the PDF file is included in the image..."
docker run --rm socket-app-test ls -la /app/resources/

echo "If you see 'sample.pdf' in the output above, the PDF file was successfully included in the Docker image."