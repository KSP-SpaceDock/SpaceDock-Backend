"""Remove KSP specific fields

Revision ID: 49cf6957c1f2
Revises: 81b6d4d30515
Create Date: 2016-10-21 15:13:44.766574

"""

# revision identifiers, used by Alembic.
revision = '49cf6957c1f2'
down_revision = '81b6d4d30515'
branch_labels = None
depends_on = None

from alembic import op
import sqlalchemy as sa


def upgrade():
    op.drop_column('user', 'dark_theme')
    op.drop_column('user', 'forumId')
    op.drop_column('mod', 'ckan')


def downgrade():
    op.add_column('user', sa.Column('dark_theme', sa.Boolean))
    op.add_column('user', sa.Column('forumId', sa.Integer))    
    op.add_column('mod', sa.Column('ckan', sa.Boolean))

