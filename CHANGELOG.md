### 12/17/2025

<h3>Updated</h3>

- Switched to using Slack inside of alertmanager with a webhook URL instead of an email address with SMTP
    - slack channel should be called 'shear'

### 12/16/2025
<h3>Added</h3>

- Cron job that runs nightly at 2am to find the time delta between today's date and the last recorded event for a user
    - New package called `scheduler`
    - Time delta threshold is 90 days
    - If a user or list of users exceeds the 90 threshold, a discord message will be sent to a particular channel ID that includese the "@" role 
        - channel ID and discord role are now .env values
- Database seeding function to seed the database with every user in a discord server on init
- New `DetectUserJoin` and `DetectUserLeave` events to detect when a user joins or leaves the server entirely and removes them from the Postgres table
    - Ensures a "best effort" to keep a 1 to 1 with discord user count and DB table congruency
- Terraform project under `/terraform`
    - Includes a DigitalOcean droplet (server) and associated SSH key
    - Backups run nightly at 4am
- Init scripts for installing and configuring Prometheus and Alertmanager for observability
    - Utilizes node_exporter for server level metrics and cAdvisor for container level metrics
    - Routes all alerts to a particular email address
- Init script (init-db.sql) that builds the Postgres table and schema

<h3>Updated</h3>

- Added rare edge case logic to check if the author of a message is the same author when running the `!shear remove-user <username>` and remove them from the table
### 12/15/2025
<h3>Added</h3>

- New bot command that can remove a user from a discord server
    - Usage: `!shear remove-user <username>`

<h3>Updated</h3>

- Reformed project directory struture
- Packages include `config`, `database`, `command`, and `event`

### 12/14/2025
<h3>Added</h3>

- First official deployment of the Shear Discord bot