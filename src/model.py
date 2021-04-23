from sqlalchemy import (
    CheckConstraint, Column, Float, ForeignKey,
    Numeric, Integer, Table, Text
)
from sqlalchemy.orm import relationship

from column import DateTime
from db import Base


class Account(Base):
    __tablename__ = 'account'

    PSEUDONYM_MIN = 800000
    PSEUDONYM_MAX = 899999

    id = Column(
        Integer,
        primary_key=True,
        index=True
    )
    pseudonym = Column(
        Integer,
        CheckConstraint(f'{PSEUDONYM_MIN} <= pseudonym <= {PSEUDONYM_MAX}'),
        unique=True,
        comment='Pseudonym identifier, for account reference by 3rd parties'
    )

    created_on = Column(
        DateTime
    )
    activated_on = Column(
        DateTime,
        nullable=True,
        comment='Time at which the activation token is used to active the account'
    )

    activation_token = Column(
        Text,
        unique=True,
        comment='Unique, random token to identify the account during activation'
    )
    session_token_hash = Column(
        Text,
        nullable=True,
        comment='Hash of random, long-lived token to identify the app session for this account'
    )

    building = relationship(
        'Building',
        uselist=False,
        back_populates='account'
    )


class Building(Base):
    __tablename__ = 'building'

    LOCATION_DECIMAL_PLACES = 10
    LOCATION_MAX_DIGITS = 15

    id = Column(
        Integer,
        primary_key=True,
        index=True
    )
    account_id = Column(
        Integer,
        ForeignKey('account.id')
    )

    longitude = Column(
        Numeric(LOCATION_MAX_DIGITS, LOCATION_DECIMAL_PLACES),
        comment='Coarse-grained longitude, for approximate location indication'
    )
    latitude = Column(
        Numeric(LOCATION_MAX_DIGITS, LOCATION_DECIMAL_PLACES),
        comment='Coarse-grained latitude, for approximate location indication'
    )
    tz_name = Column(
        Text,
        comment='Time zone name, in the IANA timezone database format'
    )
    yr_built = Column(
        Text,
        comment='Year built ("oorspronkelijk bouwjaar")'
    )
    type = Column(
        Text,
        nullable=True,
        comment='House type ("woningtype")'
    )
    floor_area = Column(
        Integer,
        nullable=True,
        comment='Floor area ("gebruiksoppervlakte"), as defined in NEN NTA 8800'
    )
    heat_loss_area = Column(
        Integer,
        nullable=True,
        comment='Heat loss area ("verliesoppervlakte"), as defined in NEN NTA 8800'
    )
    energy_label = Column(
        Text,
        nullable=True,
        comment='Energy label, as defined in NEN NTA 8800'
    )
    energy_index = Column(
        Float,
        nullable=True,
        comment='Energy index, as defined in NEN NTA 8800'
    )

    account = relationship(
        'Account',
        back_populates='building'
    )
    devices = relationship(
        'Device',
        back_populates='building'
    )


device_type_property = Table(
    'device_type_property',
    Base.metadata,
    Column(
        'device_type_id',
        Integer,
        ForeignKey('device_type.id')
    ),
    Column(
        'property_id',
        Integer,
        ForeignKey('property.id')
    ),
)


class DeviceType(Base):
    __tablename__ = 'device_type'

    id = Column(
        Integer,
        primary_key=True
    )
    name = Column(
        Text,
        unique=True
    )

    installation_manual_url = Column(
        Text,
        comment='URL to manual with installation instructions'
    )

    devices = relationship(
        'Device',
        back_populates='device_type'
    )
    properties = relationship(
        'Property',
        secondary=device_type_property,
        back_populates='device_types'
    )


class Device(Base):
    __tablename__ = 'device'

    id = Column(
        Integer,
        primary_key=True,
        index=True
    )
    device_type_id = Column(
        Integer,
        ForeignKey('device_type.id')
    )
    building_id = Column(
        Integer,
        ForeignKey('building.id'),
        nullable=True
    )

    proof_of_presence_id = Column(
        Text,
        unique=True,
        comment='Unique, random token to identify the device during activation'
    )
    session_token_hash = Column(
        Text,
        nullable=True,
        comment='Hash of random, long-lived token to identify the device session, after activation'
    )

    created_on = Column(
        DateTime
    )
    activated_on = Column(
        DateTime,
        nullable=True,
        comment='Time at which the proof-of-presence id is used to active the device'
    )

    device_type = relationship(
        'DeviceType',
        back_populates='devices'
    )
    building = relationship(
        'Building',
        back_populates='devices'
    )

    uploads = relationship(
        'Upload',
        back_populates='device'
    )
    measurements = relationship(
        'Measurement',
        back_populates='device'
    )


class Property(Base):
    __tablename__ = 'property'

    id = Column(
        Integer,
        primary_key=True,
        index=True
    )
    name = Column(
        Text
    )
    unit = Column(
        Text,
        nullable=True,
        comment='Unit of property (if any), as defined by the International System of Units'
    )

    device_types = relationship(
        'DeviceType',
        secondary=device_type_property,
        back_populates='properties'
    )

    measurements = relationship(
        'Measurement',
        back_populates='property'
    )


class Upload(Base):
    __tablename__ = 'upload'

    id = Column(
        Integer,
        primary_key=True,
        index=True
    )
    device_id = Column(
        Integer,
        ForeignKey('device.id')
    )

    server_time = Column(
        DateTime,
        comment='Upload time, as reported by the (receiving) server'
    )
    device_time = Column(
        DateTime,
        comment='Upload time, as reported by the (sending) device'
    )
    size = Column(
        Integer,
        comment='Size of upload payload, in bytes'
    )

    device = relationship(
        'Device',
        back_populates='uploads'
    )
    measurements = relationship(
        'Measurement',
        back_populates='upload'
    )


class Measurement(Base):
    __tablename__ = 'measurement'

    id = Column(
        Integer,
        primary_key=True,
        index=True
    )
    device_id = Column(
        Integer,
        ForeignKey('device.id')
    )
    property_id = Column(
        Integer,
        ForeignKey('property.id')
    )
    upload_id = Column(
        Integer,
        ForeignKey('upload.id')
    )

    timestamp = Column(
        DateTime,
        comment='Time of measurement, as reported by the device'
    )
    value = Column(
        Text
    )

    device = relationship(
        'Device',
        back_populates='measurements'
    )
    property = relationship(
        'Property',
        back_populates='measurements'
    )
    upload = relationship(
        'Upload',
        back_populates='measurements'
    )
