# Socket Application

This is a UDP socket application that allows for communication between devices on a network, including sending messages and transferring files.

## Docker Setup

### Building the Docker Image

To build the Docker image, run the following command from the root directory of the project:

```bash
docker build -t socket-app .
```

### Running the Docker Container

The application has several commands that can be run:

#### Start the UDP Server

```bash
docker run -p 8080:8080/udp socket-app serve
```

This will start the UDP server that listens for incoming connections and broadcasts its presence to the network.

#### Send a Message to a Specific IP

```bash
docker run socket-app talk --ip <destination_ip> --msg "Your message here"
```

Replace `<destination_ip>` with the IP address of the device you want to send a message to.

#### Send a File to a Specific IP

```bash
docker run -v /path/to/local/files:/files socket-app sendfile --ip <destination_ip> --path /files/your_file.txt
```

Replace `<destination_ip>` with the IP address of the device you want to send the file to, and `/path/to/local/files` with the path to the directory containing the file you want to send.

#### List All Discovered Devices

```bash
docker run socket-app devices
```

This will display a list of all devices discovered on the network, including their IP addresses, ports, and the time since the last heartbeat.

## Including PDF Files in the Docker Image

You can include PDF files in your Docker image by placing them in the `resources/` directory before building the image:

1. Place your PDF files in the `resources/` directory:
   ```bash
   cp your_document.pdf resources/
   ```

2. Build the Docker image:
   ```bash
   docker build -t socket-app .
   ```

The PDF files will be copied into the Docker image and will be available at `/app/resources/` inside the container.

To access these PDF files from within the container, you can use:

```bash
docker run --rm socket-app ls -la /app/resources/
```

This will list all files in the resources directory, including your PDF files.

For convenience, a test script is provided to verify that PDF files are correctly included in the Docker image:

```bash
chmod +x test_pdf_inclusion.sh
./test_pdf_inclusion.sh
```

This script creates a sample PDF file, builds the Docker image, and verifies that the PDF file is included in the image.

## Network Configuration

### Port Exposure

There are several ways to expose ports when running your Docker container:

#### Expose Specific Ports

You can expose specific ports using the `-p` flag:

```bash
docker run -p 8080:8080/udp -p 8081:8081/udp socket-app serve
```

This maps port 8080 and 8081 from the container to the same ports on the host.

#### Expose All Ports

To expose all ports that are defined in the Dockerfile's EXPOSE directives, use the `-P` flag:

```bash
docker run -P socket-app serve
```

This will automatically map all exposed ports to random ports on the host.

#### Expose Ports on the Go

If you need to expose ports dynamically without modifying the Dockerfile, you can:

1. Use the `-p` flag to specify ports at runtime:
   ```bash
   docker run -p 8080-8090:8080-8090/udp socket-app serve
   ```
   This exposes a range of ports from 8080 to 8090.

2. Use the host network to bypass port mapping entirely:
   ```bash
   docker run --network host socket-app serve
   ```
   This allows the container to use the host's network stack, which can be useful for UDP broadcasting and eliminates the need for port mapping.

## Notes

- The default UDP port is 8080. If you need to use a different port, you'll need to modify the Dockerfile and your docker run commands accordingly.
- For file transfers, you need to mount a volume to provide access to the files you want to send.
- The application logs to stdout in JSON format.
