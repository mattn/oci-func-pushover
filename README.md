# oci-func-pushover

OCI Function to Send Pushover Notifications

This is a Go-based Oracle Cloud Infrastructure (OCI) Function that sends push notifications to Pushover devices or groups. It uses HTTP requests to the Pushover API and is designed to be triggered by OCI Events, Functions, or external calls for real-time alerting, such as monitoring alerts or workflow notifications.

## Features

- **Simple Notifications**: Sends customizable push messages with title, priority, sound, and more.
- **Configurable**: API token and recipient token set via environment variables; extensive options like priority, sound, retry, and URL via env vars.
- **Payload-Driven**: Accepts raw text input directly as the notification message.
- **Error Handling**: Logs API responses and returns HTTP status codes.
- **Output Logging**: Confirms successful sends with receipt ID.

## Prerequisites

- **OCI Account**: Access to OCI Console with permissions to create Functions.
- **OCI CLI and Fn Project**: Installed locally for development and deployment (see [OCI Functions Quickstart](https://docs.oracle.com/en-us/iaas/Content/Functions/Concepts/funconcepts.htm)).
- **Go 1.21+**: For local building (optional, as Functions builds from source).
- **Dependencies**: The code uses `github.com/fnproject/fdk-go` for the function handler and `github.com/gregdel/pushover` for API interactions.
- **Pushover Account**: Sign up at [pushover.net](https://pushover.net) to get your app token and recipient token.

No special Dynamic Group or Policy is required, as it uses direct HTTP API calls. Ensure outbound internet access from your Functions subnet (public subnet or NAT gateway).

## Installation and Deployment

1. **Clone the Repo**:
   ```
   git clone https://github.com/mattn/oci-func-pushover
   cd oci-func-pushover
   ```

2. **Set Up OCI CLI**:
   - Configure OCI CLI with `oci setup config`.
   - Install Fn Project: Follow [OCI Functions Setup](https://docs.oracle.com/en-us/iaas/Content/Functions/Tasks/functionsoci.htm).

3. **Create Functions App**:
   ```
   fn create app <app-name> --annotation oracle.com/oci/subnetIds='["ocid1.subnet.oc1.iad.aaaa..."]'
   ```
   (Replace with your VCN subnet OCID for private access if needed.)

4. **Configure Environment Variables** (in the app; required and optional):
   ```
   fn config app <app-name> PUSHOVER_APP_TOKEN "your-pushover-app-token"
   fn config app <app-name> PUSHOVER_RECIPIENT_TOKEN "your-pushover-recipient-token"
   ```
   - `PUSHOVER_APP_TOKEN`: Your Pushover application API token (required).
   - `PUSHOVER_RECIPIENT_TOKEN`: Your Pushover user key or group key (required).

   Optional variables for customization:
   ```
   fn config app <app-name> PUSHOVER_CALLBACK_URL "https://example.com/callback"
   fn config app <app-name> PUSHOVER_DEVICE_NAME "device-name"
   fn config app <app-name> PUSHOVER_EXPIRE "300s"
   fn config app <app-name> PUSHOVER_HTML "true"
   fn config app <app-name> PUSHOVER_PRIORITY "high"
   fn config app <app-name> PUSHOVER_MONOSPACE "true"
   fn config app <app-name> PUSHOVER_RETRY "30s"
   fn config app <app-name> PUSHOVER_SOUND "siren"
   fn config app <app-name> PUSHOVER_TTL "3600s"
   fn config app <app-name> PUSHOVER_TITLE "Default Title"
   fn config app <app-name> PUSHOVER_URL "https://example.com"
   fn config app <app-name> PUSHOVER_URL_TITLE "Open Link"
   ```
   - `PUSHOVER_CALLBACK_URL`: Callback URL for acknowledgments.
   - `PUSHOVER_DEVICE_NAME`: Target device name.
   - `PUSHOVER_EXPIRE`: Expiration duration for emergency priorities (e.g., "300s").
   - `PUSHOVER_HTML`: Enable HTML formatting ("true").
   - `PUSHOVER_PRIORITY`: Priority level ("high", "low", "emergency", "lowest"; default: "normal").
   - `PUSHOVER_MONOSPACE`: Use monospace font ("true").
   - `PUSHOVER_RETRY`: Retry interval for emergency priorities (e.g., "30s").
   - `PUSHOVER_SOUND`: Notification sound name.
   - `PUSHOVER_TTL`: Message TTL duration (e.g., "3600s").
   - `PUSHOVER_TITLE`: Notification title.
   - `PUSHOVER_URL`: Attached URL.
   - `PUSHOVER_URL_TITLE`: Title for the attached URL.

5. **Deploy the Function**:
   ```
   fn deploy --app <app-name>
   ```

6. **Test Invocation**:
   ```
   echo "Test notification from OCI Function!" | fn invoke <app-name> oci-func-pushover
   ```
   Check logs in OCI Console > Developer Services > Functions > Application Logs.

## Usage

- **Triggering**: Invoke manually via Fn CLI, OCI Console, or trigger via OCI Events Service (e.g., on bucket events or alarms).
- **Behavior**:
  - Accepts raw text input directly as the notification message body.
  - Customizes the message using environment variables for title, priority, sound, etc.
  - Sends POST request to Pushover API (`https://api.pushover.net/1/messages.json`).
  - Returns 200 on success, 500 on error (e.g., invalid API key).
- **Example Input**:
  ```
  Server is down! Check the logs immediately.
  ```
  - The input text is used verbatim as the message.
- **Example Output** (on success):
  ```
  Pushover notification sent successfully: receipt=abc123
  ```

- **Scheduling**: Use OCI Resource Scheduler for periodic calls, or integrate with OCI Monitoring for alarm-based notifications.

## Limitations and Notes

- **Rate Limits**: Pushover API limits 7500 messages/month per app; monitor usage.
- **No Attachments**: Basic text notifications; extend code for images/sounds if needed.
- **Costs**: Function invocations incur OCI charges; Pushover is free for basics.
- **Security**: Store API keys securely in env vars; avoid logging sensitive data.
- **Priority Rules**: For "emergency" priority, set `PUSHOVER_RETRY` and `PUSHOVER_EXPIRE` to avoid API errors.

## Troubleshooting

- **API Errors**: Check logs for Pushover response (e.g., invalid token: 400 Bad Request).
- **Invocation Fails**: Verify input text and env vars with `fn config app <app-name>`.
- **No Notification**: Confirm device registration in Pushover app.
- **Logs**: Use `fn logs <app> -f` or OCI Console for detailed output.

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
