<h2>Overview</h2>

This README serves as the official documentation for the Shear discord bot and how the service integrates into a particular discord server and tracks user activity in the form of voice events, text channel messages, and emoji reaction events.

<h2>Description</h2>

The bot will perpetually track user activity when a user joins a voice channel, writes a message in a text channel, or appends an emoji (known as a 'reaction') to a message in a text channel. When the bot detects that a user has performed any of the three actions, it sends some related metadata depending on what action occurred into a separate goroutine to write to a PostgreSQL database and table. This should ensure that the action of discovering any of the three user events doesn't block the main thread and concurrency is implemented between indefinite user activity tracking and DB writes.

There are three columns (each of string type) in the Postgres table that the bot writes to. `username` points to the discord user's username, `update_type` points to the specific event type in which the user made (options are `voice`, `message`, or `reaction`), and `date` is the timestamp (expressed as YYYY-MM-DD) for when the event occurred.

For server hosting, Shear is standardized on utilizing the DigitalOcean (DO) Terraform provider and building an SSH key, firewall, and droplet resource. Once the server is up and running through `terraform apply`, copy the three "install" scripts (located in `/scripts`) to the server and execute them. For `install_prometheus.sh`, you will need to export two environment variables (titled `SMTP_PASSWORD` and `ALERT_EMAIL`). these correspond with the 16 digit SMTP password associated with your gmail account along with your gmail address (visit https://myaccount.google.com/apppasswords to generate a new SMTP password).

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

<h2>Usage</h2>
To deploy your bot:

1. Build your server through `terraform init/plan/apply` in the terraform directory
2. Clone this repo and copy the install scripts over located in the scripts directory and execute
    - Validate that the slack webhook URL gets injected into the alertmanager.yml file in `install_prometheus.sh`
    - Ex: `./install_prometheus.sh 'your-webhook-URL'
3. Create a .env file and paste your secrets
    - Refer to the .env.example file for variable naming conventions
4. Build the Docker image
    ```bash
    docker build -t shear .
    ```
    - If you would rather compile as a binary:
        ```bash
        go build

        ./shear
        ```
5. Deploy the image
    ```bash
    docker compose up -d
    ```

<h2>Available Commands</h2>

- `!shear get-activity`
    - Pulls all records in the main postgres table and output all results to a .csv file in a separate message. You can download and open this file with Excel, Google Sheets, etc.
- `!shear remove-user <username>`
    - Removes a user from a discord server with a parameter of 'username'