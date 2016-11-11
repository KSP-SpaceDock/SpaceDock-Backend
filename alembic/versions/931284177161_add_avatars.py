"""Add avatars

Revision ID: 931284177161
Revises: 97b3136b58b2
Create Date: 2016-11-11 19:05:10.241353

"""

# revision identifiers, used by Alembic.
revision = '931284177161'
down_revision = '97b3136b58b2'
branch_labels = None
depends_on = None

from alembic import op
import sqlalchemy as sa


def upgrade():
    op.add_column('user', sa.Column('avatar', sa.String(512)))


def downgrade():
    op.drop_column('user', 'avatar')
