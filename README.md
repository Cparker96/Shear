<h2>Overview</h2>

This README serves as the official documentation for the Shear discord bot and how the service integrates into a particular discord server and tracks user activity in the form of voice events, text channel messages, and emoji reaction events.

<h2>Description</h2>

The bot will perpetually track user activity when a user joins a voice channel, writes a message in a text channel, or appends an emoji (known as a 'reaction') to a message in a text channel. When the bot detects that a user has performed any of the three actions, it sends some related metadata depending on what action occurred into a separate goroutine to write to a PostgreSQL database and table. This should ensure that the action of discovering any of the three user events doesn't block the main thread and concurrency is implemented between indefinite user activity tracking and DB writes.

There are three columns (each of string type) in the Postgres table that the bot writes to. `username` points to the discord user's username, `update_type` points to the specific event type in which the user made (options are `voice`, `message`, or `reaction`), and `date` is the timestamp (expressed as YYYY-MM-DD) for when the event occurred.

<h2>Prerequisites</h2>
To stand up your own implementation of Shear in your discord server, refer to these steps below:

1. Register a new bot under the `Applications` tab at the [Discord Developer Portal](https://discord.com/developers/docs/intro) and sign in
2. Under `Installation`, select `Guild Install`
3. Under `OAuth2`, generate the client secret and store it somewhere safe
    - Select `bot` in the `OAuth2 URL Generator` section
    - Select `Manage Channels`, `View Channels`, `Manage Events`, `Send Messages`, `Manage Messages`, `Attach Files`, and `Read Message History` in the `Bot Permissions` section (if you want to continue developing this bot for your open source needs, you will need to update these permissions depending on your use case)
    - Copy the `Generated URL` string and paste into a browser. This will link your bot into your specific discord server
4. Under the `Bot` section, give it a friendly name and take note of the token. Store this somewhere safe. Make sure `Presence Intent`, `Server Members Intent`, and `Message Content Intent` are all selected
5. Install PostgreSQL on your computer: [Download PostgreSQL](https://www.postgresql.org/download/)
    - Optional: Install a GUI like [pgAdmin](https://www.pgadmin.org/download/) to build your database and table (table name should be `activity`)
6. Create a user in Postgres that has permissions over the DB that you create (it is recommended to not use the postgres master username/password as the authenticated user)
7. Download Docker and Docker Compose
    ```bash
    # install all applicable updates and certs
    for pkg in docker.io docker-doc docker-compose podman-docker containerd runc; do sudo apt-get remove $pkg; done
    sudo apt install -y apt-transport-https software-properties-common ca-certificates curl gnupg lsb-release

    # add Docker's GPG key
    mkdir -p /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    sudo chmod a+r /etc/apt/keyrings/docker.gpg

    # set up Docker repository
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list

    # install latest version
    sudo apt update -y
    sudo apt-get -y install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    ```

<h2>Usage</h2>
To deploy your bot:

1. Clone this repo and create a .env file. Paste your secrets in there
    - Refer to the .env.example file for variable naming conventions
2. Build the Docker image
    ```bash
    docker build -t shear .
    ```
    - If you would rather compile as a binary:
        ```bash
        go build

        ./shear
        ```
3. Deploy the image
    ```bash
    docker compose up -d
    ```

<h2>Available Commands</h2>

- `!shear get-activity`
    - This will pull all of the records in the main postgres table and output all results to a .csv file in a separate message. You can download and open this file with Excel, Google Sheets, etc.