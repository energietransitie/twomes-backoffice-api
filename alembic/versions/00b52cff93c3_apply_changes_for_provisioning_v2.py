"""Apply changes for provisioning v2

Revision ID: 00b52cff93c3
Revises: 9443f1256ebf
Create Date: 2023-02-17 14:00:16.643789

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '00b52cff93c3'
down_revision = '9443f1256ebf'
branch_labels = None
depends_on = None


def upgrade():
    # New table for apps.
    op.create_table('app',
                    sa.Column('id', sa.Integer(), nullable=False),
                    sa.Column('name', sa.Text(), nullable=False),
                    sa.Column('provisioning_url_template', sa.Text(), nullable=False,
                              comment='URL template used to create invitation links'))
    # New table for campaigns.
    op.create_table('campaign',
                    sa.Column('id', sa.Integer(), nullable=False),
                    sa.Column('name'. sa.Text(), nullable=False),
                    sa.Column('app_id', sa.Integer(), nullable=False),
                    sa.Column('info_url', sa.Text(), nullable=True,
                              comment='URL to information about a campaign'),
                    sa.Column('start', sa.DateTime(), nullable=True,
                              comment='Start datetime of the campaign'),
                    sa.Column('end', sa.DateTime(), nullable=True,
                              comment='End datetime of the campaign'))
    # Relation to campaign.
    op.add_column('account',
                  sa.Column('campaign_id', sa.Integer(), nullable=True))
    # A device type also needs an info URL.
    op.add_column('device_type',
                  sa.Column('info_url', sa.Text(), nullable=True,
                            comment='URL to information about a device type'))
    # A device_type's display name is language dependent, so we remove it from the DB.
    op.drop_column('device_type', 'display_name')
    
    # Remove building characteristics
    op.drop_column('building', 'yr_built')
    op.drop_column('building', 'type')
    op.drop_column('building', 'floor_area')
    op.drop_column('building', 'heat_loss_area')
    op.drop_column('building', 'energy_label')
    op.drop_column('building', 'energy_index')


def downgrade():
    op.add_column('building', sa.Column('yr_built', sa.Text(), nullable=True, comment='Year built ("oorspronkelijk bouwjaar")'))
    op.add_column('building', sa.Column('type', sa.Text(), nullable=True, comment='House type ("woningtype")'))
    op.add_column('building', sa.Column('floor_area', sa.Integer(), nullable=True, comment='Floor area ("gebruiksoppervlakte"), as defined in NEN NTA 8800'))
    op.add_column('building', sa.Column('heat_loss_area', sa.Integer(), nullable=True, comment='Heat loss area ("verliesoppervlakte"), as defined in NEN NTA 8800'))
    op.add_column('building', sa.Column('energy_label', sa.Text(), nullable=True, comment='Energy label, as defined in NEN NTA 8800'))
    op.add_column('building', sa.Column('energy_index', sa.Float(), nullable=True, comment='Energy index, as defined in NEN NTA 8800'))
    
    op.add_column('device_type', sa.Column('display_name', sa.Text(), nullable=True,
                                           comment='Name to show in user interfaces'))
    op.drop_column('device_type', 'info_url')
    op.drop_column('account', 'campaign_id')
    op.drop_table('campaign')
    op.drop_table('app')