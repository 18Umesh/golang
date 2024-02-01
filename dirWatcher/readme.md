# Directory Watcher

The Directory Watcher is a file monitoring system designed to keep track of changes within a specified directory. It provides real-time updates on file creations, modifications, and deletions. The project includes a RESTful API for configuration and status retrieval, as well as a database to log relevant information.

# Key Features

Real-time Monitoring: Detects and reports changes to files in the monitored directory.
Configuration API: Allows users to dynamically configure the monitoring settings using a RESTful API.
Database Integration: Logs file-related events, runtime, and other relevant information in a database.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Usage](#usage)

## Prerequisites

Ensure you have the following installed:

Go (Golang)
PostgreSQL
fsnotify library for Go (go get -u github.com/fsnotify/fsnotify)
GORM library for Go (go get -u gorm.io/gorm)

## Getting Started

Provide instructions on how to set up and run your project. Include step-by-step instructions and code snippets. For example:

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/your-project.git

2. Install Dependencies:
    go get -u github.com/fsnotify/fsnotify
    go get -u gorm.io/gorm

3. Open the config.yaml file in the project root directory.

    Customize the configuration settings based on your preferences. Key settings include:

    Directory: The path to the directory you want to monitor.
    TimeInterval: The time interval for monitoring changes (e.g., "5m" for 5 minutes).
    MagicString: The string to identify in the files.
    Save the changes.

4. Make sure PostgreSQL is running on your machine.Create a new PostgreSQL database for the project and then Replace your_database_name with the desired name.
Open the config.yaml file again and update the database connection details:
Update your_username, your_password, and your_database_name accordingly.


# USAGE

1. Run the main application file.The application will start monitoring the specified directory, and you'll see logs in the console for file changes.

2. Use the RESTFul API to interact with the Directory Watcher

# API Reference

| **Endpoint** | **Method** | **Request Body** | **Sample Response**| **Description** |
|--------------|------------|-------------------|----------------------|-------------------|
| `/start`     | POST       | None              | `{ "message": "Task started successfully." }` | Start monitoring the specified directory. |
| `/config`    | POST        | ```json { "directory": "D:/example", "timeInterval": "10m", "magicString": "NewMagicString" }``` | `{ "message": "Configuration updated successfully." }` | Update the configuration settings. |
| `/stop`      | POST       | None              | `{ "message": "Task stopped successfully." }` | Stop the directory watcher. |
| `/task-run`  | GET        | None              | ```json { "startTime": "2024-02-01T10:00:00Z", "endTime": "2024-02-01T10:15:00Z", "runtime": "15m", "filesAdded": ["file1.txt", "file2.txt"], "filesDeleted": ["file3.txt"], "magicStringHits": 3, "status": "success" }``` | Retrieve details of the latest task run. |

Notes: 
1. The Application automatically logs in changes in the directory into the console and the database.
2. On Start of the application without manually starting it, the application uploads data into the database regarding files present and magic_string occurance.
3. On deletion of a file the data is logged into the database and console automatically.
4. On Creation of a new file, We need to manually start the application for it to read the newly created file.
5. Any modifications in a file regarding its content is logged in the console and the database.
6. You can use the /configure API to change the config values.
7. The task-run will get the most recent event that has been logged into the database.
8. Start and stop will work to manually override the application process.

# END
