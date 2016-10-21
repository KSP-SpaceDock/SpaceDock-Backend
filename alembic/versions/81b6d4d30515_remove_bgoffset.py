"""Remove bgOffset

Revision ID: 81b6d4d30515
Revises: 00913c9d7fd9
Create Date: 2016-10-21 15:13:00.686561

"""

# revision identifiers, used by Alembic.
revision = '81b6d4d30515'
down_revision = '00913c9d7fd9'
branch_labels = None
depends_on = None

from alembic import op
import sqlalchemy as sa


def upgrade():
    op.drop_column('user', 'bgOffsetX')
    op.drop_column('user', 'bgOffsetY')
    op.drop_column('publisher', 'bgOffsetX')
    op.drop_column('publisher', 'bgOffsetY')
    op.drop_column('game', 'bgOffsetX')
    op.drop_column('game', 'bgOffsetY')
    op.drop_column('mod', 'bgOffsetX')
    op.drop_column('mod', 'bgOffsetY')
    op.drop_column('modlist', 'bgOffsetY')


def downgrade():
    op.add_column('user', sa.Column('bgOffsetX', sa.Integer))
    op.add_column('user', sa.Column('bgOffsetY', sa.Integer))
    op.add_column('publisher', sa.Column('bgOffsetX', sa.Integer))
    op.add_column('publisher', sa.Column('bgOffsetY', sa.Integer))
    op.add_column('game', sa.Column('bgOffsetX', sa.Integer))
    op.add_column('game', sa.Column('bgOffsetY', sa.Integer))
    op.add_column('mod', sa.Column('bgOffsetX', sa.Integer))
    op.add_column('mod', sa.Column('bgOffsetY', sa.Integer))
    op.add_column('modlist', sa.Column('bgOffsetY', sa.Integer))
