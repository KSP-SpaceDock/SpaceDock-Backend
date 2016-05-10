"""Add file size

Revision ID: a74df8caf629
Revises: f4c441491815
Create Date: 2016-05-11 06:14:35.511658

"""

# revision identifiers, used by Alembic.
revision = 'a74df8caf629'
down_revision = 'f4c441491815'
branch_labels = None
depends_on = None

from alembic import op
import sqlalchemy as sa


def upgrade():
    op.add_column('modversion', sa.Column('file_size', sa.Integer(), nullable=True))


def downgrade():
    pass
