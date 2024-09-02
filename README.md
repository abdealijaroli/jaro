# Jaro

## Overview
Jaro is a CLI app built as a developer tool for developers who don't want to leave their terminals. You can achieve two things - shorten links and transfer files, all from your terminal.

## Usage
To get started with Jaro, pull the Docker image from Docker Hub:
```
docker pull abdealijaroli/jaro
```

## Configuration
The application inside this container reads its configuration from a single environment variable. You need to provide this variable when running the container.

### Required Environment Variable
- **DB_URL**: PostgreSQL connection string (e.g., `postgresql://user:password@host:port/dbname`)

## 1. Run the server
To run the Jaro application, use the following command, replacing the `DB_URL` with your actual PostgreSQL connection string:
```
docker run -d --name jaro -e DB_URL=postgresql://user:password@host:port/dbname -p 8008:8008 abdealijaroli/jaro:latest
```

## 2. Shorten a link
```
docker run -d --name jaro -e DB_URL=postgresql://user:password@host:port/dbname -p 8008:8008 abdealijaroli/jaro:latest -s https://example.com
```

## 3. Transfer a file
```
docker run -d --name jaro -e DB_URL=postgresql://user:password@host:port/dbname -p 8008:8008 abdealijaroli/jaro:latest -t /path/to/your/file
```

## Features
- **Link Shortening**: Easily shorten URLs directly from your terminal.
- **File Transfer**: Transfer files without leaving your command line.

## Contributing
If you would like to contribute to Jaro, feel free to submit a pull request or open an issue.

## License
This project is open-source and available under the [MIT License](LICENSE).
