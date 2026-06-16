# Jellyfin Template

- A template for an alert if jellyfin begins a transcode.
- Install the webhook plugin and add in the template.
- Tick the `Playback Start.`

```text
{{#if_equals PlayMethod 'Transcode'}}
{
    "content": "{{MentionType}}",
    "username": "{{BotUsername}}",
    "embeds": [
        {
            "author": {
                {{#if_equals ItemType 'Episode'}}
                    "name": "Transcoding • {{{SeriesName}}} S{{SeasonNumber00}}E{{EpisodeNumber00}}",
                {{else}}
                    "name": "Transcoding • {{{Name}}} ({{Year}})",
                {{/if_equals}}
                "url": "{{ServerUrl}}/web/index.html#!/details?id={{ItemId}}"
            },
            "thumbnail":{
                "url": "{{ServerUrl}}/Items/{{ItemId}}/Images/Primary"
            },
            "description": "Click [**here**]({{ServerUrl}}/web/index.html#/dashboard) to open the Admin Dashboard.",
            "color": "16711680",
            "fields": [
                {
                    "name": "User",
                    "value": "{{{NotificationUsername}}}",
                    "inline": true
                },
                {
                    "name": "Device",
                    "value": "{{{DeviceName}}}",
                    "inline": true
                },
                {
                    "name": "Video Details",
                    "value": "Codec: {{Video_0_Codec}}\nProfile: {{Video_0_Profile}}\nLevel: {{Video_0_Level}}",
                    "inline": false
                }
            ],
            "footer": {
                "text": "Method: {{PlayMethod}} | {{{ServerName}}}"
            },
            "timestamp": "{{Timestamp}}"
        }
    ]
}
{{/if_equals}}
````