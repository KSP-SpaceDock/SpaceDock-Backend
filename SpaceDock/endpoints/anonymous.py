from flask import session
from flask.ext.login import current_user
from sqlalchemy import asc, desc
from SpaceDock.objects import Featured, BlogPost, Mod, ModVersion, Publisher, Game
from SpaceDock.common import *
import os.path

import math
import json

class AnonymousEndpoints:
    
    def __init__(self, cfg, db, search):
        self.cfg = cfg
        self.db = db.get_database()
        self.search = search
        
    def content(self, path):
        fullPath = self.cfg['storage'] + "/" +  path
        if not os.path.isfile(fullPath):
            abort(404)
        return send_from_directory(self.cfg['storage'] + "/", path)
    
    content.api_path = "/content/<path:path>"
    
    def get_default_game(self):
        current_game_short = self.cfg['default-game']
        game = Game.query.filter(Game.short == current_game_short).first()
        if not game:
            print('Please set a valid default-game in the config file!')
        return game
    
    def get_default_or_current_game(self):
        game = None
        if not 'game' in session:
            default_game = self.get_default_game()
            if not default_game:
                return None
            session['game'] = default_game.id
        current_game = session['game']
        return Game.query.filter(Game.id == current_game).first()
    
    def set_default_game(self, game):
        session['game'] = game.id
        
    def get_game_data(self, game):
        if not game:
            None
        mods = {}
        featured_raw = Featured.query.outerjoin(Mod).filter(Mod.published,Mod.game_id == game.id).order_by(desc(Featured.created))
        featured = []
        for featured_entry in featured_raw:
            featured.append(featured_entry.mod_id)
            mods[featured_entry.mod_id] = featured_entry.mod
        top_raw = Mod.query.filter(Mod.published,Mod.game_id == game.id).order_by(desc(Mod.download_count)).limit(6)[:6]
        top = []
        for top_entry in top_raw:
            top.append(top_entry.id)
            mods[top_entry.id] = top_entry
        new_raw = Mod.query.filter(Mod.published,Mod.game_id == game.id).order_by(desc(Mod.created)).limit(6)[:6]
        new = []
        for new_entry in new_raw:
            new.append(new_entry.id)
            mods[new_entry.id] = new_entry
        recent_raw = Mod.query.filter(Mod.published,Mod.game_id == game.id, ModVersion.query.filter(ModVersion.mod_id == Mod.id).count() > 1).order_by(desc(Mod.updated)).limit(6)[:6]
        recent = []
        for recent_entry in recent_raw:
            recent.append(recent.id)
            mods[recent.id] = recent
        user_count = User.query.count()
        mod_count = Mod.query.filter(Mod.game_id == game.id).count()
        yours = []
        if current_user:
            yours_raw = sorted(current_user.following, key=lambda m: m.updated, reverse=True)[:6]
            for yours_entry in yours_raw:
                yours.append(yours_entry.id)
                mods[yours_entry.id] = yours_entry
        send_mods = {}
        for mod_entry in mods:
            send_mods[mod.id] = mod_info(mod_entry)
                
        return({'featured': featured, 'top': top, 'new': new, 'recent': recent, 'user_count': user_count, 'mod_count': mod_count, 'yours': yours, 'mods': send_mods})
    
    def current_game_data(self):
        game = self.get_default_or_current_game()
        if not game:
            return jsonify({'error': True, 'reason': 'No game selected and default game not found'}), 400
        return jsonify({'error': False, 'game_id': game.id, 'data': self.get_game_data(game)})
    
    current_game_data.api_path = "/anon/game"
    
    def game_data(self, gameid):
        game = Game.query.filter(Game.id == gameid).first()
        if not game:
            return jsonify({'error': True, 'reason': 'Game does not exist'}), 400
        set_default_game(game)
        return jsonify({'error': False, 'game_id': game.id, 'data': self.get_game_data(game)})
    
    game_data.api_path = "/anon/game/<gameid>"
    
    def browse_new(self):
        game = self.get_default_or_current_game()
        if not game:
            return jsonify({'error': True, 'reason': 'No game selected and default game not found'}), 400
        mods = Mod.query.filter(Mod.published, Mod.game_id == game.id).order_by(desc(Mod.created)).all()
        #TODO: Pagify
        #total_pages = math.ceil(mods.count() / 30)
        #page = request.args.get('page')
        #if page:
        #    page = int(page)
        #    if page < 1:
        #        page = 1
        #    if page > total_pages:
        #        page = total_pages
        #else:
        #    page = 1
        #mods = mods.offset(30 * (page - 1)).limit(30)
        send_mods = {}
        for mod_entry in mods:
            send_mods[mod.id] = mod_info(mod_entry)
        return jsonify({'error': False, 'game_id': game.id, 'mods': send_mods})
    
    browse_new.api_path = "/anon/game/new"
    
    #@anonymous.route("/browse/new.rss")
    #def browse_new_rss():
    #    mods = Mod.query.filter(Mod.published).order_by(desc(Mod.created))
    #    mods = mods.limit(30)
    #    return Response(render_template("rss.xml", mods=mods, title="New mods on " + _cfg('site-name'),\
    #            description="The newest mods on " + _cfg('site-name'), \
    #            url="/browse/new"), mimetype="text/xml")
    
    def browse_updated(self):
        game = self.get_default_or_current_game()
        if not game:
            return jsonify({'error': True, 'reason': 'No game selected and default game not found'}), 400
        
        mods = Mod.query.filter(Mod.game_id == game.id, Mod.published, ModVersion.query.filter(ModVersion.mod_id == Mod.id).count() > 1).order_by(desc(Mod.updated)).all()
        #TODO: Pagify
        #total_pages = math.ceil(mods.count() / 30)
        #page = request.args.get('page')
        #if page:
        #    page = int(page)
        #    if page < 1:
        #        page = 1
        #    if page > total_pages:
        #        page = total_pages
        #else:
        #    page = 1
        #mods = mods.offset(30 * (page - 1)).limit(30)
        send_mods = {}
        for mod_entry in mods:
            send_mods[mod.id] = mod_info(mod_entry)
        return jsonify({'error': False, 'game_id': game.id, 'mods': send_mods})
    
    browse_updated.api_path = "/anon/game/updated"
    
    #@anonymous.route("/browse/updated.rss")
    #def browse_updated_rss():
    #    mods = Mod.query.filter(Mod.published, ModVersion.query.filter(ModVersion.mod_id == Mod.id).count() > 1).order_by(desc(Mod.updated))
    #    mods = mods.limit(30)
    #    return Response(render_template("rss.xml", mods=mods, title="Recently updated on " + _cfg('site-name'),\
    #            description="Mods on " + _cfg('site-name') + " updated recently", \
    #            url="/browse/updated"), mimetype="text/xml")
    
    def browse_top(self):
        game = self.get_default_or_current_game()
        if not game:
            return jsonify({'error': True, 'reason': 'No game selected and default game not found'}), 400
        mods = Mod.query.filter(Mod.game_id == game.id, Mod.published, ModVersion.query.filter(ModVersion.mod_id == Mod.id).count() > 1).order_by(desc(Mod.download_count)).all()
        #TODO: Pagify
        #page = request.args.get('page')
        #if page:
        #    page = int(page)
        #else:
        #    page = 1
        #mods, total_pages = search_mods(None, "", page, 30)
        send_mods = {}
        for mod_entry in mods:
            send_mods[mod.id] = mod_info(mod_entry)
        return jsonify({'error': False, 'game_id': game.id, 'mods': send_mods})
    
    browse_top.api_path = "/anon/game/top"
    
    def browse_featured(self):
        game = self.get_default_or_current_game()
        if not game:
            return jsonify({'error': True, 'reason': 'No game selected and default game not found'}), 400
        mods = [f.mod for f in Featured.query.outerjoin(Mod).filter(Mod.game_id == game.id).order_by(desc(Featured.created)).all()]
        #TODO: Pagify
        #total_pages = math.ceil(mods.count() / 30)
        #page = request.args.get('page')
        #if page:
        #    page = int(page)
        #    if page < 1:
        #        page = 1
        #    if page > total_pages:
        #        page = total_pages
        #else:
        #    page = 1
        #if page != 0:
        #    mods = mods.offset(30 * (page - 1)).limit(30)
        #mods = [f.mod for f in mods]
        send_mods = {}
        for mod_entry in mods:
            send_mods[mod.id] = mod_info(mod_entry)
        return jsonify({'error': False, 'game_id': game.id, 'mods': send_mods})
    
    browse_featured.api_path = "/anon/game/featured"
    
    #@anonymous.route("/browse/featured.rss")
    #def browse_featured_rss():
    #    mods = Featured.query.order_by(desc(Featured.created))
    #    mods = mods.limit(30)
    #    # Fix dates
    #    for f in mods:
    #        f.mod.created = f.created
    #    mods = [dumb_object(f.mod) for f in mods]
    #    db.rollback()
    #    return Response(render_template("rss.xml", mods=mods, title="Featured mods on " + _cfg('site-name'),\
    #            description="Featured mods on " + _cfg('site-name'), \
    #            url="/browse/featured"), mimetype="text/xml")
    
    def browse_all(self):
        game = self.get_default_or_current_game()
        if not game:
            return jsonify({'error': True, 'reason': 'No game selected and default game not found'}), 400
        mods = Mod.query.filter(Mod.game_id == game.id, Mod.published, ModVersion.query.filter(ModVersion.mod_id == Mod.id).count() > 1).order_by(asc(Mod.name)).all()
        #TODO: Pagify
        #
        #page = request.args.get('page')
        #if page:
        #    page = int(page)
        #else:
        #    page = 1
        #mods, total_pages = search_mods(None, "", page, 30)
        send_mods = {}
        for mod_entry in mods:
            send_mods[mod.id] = mod_info(mod_entry)
        return jsonify({'error': False, 'game_id': game.id, 'mods': send_mods})
    
    browse_all.api_path = "/anon/game/all"

    def singlegame_browse_new(self, gameid):
        game = Game.query.filter(Game.id == gameid).first()
        if not game:
            return jsonify({'error': True, 'reason': 'Game does not exist'}), 400
        mods = Mod.query.filter(Mod.published, Mod.game_id == game.id).order_by(desc(Mod.created)).all()
        #TODO: Pagify
        #total_pages = math.ceil(mods.count() / 30)
        #page = request.args.get('page')
        #if page:
        #    page = int(page)
        #    if page < 1:
        #        page = 1
        #    if page > total_pages:
        #        page = total_pages
        #else:
        #    page = 1
        #mods = mods.offset(30 * (page - 1)).limit(30)
        send_mods = {}
        for mod_entry in mods:
            send_mods[mod.id] = mod_info(mod_entry)
        return jsonify({'error': False, 'game_id': game.id, 'mods': send_mods})
    
    singlegame_browse_new.api_path = "/anon/game/<gameid>/new"
    
#    @anonymous.route("/json/<gameshort>/browse/<path:r>")
#    @json_output
#    def json_singlegame_browse_new(gameshort,r):
#        ra = r.split('/')
#        if not gameshort:
#            gameshort = 'kerbal-space-program'
#        ga = Game.query.filter(Game.short == gameshort).first()
#        session['game'] = ga.id;
#        session['gamename'] = ga.name;
#        session['gameshort'] = ga.short;
#        session['gameid'] = ga.id;
#        page = int(request.args.get('page'))
#        na = ""
#        rs = "/browse/all.rss"
#        ru = ga.short + "/browse/all"
#        if ra[0]:
#            if ra[0].lower() == "new":
#                mods = Mod.query.filter(Mod.published, Mod.game_id == ga.id).order_by(desc(Mod.created))
#                na = "Newest Mods"
#                rs = "/browse/new.rss"
#                ru = ga.short + "/browse/new"
#                total_pages = math.ceil(mods.count() / 30)
#            elif ra[0].lower() == "updated":
#                mods = Mod.query.filter(Mod.published, Mod.game_id == ga.id,ModVersion.query.filter(ModVersion.mod_id == Mod.id).count() > 1).order_by(desc(Mod.updated))
#                na = "Updated Mods"
#                rs = "/browse/updated.rss"
#                ru = ga.short + "/browse/updated"
#                total_pages = math.ceil(mods.count() / 30)
#            elif ra[0].lower() == "top":
#                na = "Top Mods"
#                rs = "/browse/top.rss"
#                ru = ga.short + "/browse/top"
#                mods = Mod.query.filter(Mod.game_id == ga.id).order_by(desc(Mod.follower_count))
#                total_pages = math.ceil(mods.count() / 30)
#            elif ra[0].lower() == "featured":
#                mods = Mod.query.filter(Mod.game_id == ga.id).join(Featured).order_by(desc(Featured.created))
#                na =" Featured Mods"
#                rs = "/browse/featured.rss"
#                ru = ga.short + "/browse/featured"
#                total_pages = math.ceil(mods.count() / 30)
#            else:
#                mods = Mod.query.filter(Mod.game_id == ga.id).order_by(desc(Mod.follower_count))
#                na = "All Mods"
#                rs = "/browse/all.rss"
#                ru = ga.short + "/browse/all"
#                total_pages = math.ceil(mods.count() / 30)
#        if page:
#            page = int(page)
#            if page < 1:
#                page = 1
#            if page > total_pages:
#                page = total_pages
#        else:
#            page = 1
#        
#        mods = mods.offset(30 * (page - 1)).limit(30)
#        mods = [e.serialize() for e in mods.all()]
#        #modsj = jsonify([e.serialize() for e in mods.all()])
#        #return { 'mods':mods, 'page':page, 'total_pages':total_pages,'ga':ga,'url':'/browse/new', 'name':'Newest Mods', 'rss':'/browse/new.rss'}
#        #return { 'mods':modsj, 'page':page, 'total_pages':total_pages,'ga':ga,'url':'/browse/new', 'name':'Newest Mods', 'rss':'/browse/new.rss'}
#    
#        return jsonify({"page":page,"total_pages":total_pages,"url":ru, "name":na, "rss":rs,"mods":mods})
#    
#    @anonymous.route("/<gameshort>/browse/new.rss")
#    def singlegame_browse_new_rss(gameshort):
#        if not gameshort:
#            gameshort = 'kerbal-space-program'
#        ga = Game.query.filter(Game.short == gameshort).first()
#        session['game'] = ga.id;
#        session['gamename'] = ga.name;
#        session['gameshort'] = ga.short;
#        session['gameid'] = ga.id;
#        mods = Mod.query.filter(Mod.published, Mod.game_id == ga.id).order_by(desc(Mod.created))
#        mods = mods.limit(30)
#        return Response(render_template("rss.xml", mods=mods, title="New mods on " + _cfg('site-name'),ga = ga,\
#                description="The newest mods on " + _cfg('site-name'), \
#                url="/browse/new"), mimetype="text/xml")
    
    def singlegame_browse_updated(self, gameid):
        game = Game.query.filter(Game.id == gameid).first()
        if not game:
            return jsonify({'error': True, 'reason': 'Game does not exist'}), 400
        
        mods = Mod.query.filter(Mod.game_id == game.id, Mod.published, ModVersion.query.filter(ModVersion.mod_id == Mod.id).count() > 1).order_by(desc(Mod.updated)).all()
        #TODO: Pagify
        #total_pages = math.ceil(mods.count() / 30)
        #page = request.args.get('page')
        #if page:
        #    page = int(page)
        #    if page < 1:
        #        page = 1
        #    if page > total_pages:
        #        page = total_pages
        #else:
        #    page = 1
        #mods = mods.offset(30 * (page - 1)).limit(30)
        send_mods = {}
        for mod_entry in mods:
            send_mods[mod.id] = mod_info(mod_entry)
        return jsonify({'error': False, 'game_id': game.id, 'mods': send_mods})
    
    singlegame_browse_updated.api_path = "/anon/game/<gameid>/updated"
    
    #@anonymous.route("/<gameshort>/browse/updated.rss")
    #def singlegame_browse_updated_rss(gameshort):
    #    if not gameshort:
    #        gameshort = 'kerbal-space-program'
    #    ga = Game.query.filter(Game.short == gameshort).first()
    #    session['game'] = ga.id;
    #    session['gamename'] = ga.name;
    #    session['gameshort'] = ga.short;
    #    session['gameid'] = ga.id;
    #    mods = Mod.query.filter(Mod.published,Mod.game_id == ga.id, ModVersion.query.filter(ModVersion.mod_id == Mod.id).count() > 1).order_by(desc(Mod.updated))
    #    mods = mods.limit(30)
    #    return Response(render_template("rss.xml", mods=mods, title="Recently updated on " + _cfg('site-name'),ga = ga,\
    #            description="Mods on " + _cfg('site-name') + " updated recently", \
    #            url="/browse/updated"), mimetype="text/xml")
    
    
    def singlegame_browse_top(self, gameid):
        game = Game.query.filter(Game.id == gameid).first()
        if not game:
            return jsonify({'error': True, 'reason': 'Game does not exist'}), 400
        
        mods = Mod.query.filter(Mod.game_id == game.id, Mod.published, ModVersion.query.filter(ModVersion.mod_id == Mod.id).count() > 1).order_by(desc(Mod.download_count)).all()
        #TODO: Pagify
        #total_pages = math.ceil(mods.count() / 30)
        #page = request.args.get('page')
        #if page:
        #    page = int(page)
        #    if page < 1:
        #        page = 1
        #    if page > total_pages:
        #        page = total_pages
        #else:
        #    page = 1
        #mods = mods.offset(30 * (page - 1)).limit(30)
        send_mods = {}
        for mod_entry in mods:
            send_mods[mod.id] = mod_info(mod_entry)
        return jsonify({'error': False, 'game_id': game.id, 'mods': send_mods})
    
    singlegame_browse_top.api_path = "/anon/game/<gameid>/top"
    
    def singlegame_browse_featured(self, gameid):
        game = Game.query.filter(Game.id == gameid).first()
        if not game:
            return jsonify({'error': True, 'reason': 'Game does not exist'}), 400
        
        mods = [f.mod for f in Featured.query.outerjoin(Mod).filter(Mod.game_id == game.id).order_by(desc(Featured.created)).all()]
        #TODO: Pagify
        #total_pages = math.ceil(mods.count() / 30)
        #page = request.args.get('page')
        #if page:
        #    page = int(page)
        #    if page < 1:
        #        page = 1
        #    if page > total_pages:
        #        page = total_pages
        #else:
        #    page = 1
        #mods = mods.offset(30 * (page - 1)).limit(30)
        send_mods = {}
        for mod_entry in mods:
            send_mods[mod.id] = mod_info(mod_entry)
        return jsonify({'error': False, 'game_id': game.id, 'mods': send_mods})
    
    singlegame_browse_featured.api_path = "/anon/game/<gameid>/featured"
    
    #@anonymous.route("/<gameshort>/browse/featured.rss")
    #def singlegame_browse_featured_rss(gameshort):
    #    if not gameshort:
    #        gameshort = 'kerbal-space-program'
    #    ga = Game.query.filter(Game.short == gameshort).first()
    #    session['game'] = ga.id;
    #    session['gamename'] = ga.name;
    #    session['gameshort'] = ga.short;
    #    session['gameid'] = ga.id;
    #    mods = Featured.query.outerjoin(Mod).filter(Mod.game_id == ga.id).order_by(desc(Featured.created))
    #    mods = mods.limit(30)
    #    # Fix dates
    #    for f in mods:
    #        f.mod.created = f.created
    #    mods = [dumb_object(f.mod) for f in mods]
    #    db.rollback()
    #    return Response(render_template("rss.xml", mods=mods, title="Featured mods on " + _cfg('site-name'),ga = ga,\
    #            description="Featured mods on " + _cfg('site-name'), \
    #            url="/browse/featured"), mimetype="text/xml")
    
    def singlegame_browse_all(self, gameid):
        game = Game.query.filter(Game.id == gameid).first()
        if not game:
            return jsonify({'error': True, 'reason': 'Game does not exist'}), 400
        mods = Mod.query.filter(Mod.game_id == game.id, Mod.published, ModVersion.query.filter(ModVersion.mod_id == Mod.id).count() > 1).order_by(asc(Mod.name)).all()
        #TODO: Pagify
        #
        #page = request.args.get('page')
        #if page:
        #    page = int(page)
        #else:
        #    page = 1
        #mods, total_pages = search_mods(None, "", page, 30)
        send_mods = {}
        for mod_entry in mods:
            send_mods[mod.id] = mod_info(mod_entry)
        return jsonify({'error': False, 'game_id': game.id, 'mods': send_mods})
    
    singlegame_browse_all.api_path = "/anon/game/<gameid>/all"
    
    
    #@anonymous.route("/about")
    #def about():
    #    return render_template("static/about.html")
    
    #@anonymous.route("/markdown")
    #def markdown_info():
    #    return render_template("static/markdown.html")
    
    #@anonymous.route("/privacy")
    #def privacy():
    #    return render_template("static/privacy.html")
    
    #@anonymous.route("/voip")
    #def voip():
    #    return render_template("static/voip.html")
    
    #@anonymous.route("/chat")
    #def chat():
    #    return render_template("static/chat.html")
    
    #@anonymous.route("/donate")
    #def donate():
    #    return render_template("static/donate.html")
    
    #@anonymous.route("/support")
    #def support():
    #    return render_template("static/support.html")
    
    def allgame_search(self):
        query = request.args.get('query')
        if not query:
            query = ''
        page = request.args.get('page')
        if page:
            page = int(page)
        else:
            page = 1
        mods, total_pages = self.search.search_mods(None, query, page, 30)
        send_mods = {}
        for mod_entry in mods:
            send_mods[mod.id] = mod_info(mod_entry)
        return jsonify({'error': False, 'page_id': page, 'pages': total_pages, 'mods': send_mods})
    allgame_search.api_path = "/anon/search"

    def singlegame_search(self, gameid):
        game = Game.query.filter(Game.id == gameid).first()
        if not game:
            return jsonify({'error': True, 'reason': 'Game does not exist'}), 400
        query = request.args.get('query')
        if not query:
            query = ''
        page = request.args.get('page')
        if page:
            page = int(page)
        else:
            page = 1
        mods, total_pages = self.search.search_mods(gameid, query, page, 30)
        send_mods = {}
        for mod_entry in mods:
            send_mods[mod.id] = mod_info(mod_entry)
        return jsonify({'error': False, 'page_id': page, 'pages': total_pages, 'mods': send_mods})
    
    singlegame_search.api_path = "/anon/search/<gameid>"
    