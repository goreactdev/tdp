# TON Developers Platform System - Backend

The backbone of the TON Developers Platform System, providing the necessary services for user management, TON connectivity, NFT deployment, reward mechanisms, activities, and more. This repository consists of numerous endpoints serving different functions.

## Monitoring System for Queues

Asynqmon is a powerful, web-based monitoring and management tool for [Asynq](https://github.com/hibiken/asynq), a distributed task queue in Go. It enables you to track and manage tasks in your Asynq queues in real-time.

### Overview

Once you have Asynqmon up and running, it provides a graphical interface for you to interact with your Asynq server. The key functionalities include:

- **Dashboard Overview**: Displays a real-time overview of all your Asynq task queues, including active, pending, and failed tasks.
- **Queue Insights**: Provides detailed information on each queue including the number of active, waiting, and completed tasks.
- **Task Management**: Allows for manual task management, including retrying failed tasks and deleting tasks.
- **Scheduler**: Displays upcoming scheduled tasks and allows you to cancel them if necessary.
- **Failure Logs**: Showcases a list of failed tasks, along with error messages and other task details.

### Usage

Asynqmon runs on a web server and is accessed through a web browser. In our system, Asynqmon is located at the **/monitoring** endpoint. You'll see the dashboard with an overview of your task queues. From here, you can click on specific queues or tasks for more information. The interface is intuitive, and actions like retrying or deleting tasks are as simple as clicking a button.

```
https://your-host/monitoring
```

In summary, Asynqmon provides an invaluable interface for monitoring and managing your Asynq tasks, making it easier to track and handle tasks within your application.
## API Endpoints Overview

Below is an overview of the available API endpoints.

#### General

- `GET /v1/status`
- `GET /v1/ton-connect/generate-payload`
- `POST /v1/ton-connect/check-proof`
- `GET /v1/manifest-ton-connect`

#### GitHub Auth

- `GET /v1/github/callback`
- `GET /v1/github/login`

#### Deployed NFT

- `GET /v1/deployed-nft/n/:base64/meta.json`
- `GET /v1/deployed-nft/c/:base64/meta.json`

#### User Related

- `GET /v1/users`
- `GET /v1/users/:username`
- `GET /v1/nfts/:username`
- `POST /v1/admin/csv/upload`
- `POST /v1/admin/media/upload`

#### Group 1 - Require Authenticated User

- `POST /v1/telegram/check_authorization`
- `GET /v1/my-account`
- `PATCH /v1/update/users`
- `PUT /v1/nft/:id/pin`
- `DELETE /v1/unlink/:provider`
- `GET /v1/incoming-achievements`
- `PUT /v1/incoming-achievements/:id`

#### Group 2 - Admin Functions

- `POST /v1/admin/existing-collection`
- `POST /v1/admin/merch`
- `GET /v1/admin/rewards`
- `GET /v1/admin/rewards/:id`
- `DELETE /v1/admin/rewards/:id`
- `GET /v1/admin/users`
- `GET /v1/admin/users/:id`
- `POST /v1/admin/users`
- `DELETE /v1/admin/users/:id`
- `PATCH /v1/admin/users/:id`
- `GET /v1/admin/collections`
- `GET /v1/admin/collections/:id`
- `POST /v1/admin/collections`
- `DELETE /v1/admin/collections/:id`
- `PATCH /v1/admin/collections/:id`
- `GET /v1/admin/prototype-nfts`
- `GET /v1/admin/prototype-nfts/:id`
- `POST /v1/admin/prototype-nfts`
- `DELETE /v1/admin/prototype-nfts/:id`
- `PATCH /v1/admin/prototype-nfts/:id`
- `GET /v1/admin/minted-nfts`
- `GET /v1/admin/minted-nfts/:id`
- `DELETE /v1/admin/minted-nfts/:id`
- `POST /v1/admin/minted-nfts`
- `PATCH /v1/admin/minted-nfts/:id`
- `GET /v1/admin/activities`
- `GET /v1/admin/activities/:id`
- `POST /v1/admin/activities`
- `DELETE /v1/admin/activities/:id`
- `PATCH /v1/admin/activities/:id`
- `GET /v1/admin/permissions`
- `GET /v1/admin/roles`
- `GET /v1/admin/roles/:id`
- `POST /v1/admin/roles`
- `DELETE /v1/admin/roles/:id`
- `PATCH /v1/admin/roles/:id`

## Integration

### POST /v1/admin/merch

This endpoint is specially designed for integrating partners or shops with our platform. It provides an interface to exchange users' ratings for merchandise or other rewards.

#### Request

The request body should include details about the merchandise and the associated user.
```
{
  "name": string,
  "amount": integer,
  "store": string,
  "user_id": integer
}
```

**"name"** - Name of the merchandise.
**"amount"** - How much it cost.
**"store"** - The store or partner providing the merchandise.
**"user_id"** - The user who is eligible to claim the merchandise.

**Note**: Only shops or partners with the required permissions can access this endpoint. Ensure your account has the necessary rights before trying to interact with this endpoint.

#### Response

Upon successful integration, the server responds with a `200 OK` status code, indicating the merchandise has been successfully linked with the specified user rating.

```
{
  "id": integer,
  "name": string,
  "amount": integer,
  "store": string,
  "user_id": integer
}
```

#### Usage

This endpoint can be used to create an interactive and dynamic environment where users are rewarded for their contributions and achievements on the platform. It encourages user engagement and incentivizes high-quality participation.

## Swagger API Documentation

For a complete API reference, please refer to our [Swagger API Documentation](https://app.swaggerhub.com/apis-docs/GOREACTDEV12/TDP/2.0.0).
