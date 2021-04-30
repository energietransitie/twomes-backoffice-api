from typing import Dict
import csv
import os

from model import (
    DeviceType,
    Property,
)
from test.common import session as db

script_dir = os.path.dirname(os.path.abspath(__file__))
sensors_tools_csv = os.path.join(script_dir, 'sensors.csv')

manual_url = 'https://energietransitiewindesheim.nl/manuals'


def csv_create_update():
    """
    Add and update DeviceType and Property instances, using `sensors.csv`.

    The `sensors.csv` file is located in the same directory as this script.
    Existing DeviceType and Property instances, not described in the .csv,
    will never be deleted. If required, these should be removed manually
    (be careful: this may also remove Measurement instances).

    A special Property, 'heartbeat', always exists, and is defined for
    all DeviceType instances.

    This function is idempotent.
    """
    device_types: Dict[str, DeviceType] = {}
    properties: Dict[str, Property] = {}

    heartbeat = db.query(Property).filter(Property.name == 'heartbeat').first()
    if not heartbeat:
        heartbeat = Property(name='heartbeat')
        db.add(heartbeat)

    with open(sensors_tools_csv, newline='') as f:
        reader = csv.DictReader(f, delimiter=';')

        rows = [r for r in reader]

    for row in rows:
        device_type_name = row['DeviceType.name']
        device_type_display_name = row['DeviceType.DisplayName']
        property_name = row['Property.name']
        property_unit = row['Property.unit']

        assert ' ' not in device_type_name, 'DeviceType.name must not contain spaces'

        try:
            device_type = device_types[device_type_name]

        except KeyError:
            device_type = db.query(DeviceType).filter(DeviceType.name == device_type_name).first()
            if device_type:
                # Existing DeviceType; reset properties
                device_type.display_name = device_type_display_name
                device_type.properties = []
            else:
                # New DeviceType
                device_type = DeviceType(
                    name=device_type_name,
                    display_name=device_type_display_name,
                    installation_manual_url=f'{manual_url}/{device_type_name}.pdf'
                )
                db.add(device_type)

            # Every DeviceType has the heartbeat Property
            device_type.properties.append(heartbeat)
            device_types[device_type_name] = device_type

        try:
            property = properties[property_name]

        except KeyError:
            property = db.query(Property).filter(Property.name == property_name).first()
            if property:
                # Existing Property; rewrite unit
                property.unit = property_unit
            else:
                # New Property
                property = Property(
                    name=property_name,
                    unit=property_unit,
                )
                db.add(property)
                db.commit()
                db.refresh(property)

            properties[property_name] = property

        if property not in device_type.properties:
            device_type.properties.append(property)

    db.commit()
