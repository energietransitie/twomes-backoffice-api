"""Rename proof_of_presence_id to activation_token, for Device

Revision ID: 9443f1256ebf
Revises: 6e528df2f1a5
Create Date: 2021-06-06 16:53:06.935359

"""
from alembic import op
from sqlalchemy.sql.expression import text
import sqlalchemy as sa
from sqlalchemy.dialects import mysql

# revision identifiers, used by Alembic.
revision = '9443f1256ebf'
down_revision = '6e528df2f1a5'
branch_labels = None
depends_on = None


def upgrade():
    op.alter_column('device', 'proof_of_presence_id',
               new_column_name='activation_token',
               existing_type=mysql.TEXT(),
               existing_comment='Unique, random token to identify the device during activation',
               existing_nullable=False)
    op.alter_column('device', 'activated_on',
               existing_type=mysql.DATETIME(),
               comment='Time at which the activation token is used to activate the device',
               existing_comment='Time at which the proof-of-presence id is used to active the device',
               existing_nullable=True)


def downgrade():
    op.alter_column('device', 'activation_token',
               new_column_name='proof_of_presence_id',
               existing_type=mysql.TEXT(),
               existing_comment='Unique, random token to identify the device during activation',
               existing_nullable=False)
    op.alter_column('device', 'activated_on',
               existing_type=mysql.DATETIME(),
               comment='Time at which the proof-of-presence id is used to active the device',
               existing_comment='Time at which the activation token is used to activate the device',
               existing_nullable=True)
