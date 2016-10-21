"""Backporting from SpaceDock

Revision ID: 088e2f07069c
Revises: 
Create Date: 2016-10-21 15:00:32.828228

"""

# revision identifiers, used by Alembic.
revision = '088e2f07069c'
down_revision = None
branch_labels = None
depends_on = None

from alembic import op
import sqlalchemy as sa


def upgrade():
    op.add_column('gameversion', sa.Column('is_beta', sa.Boolean, nullable=True))
    op.add_column('mod', sa.Column('rating_count', sa.Integer, server_default=sa.text('0'), nullable=False))
    op.add_column('mod', sa.Column('total_score', sa.Float, nullable=True))
    op.add_column('modversion', sa.Column('is_beta', sa.Boolean, nullable=True))file_size
    op.add_column('modversion', sa.Column('file_size', sa.Integer))
    op.add_column('user', sa.Column('facebookUsername', sa.String(length=128), nullable=True))
    op.add_column('user', sa.Column('showCreated', sa.Boolean, nullable=True))
    op.add_column('user', sa.Column('showEmail', sa.Boolean, nullable=True))
    op.add_column('user', sa.Column('showFacebookName', sa.Boolean, nullable=True))
    op.add_column('user', sa.Column('showForumName', sa.Boolean, nullable=True))
    op.add_column('user', sa.Column('showIRCName', sa.Boolean, nullable=True))
    op.add_column('user', sa.Column('showLocation', sa.Boolean, nullable=True))
    op.add_column('user', sa.Column('showRedditName', sa.Boolean, nullable=True))
    op.add_column('user', sa.Column('showTwitchName', sa.Boolean, nullable=True))
    op.add_column('user', sa.Column('showTwitterName', sa.Boolean, nullable=True))
    op.add_column('user', sa.Column('showYoutubeName', sa.Boolean, nullable=True))
    op.add_column('user', sa.Column('twitchUsername', sa.String(128), nullable=True))
    op.add_column('user', sa.Column('youtubeUsername', sa.String(128), nullable=True))
    op.create_table('ratings',
        sa.Column('id', sa.Integer, primary_key=True),
        sa.Column('user_id', sa.Integer, sa.ForeignKey('user.id')),
        sa.Column('mod_id', sa.Integer, sa.ForeignKey('mod.id')),
        sa.Column('score', sa.Float, nullable=False, server_default=sa.text('5')),
        sa.Column('created', sa.DateTime),
        sa.Column('updated', sa.DateTime)
    )
    op.create_table('review',
        sa.Column('id', sa.Integer, primary_key = True),
        sa.Column('user_id', sa.Integer, sa.ForeignKey('user.id')),
        sa.Column('mod_id', sa.Integer, sa.ForeignKey('mod.id')),
        sa.Column('review_title', sa.String(100), index = True),
        sa.Column('review_text', sa.Unicode(100000)),
        sa.Column('video_link', sa.String(100)),
        sa.Column('video_image', sa.String(100)),
        sa.Column('has_video', sa.Boolean),
        sa.Column('teaser', sa.Unicode(1000)),
        sa.Column('approved', sa.Boolean),
        sa.Column('published', sa.Boolean),
        sa.Column('created', sa.DateTime),
        sa.Column('updated', sa.DateTime)
    )
    op.create_table('reviewmedia',
        sa.Column('id', sa.Integer, primary_key = True),
        sa.Column('review_id', sa.Integer, sa.ForeignKey('review.id')),
        sa.Column('hash', sa.String(12)),
        sa.Column('type', sa.String(32)),
        sa.Column('data', sa.String(512))
    )

def downgrade():
    op.drop_column('user', 'youtubeUsername')
    op.drop_column('user', 'twitchUsername')
    op.drop_column('user', 'showYoutubeName')
    op.drop_column('user', 'showTwitterName')
    op.drop_column('user', 'showTwitchName')
    op.drop_column('user', 'showRedditName')
    op.drop_column('user', 'showLocation')
    op.drop_column('user', 'showIRCName')
    op.drop_column('user', 'showForumName')
    op.drop_column('user', 'showFacebookName')
    op.drop_column('user', 'showEmail')
    op.drop_column('user', 'showCreated')
    op.drop_column('user', 'facebookUsername')
    op.drop_column('modversion', 'is_beta')
    op.drop_column('modversion', 'file_size')
    op.drop_column('mod', 'total_score')
    op.drop_column('mod', 'rating_count')
    op.drop_column('gameversion', 'is_beta')
    op.drop_table('ratings')
    op.drop_table('review')
    op.drop_table('reviewmedia')
