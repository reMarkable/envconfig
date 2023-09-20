# Changelog

## V2

Environment variables explicitly set blank are not handled any differently
that missing variables. That means that required fields, always require a
value (not just the presence of the variable). This also means that default
values will override empty environment variables.

Byte slices expects environment variable values to be Base64 encoded.

Values for map types are semicolon-separated, not comma-separated. The rationale
for this is this enables us to use maps containing slices.
