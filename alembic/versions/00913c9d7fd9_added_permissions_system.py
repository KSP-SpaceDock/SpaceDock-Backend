"""Added permissions system

Revision ID: 00913c9d7fd9
Revises: a66a267f85e3
Create Date: 2016-10-21 15:12:35.217723

>> This is only schema migration! <<
>> For data migration please execute tools/migrate_roles.py <<
"""

# revision identifiers, used by Alembic.
revision = '00913c9d7fd9'
down_revision = 'a66a267f85e3'
branch_labels = None
depends_on = None

from alembic import op
import sqlalchemy as sa


def upgrade():
    op.create_table('role', 
        sa.Column('meta', sa.String(512), server_default='{}'),
        sa.Column('id', sa.Integer(), primary_key=True),
        sa.Column('name', sa.String(120), unique=True),
        sa.Column('params',  sa.String(512))
    )
    op.create_table('ability', 
        sa.Column('meta', sa.String(512), server_default='{}'),
        sa.Column('id', sa.Integer, primary_key=True),
        sa.Column('name', sa.String(120), unique=True)
    )
    op.create_table('user_role',
        sa.Column('user_id', sa.Integer(), sa.ForeignKey('user.id'), primary_key=False),
        sa.Column('role_id', sa.Integer(), sa.ForeignKey('role.id'), primary_key=False)
    )
    op.create_table('role_ability',
        sa.Column('role_id', sa.Integer(), sa.ForeignKey('role.id'), primary_key=False),
        sa.Column('ability_id', sa.Integer(), sa.ForeignKey('ability.id'), primary_key=False)
    )
    op.drop_column('user', 'admin')


def downgrade():
    op.drop_table('role')
    op.drop_table('ability')
    op.drop_table('user_role')
    op.drop_table('role_ability')
    op.add_column('user', sa.Column('admin', sa.Boolean))
