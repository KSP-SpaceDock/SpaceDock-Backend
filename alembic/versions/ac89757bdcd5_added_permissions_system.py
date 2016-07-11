"""Added Permissions System

Revision ID: ac89757bdcd5
Revises: f4c441491815
Create Date: 2016-07-10 07:07:56.174978

"""

# revision identifiers, used by Alembic.
revision = 'ac89757bdcd5'
down_revision = 'f4c441491815'
branch_labels = None
depends_on = None

from alembic import op
import sqlalchemy as sa


def upgrade():
    op.add_column('users', sa.Column('type', sa.String(50)))
    op.add_column('users', sa.Column('params', sa.String(512)))


def downgrade():
    op.drop_column('users', 'type')
    op.drop_column('users', 'params')
