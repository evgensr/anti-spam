# Anti-Spam Bot

Welcome to the Anti-Spam Bot project! This bot is designed to monitor and manage spam activities, ensuring smooth and uninterrupted conversations.

## Configuration

To configure the bot, navigate to the configuration file located at:
```
config/config.yaml
```

Inside the `config.yaml` file, you will need to specify several parameters:

- `Token`: Place your bot token here.
- `GroupID`: Input the ID of the group where the bot is a member.
- `AdminID`: List the user IDs of the admins. Admin users will receive notifications about user bans.

Here is an example of how your `config.yaml` file should look like:
```yaml
Token: "YOUR_BOT_TOKEN"
GroupID: "YOUR_GROUP_ID"
AdminID:
  - "ADMIN_USER_ID_1"
  - "ADMIN_USER_ID_2"
```

Replace the placeholder text with your actual details.

## Running the Bot

To run the bot, you have two options:

1. Use the provided Makefile:
   ```sh
   make up
   ```

2. Use Docker Compose:
   ```sh
   docker-compose up -d
   ```

Once the bot is running, it will start monitoring the specified group for any spam activities and notify the listed admin users of any bans.

## Contribution

Feel free to contribute to the development of this bot by creating pull requests or reporting any issues you encounter. Every contribution is highly appreciated!

Enjoy a spam-free experience with Anti-Spam Bot!
