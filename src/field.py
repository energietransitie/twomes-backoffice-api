from datetime import datetime
from zoneinfo import ZoneInfo, ZoneInfoNotFoundError

from pydantic.datetime_parse import parse_datetime


class Datetime(datetime):

    @classmethod
    def __get_validators__(cls):
        yield cls.validate

    @classmethod
    def validate(cls, v):
        dt = parse_datetime(v)

        # datetime must be aware
        if dt.tzinfo is None or dt.tzinfo.utcoffset(dt) is None:
            raise ValueError('datetime without timezone')

        return dt


class Timezone(str):

    @classmethod
    def __get_validators__(cls):
        yield cls.validate

    @classmethod
    def validate(cls, v):
        try:
            ZoneInfo(v)

        except ZoneInfoNotFoundError:
            raise ValueError('illegal time zone')

        return v
