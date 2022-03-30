from test.factor import DeviceTypeFactory, PropertyFactory, session


def base():
    manual_url = 'https://energiebeveiliging.nl/manuals/'

    heartbeat = PropertyFactory.create(
        name='heartbeat',
    )
    room_temperature = PropertyFactory.create(
        name='room temperature',
        unit='°C'
    )
    atmospheric_pressure = PropertyFactory.create(
        name='atmospheric pressure',
        unit='hPa'
    )

    DeviceTypeFactory.create(
        name='Gateway',
        installation_manual_url=manual_url + 'gateway.pdf',
        properties=(heartbeat,)
    )
    DeviceTypeFactory.create(
        name='Smart Meter',
        installation_manual_url=manual_url + 'smart-meter.pdf',
        properties=(heartbeat, room_temperature)
    )
    DeviceTypeFactory.create(
        name='Smart Thermostat',
        installation_manual_url=manual_url + 'smart-thermostat.pdf',
        properties=(heartbeat, room_temperature, atmospheric_pressure)
    )

    session.commit()
