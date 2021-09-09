CREATE TABLE IF NOT EXISTS webhook_events(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    webhook_event text,
    email_provider character varying,
    "to" character varying,
    email_provider_message_id character varying,
    reason text,
    event character varying,
    sd_message_id character varying,
    timestamp integer
);
