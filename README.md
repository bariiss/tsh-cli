
# File Upload Utility

This is a simple command-line utility for uploading files to a specified server. It supports setting a maximum number of days to keep the file and a maximum number of downloads.

## Features

- Upload a file to a specified server
- Set maximum number of days to keep the file
- Set maximum number of downloads
- Displays upload progress
- Copies the download URL to the clipboard upon successful upload

## Requirements

- Go 1.16 or higher
- Environment variables for the server URL, HTTP authentication username, and password

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/file-upload-utility.git
    cd file-upload-utility
    ```

2. Build the utility:

    ```sh
    go build -o file-upload
    ```

## Usage

1. Set the environment variables:

    ```sh
    export TSH_URL="https://yourserver.com/upload"
    export TSH_HTTP_AUTH_USER="yourusername"
    export TSH_HTTP_AUTH_PASS="yourpassword"
    ```

2. Run the utility with the file you want to upload:

    ```sh
    ./file-upload -max-days 7 -max-downloads 5 path/to/your/file.txt
    ```

    - `-max-days`: (Optional) Maximum number of days to keep the file
    - `-max-downloads`: (Optional) Maximum number of times the file can be downloaded

## Example

```sh
./file-upload -max-days 7 -max-downloads 5 example.txt
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.