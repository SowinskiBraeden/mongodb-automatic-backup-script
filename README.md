# MongoDB Backup Script
A script to automatically handle the backing up of MongoDB databases and optionally the storage of the archive to a github repository. Designed as a script to be run as a service on a server.

<br>

### Setup

1. Rename `.env.example` to `.env`
2. Set the parameters of your **mongoURI, database names, and optionally a gitub url/ssh link**. Ex.
    ```.env
    mongoURI=mongodb://localhost:27017
    databases=db_name_1, db_name_2, ...
    github=git@github:Username/MongoDB_Backup_Repo.git
    ```
    **Note:** If you wish to backup more than one database at a time, be sure to seperate the database names with a comma and a space as follows "`, `".

3. Setup your service on your UNIX-based OS, for example an Ubuntu Server.   
    **i.** Build the project to a binary.
    ```bash
    $ go build
    ```
    **ii.** Create a System Service file
    ```bash
    $ vim /lib/systemd/system/mongobackup.service
    ```
    Now within the `mongobackup.service` file, write the configuration for the program.
    ```service
    [Unit]
    Description=MongoDB Backup Service

    [Service]
    Type=simple
    Restart=always
    RestartSec=5
    ExecStart=/path/to/binary/file

    [Install]
    WantedBy=multi-user.target
    ```
    **iii.** Finally start the service.
    ```bash
    $ service mongobackup start
    ```

