from datetime import datetime
from sqlalchemy import Column, Integer, String, Unicode, Boolean, DateTime, ForeignKey, Table, UnicodeText, Text, text,Float
from sqlalchemy.orm import relationship, backref
from sqlalchemy.ext.associationproxy import association_proxy
from SpaceDock.database import Base, MetaObject, db
from SpaceDock.config import cfg

import os.path
import bcrypt
import json
import re

mod_followers = Table('mod_followers', Base.metadata,
    Column('mod_id', Integer, ForeignKey('mod.id')),
    Column('user_id', Integer, ForeignKey('user.id')),
)

user_role_table = Table('user_role', Base.metadata,
    Column('user_id', Integer(), ForeignKey('user.id'), primary_key=False),
    Column('role_id', Integer(), ForeignKey('role.id'), primary_key=False)
)

role_ability_table = Table('role_ability', Base.metadata,
    Column('role_id', Integer(), ForeignKey('role.id'), primary_key=False),
    Column('ability_id', Integer(), ForeignKey('ability.id'), primary_key=False)
)

def role_find_or_create(r):
    role = Role.query.filter_by(name=r).first()
    if not role:
        role = Role(name=r)
        db.add(role)
        db.commit()
    return role

def is_sequence(arg):
    return (not hasattr(arg, "strip") and
            hasattr(arg, "__getitem__") or
            hasattr(arg, "__iter__"))

class Featured(Base, MetaObject):
    __tablename__ = 'featured'
    id = Column(Integer, primary_key = True)
    mod_id = Column(Integer, ForeignKey('mod.id'))
    mod = relationship('Mod', backref=backref('mod', order_by=id))
    created = Column(DateTime)

    def __init__(self, modid):
        self.created = datetime.now()
        self.mod_id = modid


    def __repr__(self):
        return '<Featured %r>' % self.id


class BlogPost(Base, MetaObject):
    __tablename__ = 'blog'
    id = Column(Integer, primary_key = True)
    title = Column(Unicode(1024))
    text = Column(Unicode(65535))
    created = Column(DateTime)

    def __init__(self):
        self.created = datetime.now()

    def __repr__(self):
        return '<Blog Post %r>' % self.id


class User(Base, MetaObject):
    __tablename__ = 'user'
    __lock__ = ['id', 'username', 'password', 'created', 'confirmation', 'passwordReset', 'passwordResetExpiry', 'backgroundMedia', 'ratings', 'review', 'mods', 'packs', 'following', '_roles', 'roles']
    id = Column(Integer, primary_key = True)
    username = Column(String(128), nullable = False, index = True)
    email = Column(String(256), nullable = False, index = True)
    showEmail = Column(Boolean())
    public = Column(Boolean())
    password = Column(String(128))
    description = Column(Unicode(10000))
    created = Column(DateTime)
    showCreated = Column(Boolean())
    forumUsername = Column(String(128))
    showForumName = Column(Boolean())
    forumId = Column(Integer)
    ircNick = Column(String(128))
    showIRCName = Column(Boolean())
    twitterUsername = Column(String(128))
    showTwitterName = Column(Boolean())
    redditUsername = Column(String(128))
    showRedditName = Column(Boolean())
    youtubeUsername = Column(String(128))
    showYoutubeName = Column(Boolean())
    twitchUsername = Column(String(128))
    showTwitchName = Column(Boolean())
    facebookUsername = Column(String(128))
    showFacebookName = Column(Boolean())
    location = Column(String(128))
    showLocation = Column(Boolean())
    confirmation = Column(String(128))
    passwordReset = Column(String(128))
    passwordResetExpiry = Column(DateTime)
    backgroundMedia = Column(String(512))
    ratings = relationship('Rating', order_by='Rating.created')
    review = relationship('Review', order_by='Review.created')
    mods = relationship('Mod', order_by='Mod.created')
    packs = relationship('ModList', order_by='ModList.created')
    following = relationship('Mod', secondary=mod_followers, backref='user.id')
    theme = Column(String(24))
    # Permissions
    _roles = relationship('Role', secondary=user_role_table, backref='users')
    roles = association_proxy('_roles', 'name', creator=role_find_or_create)

    def set_password(self, password):
        self.password = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt()).decode('utf-8')

    def __init__(self, username, email, password, roles=None, default_role='unconfirmed'):
        self.email = email
        self.username = username
        self.password = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt()).decode('utf-8')
        self.public = False
        self.admin = False
        self.created = datetime.now()
        self.youtubeUsername = ''
        self.twitchUsername = ''
        self.facebookUsername = ''
        self.twitterUsername = ''
        self.forumUsername = ''
        self.ircNick = ''
        self.description = ''
        self.backgroundMedia = ''
        self.dark_theme = False
        if roles and isinstance(roles, basestring):
            roles = [roles]
        if roles and is_sequence(roles):
            self.roles = roles
        elif default_role:
            self.roles = [default_role]
    def __repr__(self):
        return '<User %r>' % self.username

    # Flask.Login stuff
    # We don't use most of these features
    def is_authenticated(self):
        return True
    def is_active(self):
        return confirmation == None
    def is_anonymous(self):
        return False
    def get_id(self):
        return self.username

    # Permissions
    def add_roles(self, *roles):
        for r in roles:
            if not r in self.roles:
                self.roles.append(r)
    def remove_roles(self, *roles):
        for r in roles:
            if r in self.roles:
                self.roles.remove(r)


class Role(Base, MetaObject):
    __tablename__ = 'role'
    id = Column(Integer(), primary_key=True)
    name = Column(String(120), unique=True)
    abilities = relationship('Ability', secondary=role_ability_table, backref='roles')
    params = Column(String(512))

    def __init__(self, name):
        self.name = name
        self.params = '{}'

    def add_abilities(self, *abilities):
        for ability in abilities:
            existing_ability = Ability.query.filter(Ability.name == ability).first()
            if not existing_ability:
                existing_ability = Ability(ability)
                db.add(existing_ability)
                db.commit()
            self.abilities.append(existing_ability)

    def add_abilities_re(self, pattern):
        for ability in Ability.query.all():
            if re.match(pattern, ability.name) and not ability in self.abilities:
                self.abilities.append(ability)

    def remove_abilities(self, *abilities):
        for ability in abilities:
            existing_ability = Ability.query.filter(Ability.name == ability).first()
            if existing_ability and existing_ability in self.abilities:
                self.abilities.remove(existing_ability)

    def get_param(self, ability, param):
        p = json.loads(self.params)
        if ability in p.keys():
            if param in p[ability].keys():
                return p[ability][param]
        return None
    def add_param(self, ability, param, value):
        p = json.loads(self.params)
        if not ability in p:
            p[ability] = dict()
        if not param in p[ability]:
            p[ability][param] = list()
        if not value in p[ability][param]:
            p[ability][param].append(value)
        self.params = json.dumps(p)
    def remove_param(self, ability, param, value):
        p = json.loads(self.params)
        if ability in p:
            if param in p[ability]:
                if value in p[ability][param]:
                    p[ability][param].remove(value)
        self.params = json.dumps(p)

    def __repr__(self):
        return '<Role {}>'.format(self.name)

    def __str__(self):
        return self.name


class Ability(Base, MetaObject):
    __tablename__ = 'ability'
    id = Column(Integer(), primary_key=True)
    name = Column(String(120), unique=True)

    def __init__(self, name):
        self.name = name

    def __repr__(self):
        return '<Ability {}>'.format(self.name)

    def __str__(self):
        return self.name


class UserAuth(Base, MetaObject):
    __tablename__ = 'user_auth'
    id = Column(Integer, primary_key=True)
    user_id = Column(Integer, nullable=False, index=True)
    provider = Column(String(32))  # 'github' or 'google', etc.
    remote_user = Column(String(128), index=True)  # Usually the username on the other side
    created = Column(DateTime)
    # We can keep a token here, to allow interacting with the provider's API
    # on behalf of the user.

    def __init__(self, user_id, remote_user, provider):
        self.user_id = user_id
        self.provider = provider
        self.remote_user = remote_user
        self.created = datetime.now()

    def __repr__(self):
        return '<UserAuth %r, User %r>' % (self.provider, self.user_id)


class Rating(Base, MetaObject):
    __tablename__ = 'ratings'
    id = Column(Integer, primary_key = True)
    user_id = Column(Integer, ForeignKey('user.id'))
    user = relationship('User', back_populates='ratings')
    mod_id = Column(Integer, ForeignKey('mod.id'))
    mod = relationship('Mod', back_populates='ratings')
    score = Column(Float(), nullable=False, server_default=text('5'))
    created = Column(DateTime)
    updated = Column(DateTime)

    def __init__(self, user_id, user, mod_id, mod, score):
        from SpaceDock.common import clamp_number
        self.user_id = user_id
        self.user = user
        self.mod_id = mod_id
        self.mod = mod
        self.score = clamp_number(0, 5, score)
        self.created = datetime.now()
        self.updated = datetime.now()

    def __repr__(self):
        return '<Rating %r %r>' % (self.id, self.score)


class Review(Base, MetaObject):
    __tablename__ = 'review'
    id = Column(Integer, primary_key = True)
    user_id = Column(Integer, ForeignKey('user.id'))
    user = relationship('User', back_populates='review')
    mod_id = Column(Integer, ForeignKey('mod.id'))
    mod = relationship('Mod', back_populates='review')
    review_title = Column(String(100), index = True)
    review_text = Column(Unicode(100000))
    medias = relationship('ReviewMedia')
    video_link = Column(String(100))
    video_image = Column(String(100))
    has_video = Column(Boolean())
    teaser = Column(Unicode(1000))
    approved = Column(Boolean())
    published = Column(Boolean())
    created = Column(DateTime)
    updated = Column(DateTime)

    def __init__(self):
        self.created = datetime.now()
        self.updated = datetime.now()

    def __repr__(self):
        return '<Review %r %r>' % (self.id, self.review_title)


class Publisher(Base, MetaObject):
    __tablename__ = 'publisher'
    __lock__ = ['id', 'created', 'updated', 'background', 'games']
    id = Column(Integer, primary_key = True)
    name = Column(Unicode(1024))
    short_description = Column(Unicode(1000))
    description = Column(Unicode(100000))
    created = Column(DateTime)
    updated = Column(DateTime)
    background = Column(String(512))
    link = Column(Unicode(1024))
    games = relationship('Game', back_populates='publisher')

    def __init__(self,name):
        self.created = datetime.now()
        self.name = name

    def __repr__(self):
        return '<Publisher %r %r>' % (self.id, self.name)


class Game(Base, MetaObject):
    __tablename__ = 'game'
    __lock__ = ['id', 'rating', 'short', 'publisher_id', 'publisher', 'mods', 'modlists' 'version', 'created', 'updated', 'background']
    id = Column(Integer, primary_key = True)
    name = Column(Unicode(1024))
    active = Column(Boolean())
    fileformats = Column(Unicode(1024))
    altname = Column(Unicode(1024))
    rating = Column(Float())
    releasedate = Column(DateTime)
    short = Column(Unicode(1024))
    publisher_id = Column(Integer, ForeignKey('publisher.id'))
    publisher = relationship('Publisher', back_populates='games')
    description = Column(Unicode(100000))
    short_description = Column(Unicode(1000))
    created = Column(DateTime)
    updated = Column(DateTime)
    background = Column(String(512))
    link = Column(Unicode(1024))
    mods = relationship('Mod', back_populates='game')
    modlists = relationship('ModList', back_populates='game')
    version = relationship('GameVersion', back_populates='game')

    def background_thumb(self):
        if (cfg['thumbnail_size'] == ''):
            return self.background
        thumbnailSizesStr = cfg['thumbnail_size'].split('x')
        thumbnailSize = (int(thumbnailSizesStr[0]), int(thumbnailSizesStr[1]))
        split = os.path.split(self.background)
        thumbPath = os.path.join(split[0], 'thumb_' + split[1])
        fullThumbPath = os.path.join(os.path.join(cfg['storage'], thumbPath.replace('/content/', '')))
        fullImagePath = os.path.join(cfg['storage'], self.background.replace('/content/', ''))
        if not os.path.exists(fullThumbPath):
            thumbnail.create(fullImagePath, fullThumbPath, thumbnailSize)
        return thumbPath

    def __init__(self,name,publisher_id,short):
        self.created = datetime.now()
        self.name = name
        self.publisher_id = publisher_id
        self.short = short
        self.updated = datetime.now()

    def __repr__(self):
        return '<Game %r %r>' % (self.id, self.name)


class Mod(Base, MetaObject):
    __tablename__ = 'mod'
    __lock__ = ['id', 'user_id', 'user', 'game_id', 'game', 'shared_authors', 'approved', 'votes', 'created', 'updated', 'background', 'medias', 'versions', 
                'downloads', 'follow_events', 'referrals', 'followers', 'rating', 'review', 'follower_count', 'download_count', 'total_score', 'rating_count']
    id = Column(Integer, primary_key = True)
    user_id = Column(Integer, ForeignKey('user.id'))
    user = relationship('User', backref=backref('mod', order_by=id))
    game_id = Column(Integer, ForeignKey('game.id'))
    game = relationship('Game', back_populates='mods')
    shared_authors = relationship('SharedAuthor')
    name = Column(String(100), index = True)
    description = Column(Unicode(100000))
    short_description = Column(Unicode(1000))
    approved = Column(Boolean())
    published = Column(Boolean())
    donation_link = Column(String(512))
    external_link = Column(String(512))
    license = Column(String(128))
    votes = Column(Integer())
    created = Column(DateTime)
    updated = Column(DateTime)
    background = Column(String(512))
    medias = relationship('Media')
    default_version_id = Column(Integer)
    versions = relationship('ModVersion', order_by="desc(ModVersion.sort_index)")
    downloads = relationship('DownloadEvent', order_by="desc(DownloadEvent.created)")
    follow_events = relationship('FollowEvent', order_by="desc(FollowEvent.created)")
    referrals = relationship('ReferralEvent', order_by="desc(ReferralEvent.created)")
    source_link = Column(String(256))
    follower_count = Column(Integer, nullable=False, server_default=text('0'))
    download_count = Column(Integer, nullable=False, server_default=text('0'))
    followers = relationship('User', viewonly=True, secondary=mod_followers, backref='mod.id')
    ratings = relationship('Rating', order_by='Rating.created')
    review = relationship('Review', order_by='Review.created')
    total_score = Column(Float(), nullable=True)
    rating_count = Column(Integer, nullable=False, server_default=text('0'))

    def background_thumb(self):
        if (cfg['thumbnail_size'] == ''):
            return self.background
        thumbnailSizesStr = cfg['thumbnail_size'].split('x')
        thumbnailSize = (int(thumbnailSizesStr[0]), int(thumbnailSizesStr[1]))
        split = os.path.split(self.background)
        thumbPath = os.path.join(split[0], 'thumb_' + split[1])
        fullThumbPath = os.path.join(os.path.join(cfg['storage'], thumbPath.replace('/content/', '')))
        fullImagePath = os.path.join(cfg['storage'], self.background.replace('/content/', ''))
        if not os.path.exists(fullThumbPath):
            thumbnail.create(fullImagePath, fullThumbPath, thumbnailSize)
        return thumbPath

    def default_version(self):
        versions = [v for v in self.versions if v.id == self.default_version_id]
        if len(versions) == 0:
            return None
        return versions[0]

    def __init__(self, name, user_id, game_id, license):
        self.name = name
        self.user_id = user_id
        self.game_id = game_id
        self.license = license
        self.created = datetime.now()
        self.updated = datetime.now()
        self.approved = False
        self.published = False
        self.votes = 0
        self.follower_count = 0
        self.download_count = 0
        self.user = User.query.filter(User.id == user_id).first()
        self.game = Game.query.filter(Game.id == game_id).first()

    def __repr__(self):
        return '<Mod %r %r>' % (self.id, self.name)


class ModList(Base, MetaObject):
    __tablename__ = 'modlist'
    __lock__ = ['id', 'user_id', 'user', 'created', 'game_id', 'game', 'background', 'mods']
    id = Column(Integer, primary_key = True)
    user = relationship('User', backref=backref('modlist', order_by=id))
    user_id = Column(Integer, ForeignKey('user.id'))
    created = Column(DateTime)
    game_id = Column(Integer, ForeignKey('game.id'))
    game = relationship('Game', back_populates='modlists')
    background = Column(String(32))
    description = Column(Unicode(100000))
    short_description = Column(Unicode(1000))
    name = Column(Unicode(1024))
    mods = relationship('ModListItem', order_by="asc(ModListItem.sort_index)")

    def __init__(self):
        self.created = datetime.now()

    def __repr__(self):
        return '<ModList %r %r>' % (self.id, self.name)


class ModListItem(Base, MetaObject):
    __tablename__ = 'modlistitem'
    __lock__ = ['id', 'mod_id', 'mod', 'mod_list_id', 'mod_list']
    id = Column(Integer, primary_key = True)
    mod_id = Column(Integer, ForeignKey('mod.id'))
    mod = relationship('Mod', viewonly=True, backref=backref('modlistitem'))
    mod_list_id = Column(Integer, ForeignKey('modlist.id'))
    mod_list = relationship('ModList', viewonly=True, backref=backref('modlistitem'))
    sort_index = Column(Integer)

    def __init__(self):
        self.sort_index = 0

    def __repr__(self):
        return '<ModListItem %r %r>' % (self.mod_id, self.mod_list_id)


class SharedAuthor(Base, MetaObject):
    __tablename__ = 'sharedauthor'
    id = Column(Integer, primary_key = True)
    mod_id = Column(Integer, ForeignKey('mod.id'))
    mod = relationship('Mod', viewonly=True, backref=backref('sharedauthor'))
    user_id = Column(Integer, ForeignKey('user.id'))
    user = relationship('User', backref=backref('sharedauthor', order_by=id))
    accepted = Column(Boolean)

    def __init__(self):
        self.accepted = False

    def __repr__(self):
        return '<SharedAuthor %r>' % self.user_id


class DownloadEvent(Base, MetaObject):
    __tablename__ = 'downloadevent'
    id = Column(Integer, primary_key = True)
    mod_id = Column(Integer, ForeignKey('mod.id'))
    mod = relationship('Mod', viewonly=True, backref=backref('downloadevent', order_by="desc(DownloadEvent.created)"))
    version_id = Column(Integer, ForeignKey('modversion.id'))
    version = relationship('ModVersion', backref=backref('downloadevent', order_by="desc(DownloadEvent.created)"))
    downloads = Column(Integer)
    created = Column(DateTime)

    def __init__(self):
        self.downloads = 0
        self.created = datetime.now()

    def __repr__(self):
        return '<Download Event %r>' % self.id


class FollowEvent(Base, MetaObject):
    __tablename__ = 'followevent'
    id = Column(Integer, primary_key = True)
    mod_id = Column(Integer, ForeignKey('mod.id'))
    mod = relationship('Mod', viewonly=True, backref=backref('followevent', order_by="desc(FollowEvent.created)"))
    events = Column(Integer)
    delta = Column(Integer)
    created = Column(DateTime)

    def __init__(self):
        self.delta = 0
        self.created = datetime.now()

    def __repr__(self):
        return '<Download Event %r>' % self.id


class ReferralEvent(Base, MetaObject):
    __tablename__ = 'referralevent'
    id = Column(Integer, primary_key = True)
    mod_id = Column(Integer, ForeignKey('mod.id'))
    mod = relationship('Mod', viewonly=True, backref=backref('referralevent', order_by="desc(ReferralEvent.created)"))
    host = Column(String(128))
    events = Column(Integer)
    created = Column(DateTime)

    def __init__(self):
        self.events = 0
        self.created = datetime.now()

    def __repr__(self):
        return '<Download Event %r>' % self.id


class ModVersion(Base, MetaObject):
    __tablename__ = 'modversion'
    __lock__ = ['id', 'mod_id', 'mod', 'friendly_version', 'gameversion_id', 'gameversion', 'created', 'download_path', 'file_size']
    id = Column(Integer, primary_key = True)
    mod_id = Column(Integer, ForeignKey('mod.id'))
    mod = relationship('Mod', viewonly=True, backref=backref('modversion', order_by="desc(ModVersion.created)"))
    friendly_version = Column(String(64))
    is_beta = Column(Boolean())
    gameversion_id = Column(Integer, ForeignKey('gameversion.id'))
    gameversion = relationship('GameVersion', viewonly=True, backref=backref('modversion', order_by=id))
    created = Column(DateTime)
    download_path = Column(String(512))
    changelog = Column(Unicode(10000))
    sort_index = Column(Integer)
    file_size = Column(Integer)

    def __init__(self, mod_id, friendly_version, gameversion_id, download_path,is_beta):
        self.friendly_version = friendly_version
        self.is_beta = is_beta
        self.gameversion_id = gameversion_id
        self.gameversion = GameVersion.query.filter(GameVersion.id == gameversion_id).first()
        self.download_path = download_path
        self.created = datetime.now()
        self.sort_index = 0
        self.file_size = 0
        self.mod_id = mod_id
        self.mod = Mod.query.filter(Mod.id == mod_id).first()
        
        if self.download_path:
            file_path = os.path.join(cfg['storage'], download_path)
            if os.path.isfile(file_path): self.file_size = os.path.getsize(file_path)

    def __repr__(self):
        return '<Mod Version %r>' % self.id


class Media(Base, MetaObject):
    __tablename__ = 'media'
    id = Column(Integer, primary_key = True)
    mod_id = Column(Integer, ForeignKey('mod.id'))
    mod = relationship('Mod', viewonly=True, backref=backref('media', order_by=id))
    hash = Column(String(12))
    type = Column(String(32))
    data = Column(String(512))

    def __init__(self, hash, type, data):
        self.hash = hash
        self.type = type
        self.data = data

    def __repr__(self):
        return '<Media %r>' % self.hash


class ReviewMedia(Base, MetaObject):
    __tablename__ = 'reviewmedia'
    id = Column(Integer, primary_key = True)
    review_id = Column(Integer, ForeignKey('review.id'))
    review = relationship('Review', viewonly=True, backref=backref('reviewmedia', order_by=id))
    hash = Column(String(12))
    type = Column(String(32))
    data = Column(String(512))

    def __init__(self, hash, type, data):
        self.hash = hash
        self.type = type
        self.data = data

    def __repr__(self):
        return '<ReviewMedia %r>' % self.hash


class GameVersion(Base, MetaObject):
    __tablename__ = 'gameversion'
    __lock__ = ['id', 'friendly_version', 'game_id', 'game']
    id = Column(Integer, primary_key = True)
    friendly_version = Column(String(128))
    is_beta = Column(Boolean())
    game_id = Column(Integer, ForeignKey('game.id'))
    game = relationship('Game', back_populates='version')

    def __init__(self, friendly_version, game_id, is_beta):
        self.friendly_version = friendly_version
        self.is_beta = is_beta
        self.game_id = game_id

    def __repr__(self):
        return '<Game Version %r>' % self.friendly_version