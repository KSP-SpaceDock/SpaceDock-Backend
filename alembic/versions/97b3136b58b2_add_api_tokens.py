"""Add API tokens

Revision ID: 97b3136b58b2
Revises: 49cf6957c1f2
Create Date: 2016-10-21 15:14:11.381545

"""

# revision identifiers, used by Alembic.
revision = '97b3136b58b2'
down_revision = '49cf6957c1f2'
branch_labels = None
depends_on = None

from alembic import op
import sqlalchemy as sa


def upgrade():
    op.create_table('token', 
        sa.Column('meta', sa.String(512), server_default='{}'),
        sa.Column('id', sa.Integer, primary_key=True),
        sa.Column('token', sa.String(32), unique=True)
    )


def downgrade():
    op.drop_table('token')

