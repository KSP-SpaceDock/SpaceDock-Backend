from sqlalchemy import desc
from flask_login import current_user
from werkzeug.utils import secure_filename
from SpaceDock.common import *
from SpaceDock.config import cfg
from SpaceDock.database import db
from SpaceDock.email import *
from SpaceDock.formatting import mod_info, bulk, mod_version_info, rating_info
from SpaceDock.objects import *
from SpaceDock.routing import route

import json
import os
import time
import zipfile


@route('/api/mods')
def mod_list():
    """
    Returns a list of all mods
    """
    results = list()
    for mod in Mod.query.order_by(desc(Mod.id)).filter(Mod.published):
        results.append(mod_info(mod))
    return {'error': False, 'count': len(results), 'data': results}

@route('/api/mods/<gameshort>/<modid>')
def mods_info(gameshort, modid):
    """
    Returns information for one mod
    """
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        return {'error': True, 'reasons': ['The modid is invalid']}, 400
    if not Mod.query.filter(Mod.id == int(modid)).filter(Mod.game_id == game_id(gameshort)).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400
    # Get the mod
    mod = Mod.query.filter(Mod.id == int(modid)).first()
    return {'error': False, 'count': 1, 'data': mod_info(mod)}

@route('/api/mods/<gameshort>/<modid>/download/<versionname>')
@with_session
def mods_download(gameshort, modid, versionname):
    """
    Downloads the latest non-beta version of the mod
    """
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        return {'error': True, 'reasons': ['The modid is invalid']}, 400
    if not Mod.query.filter(Mod.id == int(modid)).filter(Mod.game_id == game_id(gameshort)).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400
    if not ModVersion.query.filter(ModVersion.mod_id == int(modid)).filter(ModVersion.friendly_version == versionname).first():
        return {'error': True, 'reasons': ['The version is invalid.']}, 400
    # Get the mod
    mod = Mod.query.filter(Mod.id == int(modid)).first()
    version = ModVersion.query.filter(ModVersion.mod_id == modid).filter(ModVersion.friendly_version == versionname).first()
    download = DownloadEvent.query\
            .filter(DownloadEvent.mod_id == mod.id and DownloadEvent.version_id == version.id)\
            .order_by(desc(DownloadEvent.created))\
            .first()

    # Check whether the path exists
    if not os.path.isfile(os.path.join(cfg['storage'], version.download_path)):
        return {'error': True, 'reasons': ['The file you tried to access doesn\'t exist.']}, 404

    if not 'Range' in request.headers:
        # Events are aggregated hourly
        if not download or ((datetime.now() - download.created).seconds / 60 / 60) >= 1:
            download = DownloadEvent()
            download.mod = mod
            download.version = version
            download.downloads = 1
            db.add(download)
            mod.downloads.append(download)
        else:
            download.downloads += 1
        mod.download_count += 1
    return redirect('/content/' + version.download_path)

@route('/api/mods/<gameshort>/<modid>/edit', methods=['POST'])
@user_has('mods-edit', params=['gameshort', 'modid'])
@with_session
def mod_edit(gameshort, modid):
    """
    Edits a mod, based on the request parameters. Required fields: data
    """
    errors = list()
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        errors.append('The Mod ID is invalid.')
    if not any(errors) and not Mod.query.filter(Mod.id == int(modid)).filter(Mod.game_id == game_id(gameshort)).first():
       errors.append('The gameshort is invalid.')
    if not request.form.get('data') or not is_json(request.form.get('data')):
        errors.append('The patch data is invalid.')
    if any(errors):
        return {'error': True, 'reasons': errors}, 400

    # Get variables
    parameters = json.loads(request.form.get('data'))

    # Get the matching mod and edit it
    mod = Mod.query.filter(Mod.id == int(modid)).first()
    code = edit_object(mod, parameters)

    # Error check
    if code == 2:
        return {'error': True, 'reasons': ['You tried to edit a value that doesn\'t exist.']}, 400
    elif code == 1:
        return {'error': True, 'reasons': ['You tried to edit a value that is marked as read-only.']}, 400
    else:
        return {'error': False, 'count': 1, 'data': mod_info(mod)}

@route('/api/mods/add', methods=['POST'])
@user_has('mods-add', params=['gameshort'])
@with_session
def add_mod():
    """
    Adds a mod, based on the request parameters. Required fields: name, gameshort, license
    """
    # Get variables
    name = request.form.get('name')
    short = request.form.get('gameshort')
    license = request.form.get('license')

    # Check the vars
    errors = list()
    if not name:
        errors.append('Invalid mod name.')
    if Mod.query.filter(Mod.name == name).first():
        errors.append('A mod with this name does already exist.')
    if not short or not game_id(short):
        errors.append('Invalid gameshort.')
    if not license:
        errors.append('Invalid License.')
    if any(errors):
        return {'error': True, 'reasons': errors}, 400

    # Add new mod
    mod = Mod(name, current_user.id, game_id(short), license)
    db.add(mod)
    current_user.add_roles(name)
    role = Role.query.filter(Role.name == name).first()
    role.add_abilities('mods-edit', 'mods-remove')
    role.add_param('mods-edit', 'modid', str(mod.id))
    role.add_param('mods-remove', 'name', name)
    db.add(role)
    return {'error': False, 'count': 1, 'data': mod_info(mod)}

@route('/api/mods/publish', methods=['POST'])
@user_has('mods-edit', params=['gameshort', 'name'])
@with_session
def publish_mod():
    """
    Makes a mod public. Required fields: name, gameshort
    """
    # Get variables
    name = request.form.get('name')
    short = request.form.get('gameshort')

    # Check the vars
    errors = list()
    if not name:
        errors.append('Invalid mod name.')
    if not short or not game_id(short):
        errors.append('Invalid gameshort.')
    if name and short and not Mod.query.filter(Mod.name == name).filter(Mod.game_id == game_id(short)).first():
        errors.append('A mod with theese parameters does not exist.')
    if any(errors):
        return {'error': True, 'reasons': errors}, 400

    # Publish
    mod = Mod.query.filter(Mod.name == name).filter(Mod.game_id == game_id(short)).first()
    mod.published = True
    mod.updated = datetime.now()
    return {'error': False}

@route('/api/mods/remove', methods=['POST'])
@user_has('mods-remove', params=['gameshort', 'name']) # We might want to allow deletion of own mods. Gameshort is here to allow per-game moderators.
@with_session
def remove_mod():
    """
    Removes a mod, based on the request parameters. Required fields: name, gameshort
    """
    # Get variables
    name = request.form.get('name')
    short = request.form.get('gameshort')

    # Check the vars
    errors = list()
    if not name:
        errors.append('Invalid mod name.')
    if not short or not game_id(short):
        errors.append('Invalid gameshort.')
    if name and short and not Mod.query.filter(Mod.name == name).filter(Mod.game_id == game_id(short)).first():
        errors.append('A mod with theese parameters does not exist.')
    if any(errors):
        return {'error': True, 'reasons': errors}, 400

    # Add new mod
    mod = Mod.query.filter(Mod.name == name).filter(Mod.game_id == game_id(short)).first()
    db.delete(mod)
    current_user.remove_roles(name)
    role = Role.query.filter(Role.name == name).first()
    role.remove_abilities('mods-edit', 'mods-remove')
    role.remove_param('mods-edit', 'modid', str(mod.id))
    role.remove_param('mods-remove', 'name', name)
    db.delete(role)
    return {'error': False}

@route('/api/mods/<gameshort>/<modid>/update-bg', methods=['POST'])
@user_has('mods-edit', params=['gameshort', 'modid'])
@with_session
def mod_updateBG(gameshort, modid):
    """
    Updates a mod background. Required fields: image
    """
    errors = list()
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        errors.append('The Mod ID is invalid.')
    if not any(errors) and not Mod.query.filter(Mod.id == int(modid)).filter(Mod.game_id == game_id(gameshort)).first():
       errors.append('The gameshort is invalid.')
    if not request.files.get('image'):
        errors.append('The background is invalid.')
    if any(errors):
        return {'error': True, 'reasons': errors}, 400

    # Find the mod
    mod = Mod.query.filter(Mod.id == int(modid)).first()

    # Get the file and save it to disk
    f = request.files['image']
    filetype = os.path.splitext(os.path.basename(f.filename))[1]
    if not filetype in ['.png', '.jpg']:
        return {'error': True, 'reasons': ['This file type is not acceptable.']}, 400
    filename = secure_filename(mod.name) + '-' + str(time.time()) + filetype
    base_path = os.path.join(secure_filename(mod.user.username) + '_' + str(mod.user.id), secure_filename(mod.name))
    full_path = os.path.join(cfg['storage'], base_path)
    if not os.path.exists(full_path):
        os.makedirs(full_path)
    path = os.path.join(full_path, filename)
    try:
        os.remove(os.path.join(cfg['storage'], mod.background))
    except:
        pass # who cares
    f.save(path)
    mod.background = os.path.join(base_path, filename)
    return {'error': False, 'count': 1, 'data': mod_info(mod)}

# Versions

@route('/api/mods/<gameshort>/<modid>/versions')
def mod_versions(gameshort, modid):
    """
    Returns a list of mod versions including their data.
    """
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        return {'error': True, 'reasons': ['The modid is invalid']}, 400
    if not Mod.query.filter(Mod.id == int(modid)).filter(Mod.game_id == game_id(gameshort)).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400
    # Get the versions
    versions = ModVersion.query.filter(ModVersion.mod_id == int(modid)).all()
    return {'error': False, 'count': len(versions), 'data': bulk(versions, mod_version_info)}

@route('/api/mods/<gameshort>/<modid>/versions/add', methods=['POST'])
@user_has('mods-edit', params=['gameshort', 'modid'])
@with_session
def mod_update(gameshort, modid):
    """
    Releases a new version of your mod. Required fields: version, game-version, notify-followers, is-beta, zipball. Optional fields: changelog
    """
    version = request.form.get('version')
    changelog = request.form.get('changelog')
    game_version = request.form.get('game-version')
    notify = request.form.get('notify-followers')
    beta = request.form.get('is-beta')
    zipball = request.files.get('zipball')

    # Get the mod
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        return {'error': True, 'reasons': ['The modid is invalid']}, 400
    if not Mod.query.filter(Mod.id == int(modid)).filter(Mod.game_id == game_id(gameshort)).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400
    mod = Mod.query.filter(Mod.id == int(modid)).first()

    # Process fields
    if not version or not game_version or not zipball:
        return {'error': True, 'reasons': ['All fields are required.']}, 400
    test_gameversion = GameVersion.query.filter(GameVersion.game_id == mod.game_id).filter(GameVersion.friendly_version == game_version).first()
    if not test_gameversion:
        return {'error': True, 'reasons': ['Game version does not exist.']}, 400
    game_version_id = test_gameversion.id
    notify = boolean(notify)

    # Save the file
    filename = secure_filename(mod.name) + '-' + secure_filename(version) + '.zip'
    base_path = os.path.join(secure_filename(current_user.username) + '_' + str(current_user.id), secure_filename(mod.name))
    full_path = os.path.join(cfg['storage'], base_path)
    if not os.path.exists(full_path):
        os.makedirs(full_path)
    path = os.path.join(full_path, filename)
    for v in mod.versions:
        if v.friendly_version == secure_filename(version):
            return {'error': True, 'reasons': ['We already have this version. Did you mistype the version number?']}, 400
    if os.path.isfile(path):
        os.remove(path)
    zipball.save(path)
    if not zipfile.is_zipfile(path):
        os.remove(path)
        return {'error': True, 'reasons': ['This is not a valid zip file.']}, 400
    version = ModVersion(mod.id, secure_filename(version), game_version_id, os.path.join(base_path, filename))
    version.changelog = changelog
    # Assign a sort index
    if len(mod.versions) == 0:
        version.sort_index = 0
    else:
        version.sort_index = max([v.sort_index for v in mod.versions]) + 1
    mod.versions.append(version)
    mod.updated = datetime.now()
    if notify:
        send_update_notification(mod, version, current_user)
    db.add(version)
    mod.default_version_id = version.id
    return {'error': False, 'count': 1, 'data': mod_version_info(version)}

@route('/api/mods/<gameshort>/<modid>/versions/delete', methods=['POST'])
@user_has('mods-edit', params=['gameshort', 'modid'])
@with_session
def delete_version(gameshort, modid):
    """
    Deletes a released version of the mod. Required fields: version-id
    """
    # Parameters
    versionid = request.form.get('version-id')

    # Error check
    errors = list()
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        errors.append('The mod ID is invalid.')
    if not any(errors) and not Mod.query.filter(Mod.id == int(modid)).filter(Mod.game_id == game_id(gameshort)).first():
       errors.append('The gameshort is invalid.')
    if not versionid.isdigit() or not ModVersion.query.filter(ModVersion.id == int(versionid)).first():
        errors.append('The version ID is invalid.')
    if not any(errors) and not ModVersion.query.filter(ModVersion.mod_id == int(modid)).filter(ModVersion.id == int(versionid)).first():
        errors.append('The mod ID and the version ID don\'t match.')
    if any(errors):
        return {'error': True, 'reasons': errors}, 400
    mod = Mod.query.filter(Mod.id == int(modid)).first()
    version = [v for v in mod.versions if v.id == int(versionid)]

    # Checks
    if len(mod.versions) == 1:
        return {'error': True, 'reasons': ['There is only one version left. You cant delete this one.']}, 400
    if len(version) == 0:
        return {'error': True, 'reasons': ['Something went wrong.']}, 404
    if version[0].id == mod.default_version_id:
        return {'error': True, 'reasons': ['You cannot delete the default version of a mod.']}, 400
    db.delete(version[0])
    mod.versions = [v for v in mod.versions if v.id != int(version_id)]
    return {'error': False}

@route('/api/mods/<gameshort>/<modid>/follow')
@user_has('logged-in', public=False)
@with_session
def mods_follow(gameshort, modid):
    """
    Registers a user for automated email sending when a new mod version is released
    """
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        return {'error': True, 'reasons': ['The modid is invalid']}, 400
    if not Mod.query.filter(Mod.id == int(modid)).filter(Mod.game_id == game_id(gameshort)).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 
    if any(m.id == int(modid) for m in current_user.following):
        return {'error': True, 'reasons': ['You are already following this mod.']}, 400

    # Get the mod
    mod = Mod.query.filter(Mod.id == int(modid)).first()

    # Follow
    event = FollowEvent.query\
            .filter(FollowEvent.mod_id == mod.id)\
            .order_by(desc(FollowEvent.created))\
            .first()
    # Events are aggregated hourly
    if not event or ((datetime.now() - event.created).seconds / 60 / 60) >= 1:
        event = FollowEvent()
        event.mod = mod
        event.delta = 1
        event.events = 1
        db.add(event)
        mod.follow_events.append(event)
    else:
        event.delta += 1
        event.events += 1
    mod.follower_count += 1
    current_user.following.append(mod)
    return {'error': False}

@route('/api/mods/<gameshort>/<modid>/unfollow')
@user_has('logged-in', public=False)
@with_session
def mods_unfollow(gameshort, modid):
    """
    Unregisters a user for automated email sending when a new mod version is released
    """
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        return {'error': True, 'reasons': ['The modid is invalid']}, 400
    if not Mod.query.filter(Mod.id == int(modid)).filter(Mod.game_id == game_id(gameshort)).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 
    if not any(m.id == int(modid) for m in current_user.following):
        return {'error': True, 'reasons': ['You are not following this mod.']}, 400

    # Get the mod
    mod = Mod.query.filter(Mod.id == int(modid)).first()

    # Follow
    event = FollowEvent.query\
            .filter(FollowEvent.mod_id == mod.id)\
            .order_by(desc(FollowEvent.created))\
            .first()
    # Events are aggregated hourly
    if not event or ((datetime.now() - event.created).seconds / 60 / 60) >= 1:
        event = FollowEvent()
        event.mod = mod
        event.delta = -1
        event.events = 1
        mod.follow_events.append(event)
        db.add(event)
    else:
        event.delta -= 1
        event.events += 1
    mod.follower_count -= 1
    current_user.following = [m for m in current_user.following if m.id != int(modid)]
    current_user.following.append(mod)
    return {'error': False}

@route('/api/mods/<gameshort>/<modid>/ratings/add')
@user_has('logged-in', public=False)
@with_session
def mods_rate(gameshort, modid):
	"""
	Rates a mod. Required fields: rating
	"""
	# Get variables
	score = request.form.get('rating')

	errors = list()
	if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
		errors.append('The Mod ID is invalid.')
	if not any(errors) and not Mod.query.filter(Mod.id == int(modid)).filter(Mod.game_id == game_id(gameshort)).first():
		errors.append('The gameshort is invalid.')
	if not score or not score.isdigit():
		 errors.append('The score is invalid.')
	if Rating.query.filter(Rating.mod_id == int(modid)).filter(Rating.user_id == current_user.id).first():
		errors.append('You already have a rating for this mod.')
	if any(errors):
		return {'error': True, 'reasons': errors}, 400

	# Find the mod
	mod = Mod.query.filter(Mod.id == int(modid)).first()
	# Create rating
	rating = Rating(current_user.id, current_user, mod.id, mod, int(score))
	db.add(rating)

	# Add rating to user and increase mod rating count
	current_user.ratings.append(rating)
	mod.rating_count += 1

	return {'error': False, 'count': 1, 'data': rating_info(rating)}

@route('/api/mods/<gameshort>/<modid>/ratings/remove')
@user_has('logged-in', public=False)
@with_session
def mods_unrate(gameshort, modid):
	"""
	Removes a rating for a mod.
	"""
	
	errors = list()
	if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
		errors.append('The Mod ID is invalid.')
	if not any(errors) and not Mod.query.filter(Mod.id == int(modid)).filter(Mod.game_id == game_id(gameshort)).first():
		errors.append('The gameshort is invalid.')
	if not Rating.query.filter(Rating.mod_id == int(modid)).filter(Rating.user_id == current_user.id).first():
		errors.append('You can\'t remove a rating you don\'t have, right?')
	if any(errors):
		return {'error': True, 'reasons': errors}, 400

	# Find the mod
	mod = Mod.query.filter(Mod.id == int(modid)).first()

	# Find the rating
	rating = Rating.query.filter(Rating.mod_id == mod.id).filter(Rating.user_id == current_user.id).first()

	# Remove the rating
	current_user.ratings.remove(rating)
	mod.rating_count -= 1

	return {'error': False}

@route('/api/mods/<gameshort>/<modid>/grant', methods=['POST'])
@user_has('mods-invite', params=['gameshort', 'modid'])
@with_session
def mods_grant(gameshort, modid):
    """
    Adds a new author to a mod. Required fields: username
    """
    username = request.form.get('username')

    # Check params
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        return {'error': True, 'reasons': ['The modid is invalid']}, 400
    if not Mod.query.filter(Mod.id == int(modid)).filter(Mod.game_id == game_id(gameshort)).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400
    if not User.query.filter(User.username == username).first():
        return {'error': true, 'reasons': ['The username is invalid']}, 400
    # Get the mod
    mod = Mod.query.filter(Mod.id == int(modid)).first()
    user = User.query.filter(User.username == username).first()

    # More checks
    if mod.user == user:
        return {'error': True, 'reasons': ['This user has already been added.']}, 400
    if any(m.user == user for m in mod.shared_authors):
        return {'error': True, 'reasons': ['This user has already been added.']}, 400
    if not user.public:
        return {'error': True, 'reasons': ['This user has not made their profile public.']}, 400
    if not mod.user == current_user:
        return {'error': True, 'reasons': ['You dont have the permission to add new authors.']}, 400
    author = SharedAuthor()
    author.mod = mod
    author.user = new_user
    mod.shared_authors.append(author)
    db.add(author)
    db.commit()
    send_grant_notice(mod, new_user)
    return {'error': False, 'count': 1, 'data': mod_info(mod)}

@route('/api/mods/<gameshort>/<modid>/accept_grant', methods=['POST'])
@user_has('logged-in')
@with_session
def mods_accept_grant(gameshort, modid):
    """
    Accepts a pending authorship grant for a mod. 
    """
    # Check params
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        return {'error': True, 'reasons': ['The modid is invalid']}, 400
    if not Mod.query.filter(Mod.id == int(modid)).filter(Mod.game_id == game_id(gameshort)).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400

    # Get the mod
    mod = Mod.query.filter(Mod.id == int(modid)).first()
    author = [a for a in mod.shared_authors if a.user == current_user and a.accepted]
    if len(author) == 0:
        return {'error': True, 'reasons': ['You do not have a pending authorship invite.']}, 400
    author = author[0]
    author.accepted = True
    current_user.add_role(mod.name)
    return {'error': False, 'count': 1, 'data': mod_info(mod)}

@route('/api/mods/<gameshort>/<modid>/reject_grant', methods=['POST'])
@user_has('logged-in')
@with_session
def mods_reject_grant(gameshort, modid):
    """
    Rejects a pending authorship grant for a mod. 
    """
    # Check params
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        return {'error': True, 'reasons': ['The modid is invalid']}, 400
    if not Mod.query.filter(Mod.id == int(modid)).filter(Mod.game_id == game_id(gameshort)).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400

    # Get the mod
    mod = Mod.query.filter(Mod.id == int(modid)).first()
    author = [a for a in mod.shared_authors if a.user == current_user and a.accepted]
    if len(author) == 0:
        return {'error': True, 'reasons': ['You do not have a pending authorship invite.']}, 400
    mod.shared_authors = [a for a in mod.shared_authors if a.user != current_user]
    db.delete(author)
    return {'error': False}

@api.route('/api/mods/<gameshort>/<modid>/revoke', methods=['POST'])
@user_has('mods-invite', params=['gameshort', 'modid'])
@with_session
def mods_revoke(gameshort, modid):
    """
    Removes an author from a mod. Required fields: username
    """
    username = request.form.get('username')

    # Check params
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        return {'error': True, 'reasons': ['The modid is invalid']}, 400
    if not Mod.query.filter(Mod.id == int(modid)).filter(Mod.game_id == game_id(gameshort)).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400
    if not User.query.filter(User.username == username).first():
        return {'error': true, 'reasons': ['The username is invalid']}, 400
    # Get the mod
    mod = Mod.query.filter(Mod.id == int(modid)).first()
    user = User.query.filter(User.username == username).first()

    # More checks
    if not mod.user == current_user:
        return {'error': True, 'reasons': ['You dont have the permission to remove authors.']}, 400
    if current_user == user:
        return {'error': True, 'reasons': ['You can\'t remove yourself.']}, 400
    if not any(m.user == user for m in mod.shared_authors):
        return { 'error': True, 'reason': 'This user is not an author.' }, 400

    # Remove
    author = [a for a in mod.shared_authors if a.user == new_user][0]
    mod.shared_authors = [a for a in mod.shared_authors if a.user != current_user]
    user.remove_role(mod.name)
    db.delete(author)
    return {'error': False}

