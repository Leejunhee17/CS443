import datetime


def now():
    return datetime.datetime.isoformat(datetime.datetime.utcnow())[:-7]
