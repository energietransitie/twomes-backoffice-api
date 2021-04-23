from datetime import datetime, timezone

import sqlalchemy.types


class DateTime(sqlalchemy.types.TypeDecorator):
    """
    Custom DateTime, to make sure we always have aware datetime instances
    at the python side, and always store timestamps in the database in UTC.

    This is necessary, as MariaDB, contrary to PostgreSQL, does not support
    storing the timestamp *including* the timezone, in the database.
    """
    impl = sqlalchemy.types.DateTime

    def process_bind_param(self, value: datetime, dialect) -> datetime:
        if value is not None:
            if not value.tzinfo:
                raise TypeError("tzinfo is required")

            value = value.astimezone(timezone.utc).replace(tzinfo=None)

        return value

    def process_result_value(self, value: datetime, dialect) -> datetime:
        if value is not None:
            assert not value.tzinfo
            value = value.replace(tzinfo=timezone.utc)

        return value
