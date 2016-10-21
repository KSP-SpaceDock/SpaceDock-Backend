"""Added meta fields

Revision ID: a66a267f85e3
Revises: 088e2f07069c
Create Date: 2016-10-21 15:12:00.021483

"""

# revision identifiers, used by Alembic.
revision = 'a66a267f85e3'
down_revision = '088e2f07069c'
branch_labels = None
depends_on = None

from alembic import op
import sqlalchemy as sa


def upgrade():
    op.add_column('featured', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('blog', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('user', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('user_auth', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('ratings', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('review', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('publisher', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('game', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('mod', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('modlist', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('modlistitem', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('sharedauthor', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('downloadevent', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('followevent', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('referralevent', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('modversion', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('media', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('reviewmedia', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('gameversion', sa.Column('meta', sa.String(512), server_default='{}'))


def downgrade():
    op.drop_column('featured', 'meta')
    op.drop_column('blog', 'meta')
    op.drop_column('user', 'meta')
    op.drop_column('user_auth', 'meta')
    op.drop_column('rating', 'meta')
    op.drop_column('review', 'meta')
    op.drop_column('publisher', 'meta')
    op.drop_column('game', 'meta')
    op.drop_column('mod', 'meta')
    op.drop_column('modlist', 'meta')
    op.drop_column('modlistitem', 'meta')
    op.drop_column('sharedauthor', 'meta')
    op.drop_column('downloadevent', 'meta')
    op.drop_column('followevent', 'meta')
    op.drop_column('referralevent', 'meta')
    op.drop_column('modversion', 'meta')
    op.drop_column('media', 'meta')
    op.drop_column('reviewmedia', 'meta')
    op.drop_column('gameversion', 'meta')
