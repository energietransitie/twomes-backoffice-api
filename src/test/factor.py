from factory import post_generation, Sequence
from factory.alchemy import SQLAlchemyModelFactory

from test.common import session

import model


class DeviceTypeFactory(SQLAlchemyModelFactory):
    class Meta:
        model = model.DeviceType
        sqlalchemy_session = session

    name = Sequence(lambda n: f'device type {n}')
    installation_manual_url = Sequence(lambda n: f'https://manual.com/{n}')

    @post_generation
    def properties(self, create, extracted, **kwargs):
        if not create:
            # Simple build, do nothing
            return

        if extracted:
            # A list of properties were passed in, use them
            for property in extracted:
                self.properties.append(property)


class PropertyFactory(SQLAlchemyModelFactory):
    class Meta:
        model = model.Property
        sqlalchemy_session = session

    name = Sequence(lambda n: f'property {n}')
