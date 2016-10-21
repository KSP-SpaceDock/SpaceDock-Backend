"""SpaceDock Backend Migration

Revision ID: 60264098e5f0
Revises: a74df8caf629
Create Date: 2016-08-27 10:41:08.171557

"""

# revision identifiers, used by Alembic.
revision = '60264098e5f0'
down_revision = 'a74df8caf629'
branch_labels = None
depends_on = None

from alembic import op
import sqlalchemy as sa


def upgrade():
    
    # Featured
    op.add_column('featured', sa.Column('meta', sa.String(512), server_default='{}'))
    
    # Blog Post
    op.add_column('blog', sa.Column('meta', sa.String(512), server_default='{}'))
    
    # User
    op.drop_column('user', 'bgOffsetX')
    op.drop_column('user', 'bgOffsetY')
    op.drop_column('user', 'dark_theme')
    op.drop_column('user', 'admin')
    op.drop_column('user', 'forumId')
    op.add_column('user', sa.Column('meta', sa.String(512), server_default='{}'))
    
    # Roles
    op.create_table('role', 
        sa.Column('meta', sa.String(512), server_default='{}'),
        sa.Column('id', sa.Integer(), primary_key=True),
        sa.Column('name', sa.String(120), unique=True),
        sa.Column('params',  sa.String(512))
    )
    op.create_table('user_role',
        sa.Column('user_id', sa.Integer(), sa.ForeignKey('user.id'), primary_key=False),
        sa.Column('role_id', sa.Integer(), sa.ForeignKey('role.id'), primary_key=False)
    )
    op.create_table('role_ability',
        sa.Column('role_id', sa.Integer(), sa.ForeignKey('role.id'), primary_key=False),
        sa.Column('ability_id', sa.Integer(), sa.ForeignKey('ability.id'), primary_key=False)
    )
    
    # Abilities
    op.create_table('ability', 
        sa.Column('meta', sa.String(512), server_default='{}'),
        sa.Column('id', sa.Integer, primary_key=True),
        sa.Column('name', sa.String(120), unique=True)
    )
    
    # User Auth
    op.add_column('user_auth', sa.Column('meta', sa.String(512), server_default='{}'))
    
    # Ratings
    op.add_column('rating', sa.Column('meta', sa.String(512), server_default='{}'))
    
    # Reviews
    op.add_column('review', sa.Column('meta', sa.String(512), server_default='{}'))
    
    # Publisher    
    op.drop_column('publisher', 'bgOffsetX')
    op.drop_column('publisher', 'bgOffsetY')
    op.add_column('publisher', sa.Column('meta', sa.String(512), server_default='{}'))
    
    # Game
    op.drop_column('game', 'bgOffsetX')
    op.drop_column('game', 'bgOffsetY')
    op.add_column('game', sa.Column('meta', sa.String(512), server_default='{}'))
    
    # Mods
    op.drop_column('mod', 'bgOffsetX')
    op.drop_column('mod', 'bgOffsetY')
    op.drop_column('mod', 'ckan')
    op.add_column('mod', sa.Column('meta', sa.String(512), server_default='{}'))
    
    # Modlist
    op.drop_column('modlist', 'bgOffsetY')
    op.add_column('modlist', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('modlistitem', sa.Column('meta', sa.String(512), server_default='{}'))
    
    # Shared Author    
    op.add_column('sharedauthor', sa.Column('meta', sa.String(512), server_default='{}'))
    
    # Events
    op.add_column('downloadevent', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('followevent', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('referralevent', sa.Column('meta', sa.String(512), server_default='{}'))
    
    # Modversion    
    op.add_column('modversion', sa.Column('meta', sa.String(512), server_default='{}'))
    
    # Media
    op.add_column('media', sa.Column('meta', sa.String(512), server_default='{}'))
    op.add_column('reviewmedia', sa.Column('meta', sa.String(512), server_default='{}'))
    
    # GameVersion
    op.add_column('gameversion', sa.Column('meta', sa.String(512), server_default='{}'))
       
    # Tokens
    op.create_table('token', 
        sa.Column('meta', sa.String(512), server_default='{}'),
        sa.Column('id', sa.Integer, primary_key=True),
        sa.Column('token', sa.String(32), unique=True)
    )


def downgrade():
    
    # Featured
    op.drop_column('featured', 'meta')
    
    # Blog Post
    op.drop_column('blog', 'meta')
    
    # User
    op.add_column('user', sa.Column('bgOffsetX', sa.Integer))
    op.add_column('user', sa.Column('bgOffsetY', sa.Integer))
    op.add_column('user', sa.Column('dark_theme', sa.Boolean))
    op.add_column('user', sa.Column('forumId', sa.Integer))
    op.add_column('user', sa.Column('admin', sa.Boolean))
    op.drop_column('user', 'meta')
    
    # Roles
    op.drop_table('role')
    op.drop_table('user_role')
    op.drop_table('role_ability')
    
    # Abilities
    op.drop_table('ability')
    
    # User Auth
    op.drop_column('user_auth', 'meta')
    
    # Ratings
    op.drop_column('rating', 'meta')
    
    # Reviews
    op.drop_column('review', 'meta')
    
    # Publisher    
    op.add_column('publisher', sa.Column('bgOffsetX', sa.Integer))
    op.add_column('publisher', sa.Column('bgOffsetY', sa.Integer))
    op.drop_column('publisher', 'meta')
    
    # Game
    op.add_column('game', sa.Column('bgOffsetX', sa.Integer))
    op.add_column('game', sa.Column('bgOffsetY', sa.Integer))
    op.drop_column('game', 'meta')
    
    # Mods
    op.add_column('mod', sa.Column('bgOffsetX', sa.Integer))
    op.add_column('mod', sa.Column('bgOffsetY', sa.Integer))
    op.add_column('mod', sa.Column('ckan', sa.Boolean))
    op.drop_column('mod', 'meta')
    
    # Modlist
    op.add_column('modlist', sa.Column('bgOffsetY', sa.Integer))
    op.drop_column('modlist', 'meta')
    op.drop_column('modlistitem', 'meta')
    
    # Shared Author    
    op.drop_column('sharedauthor', 'meta')
    
    # Events
    op.drop_column('downloadevent', 'meta')
    op.drop_column('followevent', 'meta')
    op.drop_column('referralevent', 'meta')
    
    # Modversion    
    op.drop_column('modversion', 'meta')
    
    # Media
    op.drop_column('media', 'meta')
    op.drop_column('reviewmedia', 'meta')
    
    # GameVersion
    op.drop_column('gameversion', 'meta')

    # Tokens
    op.drop_table('token')
