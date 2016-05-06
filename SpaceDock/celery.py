import smtplib
from celery import Celery
from email.mime.text import MIMEText
import redis
import requests
import time
import json
from SpaceDock.config import Config

cfg = Config()

app = Celery("tasks", broker=cfg["redis-connection"])
donation_cache = redis.Redis(host=cfg['patreon-host'], port=cfg['patreon-port'], db=cfg['patreon-db'])

def chunks(l, n):
    """ Yield successive n-sized chunks from l.
    """
    for i in range(0, len(l), n):
        yield l[i:i+n]

@app.task
def send_mail(sender, recipients, subject, message, important=False):
    if cfg["smtp-host"] == "":
        return
    smtp = smtplib.SMTP(host=cfg["smtp-host"], port=_cfgi("smtp-port"))
    if _cfgb("smtp-tls"):
        smtp.starttls()
    if cfg["smtp-user"] != "":
        smtp.login(cfg["smtp-user"], cfg["smtp-password"])
    message = MIMEText(message)
    if important:
        message['X-MC-Important'] = "true"
    message['X-MC-PreserveRecipients'] = "false"
    message['Subject'] = subject
    message['From'] = sender
    if len(recipients) > 1:
        message['Precedence'] = 'bulk'
    for group in chunks(recipients, 100):
        if len(group) > 1:
            message['To'] = "undisclosed-recipients:;"
        else:
            message['To'] = ";".join(group)
        print("Sending email from {} to {} recipients".format(sender, len(group)))
        smtp.sendmail(sender, group, message.as_string())
    smtp.quit()

@app.task
def update_patreon():
    donation_cache.set('patreon_update_time', time.time())
    if cfg['patreon_user_id'] != '' and cfg['patreon_campaign'] != '':
        r = requests.get("https://api.patreon.com/user/" + cfg['patreon_user_id'])
        if r.status_code == 200:
            patreon = json.loads(r.text)
            for linked_data in patreon['linked']:
                if 'creation_name' in linked_data and 'pledge_sum' in linked_data:
                    if linked_data['creation_name'] == cfg['patreon_campaign']:
                        donation_cache.set('patreon_donation_amount', linked_data['pledge_sum'])
